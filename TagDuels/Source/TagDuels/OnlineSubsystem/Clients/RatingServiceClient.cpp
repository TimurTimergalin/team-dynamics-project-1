#include "RatingServiceClient.h"

#include "HttpModule.h"
#include "Interfaces/IHttpResponse.h"
#include "TagDuels/OnlineSubsystem/utils.h"
#include "Dom/JsonObject.h"
#include "Serialization/JsonSerializer.h"

namespace
{
    // Convert HTTP response to optional int64 (display rating).
    TOptional<int64> ConvertGetRatingResponse(FHttpResponsePtr Response)
    {
        if (!Response.IsValid())
        {
            UE_LOG(LogTemp, Error, TEXT("GetRating: invalid response"));
            return {};
        }

        if (Response->GetResponseCode() != 200)
        {
            UE_LOG(LogTemp, Error, TEXT("GetRating: HTTP error %d, body: %s"),
                   Response->GetResponseCode(), *Response->GetContentAsString());
            return {};
        }

        FString JsonString = Response->GetContentAsString();
        TSharedPtr<FJsonObject> JsonObject;
        TSharedRef<TJsonReader<>> Reader = TJsonReaderFactory<>::Create(JsonString);

        if (!FJsonSerializer::Deserialize(Reader, JsonObject) || !JsonObject.IsValid())
        {
            UE_LOG(LogTemp, Error, TEXT("GetRating: failed to parse JSON: %s"), *JsonString);
            return {};
        }

        const TSharedPtr<FJsonObject>* RatingObject = nullptr;
        if (!JsonObject->TryGetObjectField(TEXT("rating"), RatingObject) || !RatingObject)
        {
            UE_LOG(LogTemp, Error, TEXT("GetRating: missing 'rating' field"));
            return {};
        }

        int64 DisplayValue = 0;
        if (!(*RatingObject)->TryGetNumberField(TEXT("display_value"), DisplayValue))
        {
            UE_LOG(LogTemp, Error, TEXT("GetRating: missing or invalid 'display_value' field"));
            return {};
        }

        return DisplayValue;
    }
}

RatingServiceClient::RatingServiceClient(const FString& Address)
    : Address(Address)
{
}

TSharedPtr<IHttpRequest> RatingServiceClient::CreateGetRatingRequest(int64 UserId) const
{
    FString Url = FString::Printf(TEXT("http://%s/v1/ratings/%lld"), *Address, UserId);

    TSharedRef<IHttpRequest> Request = FHttpModule::Get().CreateRequest();
    Request->SetURL(Url);
    Request->SetVerb(TEXT("GET"));
    Request->SetHeader(TEXT("Accept"), TEXT("application/json"));

    return Request;
}

TFuture<TOptional<int64>> RatingServiceClient::GetRating(int64 UserId) const
{
    TSharedPtr<IHttpRequest> Request = CreateGetRatingRequest(UserId);
    return MakeHttpRequest(Request).Next(ConvertGetRatingResponse);
}

TOptional<RatingServiceClient> CreateRatingServiceClient()
{
    FString Address;
    if (!GConfig)
    {
        UE_LOG(LogTemp, Error, TEXT("CreateRatingServiceClient: GConfig absent"));
        return {};
    }
    if (!GConfig->GetString(TEXT("OnlineSubsystemAddresses"), TEXT("RatingService"), Address, GGameIni))
    {
        UE_LOG(LogTemp, Error, TEXT("CreateRatingServiceClient: OnlineSubsystemAddresses:RatingService not set"));
        return {};
    }

    return RatingServiceClient(Address);
}
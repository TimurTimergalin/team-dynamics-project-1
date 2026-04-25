#include "MatchServiceClient.h"

#include "HttpModule.h"
#include "Interfaces/IHttpResponse.h"
#include "TagDuels/OnlineSubsystem/utils.h"
#include "Dom/JsonObject.h"
#include "Dom/JsonValue.h"
#include "Serialization/JsonSerializer.h"
#include "Serialization/JsonWriter.h"

namespace
{
    TOptional<FEndMatchResponse> ConvertEndMatchResponse(FHttpResponsePtr Response)
    {
        if (!Response.IsValid())
        {
            UE_LOG(LogTemp, Error, TEXT("EndMatch: invalid response"));
            return {};
        }

        if (Response->GetResponseCode() != 200)
        {
            UE_LOG(LogTemp, Error, TEXT("EndMatch: HTTP error %d, body: %s"),
                   Response->GetResponseCode(), *Response->GetContentAsString());
            return {};
        }

        FString JsonString = Response->GetContentAsString();
        TSharedPtr<FJsonObject> JsonObject;
        TSharedRef<TJsonReader<>> Reader = TJsonReaderFactory<>::Create(JsonString);

        if (!FJsonSerializer::Deserialize(Reader, JsonObject) || !JsonObject.IsValid())
        {
            UE_LOG(LogTemp, Error, TEXT("EndMatch: failed to parse JSON: %s"), *JsonString);
            return {};
        }

        FEndMatchResponse Result;
        JsonObject->TryGetNumberField(TEXT("newRating1"), Result.NewRating1);
        JsonObject->TryGetNumberField(TEXT("newRating2"), Result.NewRating2);
        return Result;
    }

    bool ConvertEmptyResponse(FHttpResponsePtr Response, const TCHAR* OperationName)
    {
        if (!Response.IsValid())
        {
            UE_LOG(LogTemp, Error, TEXT("%s: invalid response"), OperationName);
            return false;
        }

        if (Response->GetResponseCode() != 200)
        {
            UE_LOG(LogTemp, Error, TEXT("%s: HTTP error %d, body: %s"),
                   OperationName, Response->GetResponseCode(), *Response->GetContentAsString());
            return false;
        }

        return true;
    }

    TOptional<FString> ConvertRenewMatchResponse(FHttpResponsePtr Response)
    {
        if (!Response.IsValid())
        {
            UE_LOG(LogTemp, Error, TEXT("RenewMatch: invalid response"));
            return {};
        }

        if (Response->GetResponseCode() != 200)
        {
            UE_LOG(LogTemp, Error, TEXT("RenewMatch: HTTP error %d, body: %s"),
                   Response->GetResponseCode(), *Response->GetContentAsString());
            return {};
        }

        FString JsonString = Response->GetContentAsString();
        TSharedPtr<FJsonObject> JsonObject;
        TSharedRef<TJsonReader<>> Reader = TJsonReaderFactory<>::Create(JsonString);

        if (!FJsonSerializer::Deserialize(Reader, JsonObject) || !JsonObject.IsValid())
        {
            UE_LOG(LogTemp, Error, TEXT("RenewMatch: failed to parse JSON: %s"), *JsonString);
            return {};
        }

        FString MatchId;
        JsonObject->TryGetStringField(TEXT("matchId"), MatchId);
        return MatchId;
    }
}

MatchServiceClient::MatchServiceClient(const FString& Address)
    : Address(Address)
{
}

TSharedPtr<IHttpRequest> MatchServiceClient::CreateEndMatchRequest(const FEndMatchResult& MatchResult) const
{
    FString Url = FString::Printf(TEXT("http://%s/v1/matches/%s/end"), *Address, *MatchResult.MatchId);

    TSharedRef<FJsonObject> Body = MakeShared<FJsonObject>();

    if (MatchResult.WinnerId.IsSet())
    {
        Body->SetNumberField(TEXT("winnerId"), static_cast<double>(MatchResult.WinnerId.GetValue()));
    }

    TArray<TSharedPtr<FJsonValue>> RoundsArray;
    for (const FRoundData& Round : MatchResult.Rounds)
    {
        TSharedRef<FJsonObject> RoundObj = MakeShared<FJsonObject>();
        RoundObj->SetBoolField(TEXT("isPlayer1Killer"), Round.RoundKiller == RoundKiller::First);
        RoundObj->SetNumberField(TEXT("timeMillis"), static_cast<double>(Round.Duration.GetTotalMilliseconds()));
        RoundsArray.Add(MakeShared<FJsonValueObject>(RoundObj));
    }
    Body->SetArrayField(TEXT("rounds"), RoundsArray);

    FString BodyString;
    TSharedRef<TJsonWriter<>> Writer = TJsonWriterFactory<>::Create(&BodyString);
    FJsonSerializer::Serialize(Body, Writer);

    TSharedRef<IHttpRequest> Request = FHttpModule::Get().CreateRequest();
    Request->SetURL(Url);
    Request->SetVerb(TEXT("POST"));
    Request->SetHeader(TEXT("Content-Type"), TEXT("application/json"));
    Request->SetHeader(TEXT("Accept"), TEXT("application/json"));
    Request->SetContentAsString(BodyString);

    return Request;
}

TSharedPtr<IHttpRequest> MatchServiceClient::CreateCancelMatchRequest(const FString& MatchId) const
{
    FString Url = FString::Printf(TEXT("http://%s/v1/matches/%s/cancel"), *Address, *MatchId);

    TSharedRef<IHttpRequest> Request = FHttpModule::Get().CreateRequest();
    Request->SetURL(Url);
    Request->SetVerb(TEXT("POST"));
    Request->SetHeader(TEXT("Content-Type"), TEXT("application/json"));
    Request->SetContentAsString(TEXT("{}"));

    return Request;
}

TSharedPtr<IHttpRequest> MatchServiceClient::CreateRenewMatchRequest(const FString& MatchId) const
{
    FString Url = FString::Printf(TEXT("http://%s/v1/matches/%s/renew"), *Address, *MatchId);

    TSharedRef<IHttpRequest> Request = FHttpModule::Get().CreateRequest();
    Request->SetURL(Url);
    Request->SetVerb(TEXT("POST"));
    Request->SetHeader(TEXT("Content-Type"), TEXT("application/json"));
    Request->SetContentAsString(TEXT("{}"));

    return Request;
}

TFuture<TOptional<FEndMatchResponse>> MatchServiceClient::EndMatch(const FEndMatchResult& MatchResult) const
{
    TSharedPtr<IHttpRequest> Request = CreateEndMatchRequest(MatchResult);
    return MakeHttpRequest(Request).Next(ConvertEndMatchResponse);
}

TFuture<bool> MatchServiceClient::CancelMatch(const FString& MatchId) const
{
    TSharedPtr<IHttpRequest> Request = CreateCancelMatchRequest(MatchId);
    return MakeHttpRequest(Request).Next([](FHttpResponsePtr Response)
    {
        return ConvertEmptyResponse(Response, TEXT("CancelMatch"));
    });
}

TFuture<TOptional<FString>> MatchServiceClient::RenewMatch(const FString& MatchId) const
{
    TSharedPtr<IHttpRequest> Request = CreateRenewMatchRequest(MatchId);
    return MakeHttpRequest(Request).Next(ConvertRenewMatchResponse);
}

TOptional<MatchServiceClient> CreateMatchServiceClient()
{
    FString Address;
    if (!GConfig)
    {
        UE_LOG(LogTemp, Error, TEXT("CreateMatchServiceClient: GConfig absent"));
        return {};
    }
    if (!GConfig->GetString(TEXT("OnlineSubsystemAddresses"), TEXT("MatchService"), Address, GGameIni))
    {
        UE_LOG(LogTemp, Error, TEXT("CreateMatchServiceClient: OnlineSubsystemAddresses:MatchService not set"));
        return {};
    }

    return MatchServiceClient(Address);
}

#include "UserServiceClient.h"

#include "HttpModule.h"
#include "interfaces/IHttpResponse.h"
#include "TagDuels/OnlineSubsystem/utils.h"

namespace
{
	TOptional<FUserPlayerData> ConvertGetSelfDataResponse(FHttpResponsePtr Response)
	{
		if (!Response.IsValid())
		{
			UE_LOG(LogTemp, Error, TEXT("GetSelfData: invalid response"));
			return {};
		}

		if (Response->GetResponseCode() != 200)
		{
			UE_LOG(LogTemp, Error, TEXT("GetSelfData: HTTP error %d, body: %s"), Response->GetResponseCode(), *Response->GetContentAsString());
			return {};
		}

		FString JsonString = Response->GetContentAsString();
		TSharedPtr<FJsonObject> JsonObject;
		TSharedRef<TJsonReader<>> Reader = TJsonReaderFactory<>::Create(JsonString);

		if (!FJsonSerializer::Deserialize(Reader, JsonObject) || !JsonObject.IsValid())
		{
			UE_LOG(LogTemp, Error, TEXT("GetSelfData: failed to parse JSON: %s"), *JsonString);
			return {};
		}

		const TSharedPtr<FJsonObject>* UserDataObject = nullptr;
		if (!JsonObject->TryGetObjectField(TEXT("user_data"), UserDataObject) || !UserDataObject)
		{
			UE_LOG(LogTemp, Error, TEXT("GetSelfData: missing 'user_data' field"));
			return {};
		}

		FUserPlayerData Result;
		if (!(*UserDataObject)->TryGetNumberField(TEXT("id"), Result.Id))
		{
			UE_LOG(LogTemp, Error, TEXT("GetSelfData: missing or invalid 'id' field"));
			return {};
		}

		if (!(*UserDataObject)->TryGetStringField(TEXT("name"), Result.Name))
		{
			UE_LOG(LogTemp, Error, TEXT("GetSelfData: missing or invalid 'name' field"));
			return {};
		}

		return Result;
	}
}

UserServiceClient::UserServiceClient(const FString& Address): Address(Address) {}

TSharedPtr<IHttpRequest> UserServiceClient::GetSelfDataRequest(int64 SteamId) const
{
	// Build URL: http://<Address>/v1/self?key.steam_id=<SteamId>
	FString Url = FString::Printf(TEXT("http://%s/v1/self?key.steam_id=%lld"), *Address, SteamId);

	TSharedRef<IHttpRequest> Request = FHttpModule::Get().CreateRequest();
	Request->SetURL(Url);
	Request->SetVerb(TEXT("GET"));
	Request->SetHeader(TEXT("Accept"), TEXT("application/json"));

	return Request;
}

TFuture<TOptional<FUserPlayerData>> UserServiceClient::GetSelfData(int64 SteamId) const
{
	return MakeHttpRequest(GetSelfDataRequest(SteamId)).Next(ConvertGetSelfDataResponse);
}

TOptional<UserServiceClient> CreateUserServiceClient()
{
	FString Address;
	if (!GConfig)
	{
		UE_LOG(LogTemp, Error, TEXT("CreateUserServiceClient: GConfig absent"));
		return {};
	}
	if (!GConfig->GetString(TEXT("/Script/OnlineServices.Addresses"), TEXT("UserService"), Address, GGameIni))
	{
		UE_LOG(LogTemp, Error, TEXT("CreateUserServiceClient: /Script/OnlineServices.Addresses:UserService not set"));
		return {};
	}

	return UserServiceClient(Address);
}

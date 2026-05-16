#include "AuthServiceClient.h"

#include "HttpModule.h"
#include "Interfaces/IHttpResponse.h"
#include "Dom/JsonObject.h"
#include "Serialization/JsonSerializer.h"
#include "TagDuels/OnlineSubsystem/utils.h"

namespace
{
	TOptional<FAuthExternalResponse> ConvertAuthExternalResponse(FHttpResponsePtr Response)
	{
		if (!Response.IsValid())
		{
			UE_LOG(LogTemp, Error, TEXT("AuthExternal: invalid response"));
			return {};
		}

		if (Response->GetResponseCode() != 200)
		{
			UE_LOG(LogTemp, Error, TEXT("AuthExternal: HTTP error %d, body: %s"),
				Response->GetResponseCode(), *Response->GetContentAsString());
			return {};
		}

		TSharedPtr<FJsonObject> JsonObject;
		TSharedRef<TJsonReader<>> Reader = TJsonReaderFactory<>::Create(Response->GetContentAsString());
		if (!FJsonSerializer::Deserialize(Reader, JsonObject) || !JsonObject.IsValid())
		{
			UE_LOG(LogTemp, Error, TEXT("AuthExternal: failed to parse JSON: %s"), *Response->GetContentAsString());
			return {};
		}

		FAuthExternalResponse Result;
		JsonObject->TryGetStringField(TEXT("access"), Result.Access);
		JsonObject->TryGetStringField(TEXT("refresh"), Result.Refresh);

		const TSharedPtr<FJsonObject>* UserDataObject = nullptr;
		if (JsonObject->TryGetObjectField(TEXT("userData"), UserDataObject) && UserDataObject)
		{
			double Id = 0;
			if ((*UserDataObject)->TryGetNumberField(TEXT("id"), Id))
				Result.UserId = static_cast<int64>(Id);
			(*UserDataObject)->TryGetStringField(TEXT("name"), Result.UserName);
		}

		return Result;
	}

	TOptional<FRefreshResponse> ConvertRefreshResponse(FHttpResponsePtr Response)
	{
		if (!Response.IsValid())
		{
			UE_LOG(LogTemp, Error, TEXT("Refresh: invalid response"));
			return {};
		}

		if (Response->GetResponseCode() != 200)
		{
			UE_LOG(LogTemp, Error, TEXT("Refresh: HTTP error %d, body: %s"),
				Response->GetResponseCode(), *Response->GetContentAsString());
			return {};
		}

		TSharedPtr<FJsonObject> JsonObject;
		TSharedRef<TJsonReader<>> Reader = TJsonReaderFactory<>::Create(Response->GetContentAsString());
		if (!FJsonSerializer::Deserialize(Reader, JsonObject) || !JsonObject.IsValid())
		{
			UE_LOG(LogTemp, Error, TEXT("Refresh: failed to parse JSON: %s"), *Response->GetContentAsString());
			return {};
		}

		FRefreshResponse Result;
		JsonObject->TryGetStringField(TEXT("access"), Result.Access);
		JsonObject->TryGetStringField(TEXT("refresh"), Result.Refresh);
		return Result;
	}
}

AuthServiceClient::AuthServiceClient(const FString& Address) : Address(Address) {}

TFuture<TOptional<FAuthExternalResponse>> AuthServiceClient::AuthExternal(int64 SteamId, const FString& AuthToken) const
{
	FString Body = FString::Printf(TEXT("{\"external_key\":{\"steam_id\":%lld},\"auth_token\":\"%s\"}"), SteamId, *AuthToken);

	TSharedRef<IHttpRequest> Request = FHttpModule::Get().CreateRequest();
	Request->SetURL(FString::Printf(TEXT("http://%s/v1/auth/external"), *Address));
	Request->SetVerb(TEXT("POST"));
	Request->SetHeader(TEXT("Content-Type"), TEXT("application/json"));
	Request->SetContentAsString(Body);

	return MakeHttpRequest(Request).Next(ConvertAuthExternalResponse);
}

TFuture<TOptional<FAuthExternalResponse>> AuthServiceClient::AuthExternal(const FString& EosId, const FString& AuthToken) const
{
	FString Body = FString::Printf(TEXT("{\"external_key\":{\"eos_id\":\"%s\"},\"auth_token\":\"%s\"}"), *EosId, *AuthToken);

	TSharedRef<IHttpRequest> Request = FHttpModule::Get().CreateRequest();
	Request->SetURL(FString::Printf(TEXT("http://%s/v1/auth/external"), *Address));
	Request->SetVerb(TEXT("POST"));
	Request->SetHeader(TEXT("Content-Type"), TEXT("application/json"));
	Request->SetContentAsString(Body);

	return MakeHttpRequest(Request).Next(ConvertAuthExternalResponse);
}

TFuture<TOptional<FRefreshResponse>> AuthServiceClient::Refresh(const FString& RefreshToken) const
{
	FString Body = FString::Printf(TEXT("{\"refresh\":\"%s\"}"), *RefreshToken);

	TSharedRef<IHttpRequest> Request = FHttpModule::Get().CreateRequest();
	Request->SetURL(FString::Printf(TEXT("http://%s/v1/auth/refresh"), *Address));
	Request->SetVerb(TEXT("POST"));
	Request->SetHeader(TEXT("Content-Type"), TEXT("application/json"));
	Request->SetContentAsString(Body);

	return MakeHttpRequest(Request).Next(ConvertRefreshResponse);
}

TOptional<AuthServiceClient> CreateAuthServiceClient()
{
	FString Address;
	if (!GConfig)
	{
		UE_LOG(LogTemp, Error, TEXT("CreateAuthServiceClient: GConfig absent"));
		return {};
	}
	if (!GConfig->GetString(TEXT("OnlineSubsystemAddresses"), TEXT("AuthService"), Address, GGameIni))
	{
		UE_LOG(LogTemp, Error, TEXT("CreateAuthServiceClient: OnlineSubsystemAddresses:AuthService not set"));
		return {};
	}

	return AuthServiceClient(Address);
}

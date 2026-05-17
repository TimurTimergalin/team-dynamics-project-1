#include "UserServiceClient.h"

#include "HttpModule.h"
#include "GenericPlatform/GenericPlatformHttp.h"
#include "Interfaces/IHttpResponse.h"
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
		if (!JsonObject->TryGetObjectField(TEXT("userData"), UserDataObject) || !UserDataObject)
		{
			UE_LOG(LogTemp, Error, TEXT("GetSelfData: missing 'userData' field"));
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

TFuture<TOptional<FUserPlayerData>> UserServiceClient::GetUserData(int64 UserId, const FString& AuthToken) const
{
	FString Url = FString::Printf(TEXT("http://%s/v1/users/%lld"), *Address, UserId);
	TSharedRef<IHttpRequest> Request = FHttpModule::Get().CreateRequest();
	Request->SetURL(Url);
	Request->SetVerb(TEXT("GET"));
	Request->SetHeader(TEXT("Accept"), TEXT("application/json"));
	Request->SetHeader(TEXT("x-auth-token"), AuthToken);
	return MakeHttpRequest(Request).Next(ConvertGetSelfDataResponse);
}

TSharedPtr<IHttpRequest> UserServiceClient::GetFriendsRequest(const FString& Path, int64 UserId, const FString& PageKey, const FString& AuthToken) const
{
	FString Url = FString::Printf(TEXT("http://%s%s?user_id=%lld"), *Address, *Path, UserId);
	if (!PageKey.IsEmpty())
	{
		Url += FString::Printf(TEXT("&pagekey=%s"), *FGenericPlatformHttp::UrlEncode(PageKey));
	}
	TSharedRef<IHttpRequest> Request = FHttpModule::Get().CreateRequest();
	Request->SetURL(Url);
	Request->SetVerb(TEXT("GET"));
	Request->SetHeader(TEXT("Accept"), TEXT("application/json"));
	Request->SetHeader(TEXT("x-auth-token"), AuthToken);
	return Request;
}

namespace
{
	TOptional<FPlayersList> ConvertFriendsResponse(FHttpResponsePtr Response)
	{
		if (!Response.IsValid() || Response->GetResponseCode() != 200)
		{
			UE_LOG(LogTemp, Error, TEXT("UserServiceClient: friends request failed: %d"), Response.IsValid() ? Response->GetResponseCode() : -1);
			return {};
		}
		TSharedPtr<FJsonObject> JsonObject;
		TSharedRef<TJsonReader<>> Reader = TJsonReaderFactory<>::Create(Response->GetContentAsString());
		if (!FJsonSerializer::Deserialize(Reader, JsonObject) || !JsonObject.IsValid())
		{
			UE_LOG(LogTemp, Error, TEXT("UserServiceClient: failed to parse friends JSON"));
			return {};
		}
		FPlayersList Result;
		JsonObject->TryGetStringField(TEXT("pagekey"), Result.NextPageKey);
		const TArray<TSharedPtr<FJsonValue>>* FriendsArray;
		if (JsonObject->TryGetArrayField(TEXT("friends"), FriendsArray))
		{
			for (const TSharedPtr<FJsonValue>& FriendValue : *FriendsArray)
			{
				const TSharedPtr<FJsonObject>* FriendObj;
				if (!FriendValue->TryGetObject(FriendObj)) continue;
				const TSharedPtr<FJsonObject>* UserObj;
				if (!(*FriendObj)->TryGetObjectField(TEXT("user"), UserObj)) continue;
				FUserPlayerData Player;
				double Id;
				if ((*UserObj)->TryGetNumberField(TEXT("id"), Id))
					Player.Id = static_cast<int64>(Id);
				(*UserObj)->TryGetStringField(TEXT("name"), Player.Name);
				Result.Players.Add(Player);
			}
		}
		return Result;
	}
}

TFuture<TOptional<FPlayersList>> UserServiceClient::GetFriends(int64 UserId, const FString& PageKey, const FString& AuthToken) const
{
	return MakeHttpRequest(GetFriendsRequest(TEXT("/v1/friends"), UserId, PageKey, AuthToken)).Next(ConvertFriendsResponse);
}

TFuture<TOptional<FPlayersList>> UserServiceClient::GetIncomingRequests(int64 UserId, const FString& PageKey, const FString& AuthToken) const
{
	return MakeHttpRequest(GetFriendsRequest(TEXT("/v1/friends/incoming"), UserId, PageKey, AuthToken)).Next(ConvertFriendsResponse);
}

TFuture<TOptional<FPlayersList>> UserServiceClient::GetOutgoingRequests(int64 UserId, const FString& PageKey, const FString& AuthToken) const
{
	return MakeHttpRequest(GetFriendsRequest(TEXT("/v1/friends/outgoing"), UserId, PageKey, AuthToken)).Next(ConvertFriendsResponse);
}

TSharedPtr<IHttpRequest> UserServiceClient::MutationRequest(const FString& Verb, const FString& Path, int64 UserId, int64 OtherUserId, const FString& AuthToken) const
{
	FString Url = FString::Printf(TEXT("http://%s%s"), *Address, *Path);
	FString Body = FString::Printf(TEXT("{\"user_id\":%lld,\"other_user_id\":%lld}"), UserId, OtherUserId);
	TSharedRef<IHttpRequest> Request = FHttpModule::Get().CreateRequest();
	Request->SetURL(Url);
	Request->SetVerb(Verb);
	Request->SetHeader(TEXT("Content-Type"), TEXT("application/json"));
	Request->SetHeader(TEXT("x-auth-token"), AuthToken);
	Request->SetContentAsString(Body);
	return Request;
}

namespace
{
	bool ConvertMutationResponse(FHttpResponsePtr Response)
	{
		if (!Response.IsValid() || Response->GetResponseCode() != 200)
		{
			UE_LOG(LogTemp, Error, TEXT("UserServiceClient: mutation failed: %d"), Response.IsValid() ? Response->GetResponseCode() : -1);
			return false;
		}
		return true;
	}
}

TFuture<bool> UserServiceClient::AddFriend(int64 UserId, int64 OtherUserId, const FString& AuthToken) const
{
	return MakeHttpRequest(MutationRequest(TEXT("POST"), TEXT("/v1/friends"), UserId, OtherUserId, AuthToken)).Next(ConvertMutationResponse);
}

TFuture<bool> UserServiceClient::RemoveFriend(int64 UserId, int64 OtherUserId, const FString& AuthToken) const
{
	return MakeHttpRequest(MutationRequest(TEXT("POST"), TEXT("/v1/friends/remove"), UserId, OtherUserId, AuthToken)).Next(ConvertMutationResponse);
}

TOptional<UserServiceClient> CreateUserServiceClient()
{
	FString Address;
	if (!GConfig)
	{
		UE_LOG(LogTemp, Error, TEXT("CreateUserServiceClient: GConfig absent"));
		return {};
	}
	if (!GConfig->GetString(TEXT("OnlineSubsystemAddresses"), TEXT("UserService"), Address, GGameIni))
	{
		UE_LOG(LogTemp, Error, TEXT("CreateUserServiceClient: OnlineSubsystemAddresses:UserService not set"));
		return {};
	}

	return UserServiceClient(Address);
}

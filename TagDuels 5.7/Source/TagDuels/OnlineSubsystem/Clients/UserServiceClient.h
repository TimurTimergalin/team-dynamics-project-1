#pragma once
#include "TagDuels/OnlineSubsystem/Contract/UserData.h"
#include "Interfaces/IHttpRequest.h"

class UserServiceClient
{
public:
	TFuture<TOptional<FUserPlayerData>> GetSelfData(int64 SteamId) const;
	TFuture<TOptional<FPlayersList>> GetFriends(int64 UserId, const FString& PageKey) const;
	TFuture<TOptional<FPlayersList>> GetIncomingRequests(int64 UserId, const FString& PageKey) const;
	TFuture<TOptional<FPlayersList>> GetOutgoingRequests(int64 UserId, const FString& PageKey) const;
	TFuture<bool> AddFriend(int64 UserId, int64 OtherUserId) const;
	TFuture<bool> RemoveFriend(int64 UserId, int64 OtherUserId) const;
private:
	explicit UserServiceClient(const FString& Address);
	TSharedPtr<IHttpRequest> GetSelfDataRequest(int64 SteamId) const;
	TSharedPtr<IHttpRequest> GetFriendsRequest(const FString& Path, int64 UserId, const FString& PageKey) const;
	TSharedPtr<IHttpRequest> MutationRequest(const FString& Verb, const FString& Path, int64 UserId, int64 OtherUserId) const;

	FString Address;
	friend TOptional<UserServiceClient> CreateUserServiceClient();
};

TOptional<UserServiceClient> CreateUserServiceClient();

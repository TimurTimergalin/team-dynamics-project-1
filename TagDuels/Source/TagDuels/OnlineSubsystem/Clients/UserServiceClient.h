#pragma once
#include "TagDuels/OnlineSubsystem/Contract/UserData.h"
#include "Interfaces/IHttpRequest.h"

class UserServiceClient
{
public:
	TFuture<TOptional<FUserPlayerData>> GetSelfData(int64 SteamId) const;
private:
	explicit UserServiceClient(const FString& Address);
	TSharedPtr<IHttpRequest> GetSelfDataRequest(int64 SteamId) const;

	FString Address;
	friend TOptional<UserServiceClient> CreateUserServiceClient();
};

TOptional<UserServiceClient> CreateUserServiceClient();

#pragma once

#include "CoreMinimal.h"
#include "Interfaces/IHttpRequest.h"

struct FAuthExternalResponse
{
	FString Access;
	FString Refresh;
	int64 UserId = 0;
	FString UserName;
};

struct FRefreshResponse
{
	FString Access;
	FString Refresh;
};

class AuthServiceClient
{
public:
	TFuture<TOptional<FAuthExternalResponse>> AuthExternal(int64 SteamId, const FString& AuthToken) const;
	TFuture<TOptional<FAuthExternalResponse>> AuthExternal(const FString& EosId, const FString& AuthToken) const;
	TFuture<TOptional<FRefreshResponse>> Refresh(const FString& RefreshToken) const;

private:
	explicit AuthServiceClient(const FString& Address);

	FString Address;

	friend TOptional<AuthServiceClient> CreateAuthServiceClient();
};

TOptional<AuthServiceClient> CreateAuthServiceClient();

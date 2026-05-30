#pragma once

#include "CoreMinimal.h"
#include "Interfaces/IHttpRequest.h"
#include "TagDuels/OnlineSubsystem/Contract/MatchHistory.h"

class MatchHistoryServiceClient
{
public:
	TFuture<TOptional<FMatchHistoryPage>> GetMatchHistory(int64 UserId, const FString& PageKey = FString()) const;

private:
	explicit MatchHistoryServiceClient(const FString& Address);

	TSharedPtr<IHttpRequest> CreateMatchHistoryRequest(int64 UserId, const FString& PageKey) const;

	FString Address;

	friend TOptional<MatchHistoryServiceClient> CreateMatchHistoryServiceClient();
};

TOptional<MatchHistoryServiceClient> CreateMatchHistoryServiceClient();

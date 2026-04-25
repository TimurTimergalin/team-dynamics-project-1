#pragma once

#include "CoreMinimal.h"
#include "Interfaces/IHttpRequest.h"
#include "TagDuels/OnlineSubsystem/Contract/Match.h"

class MatchServiceClient
{
public:
	TFuture<TOptional<FEndMatchResponse>> EndMatch(const FEndMatchResult& MatchResult) const;
	TFuture<bool> CancelMatch(const FString& MatchId) const;
	TFuture<TOptional<FString>> RenewMatch(const FString& MatchId) const;

private:
	explicit MatchServiceClient(const FString& Address);

	TSharedPtr<IHttpRequest> CreateEndMatchRequest(const FEndMatchResult& MatchResult) const;
	TSharedPtr<IHttpRequest> CreateCancelMatchRequest(const FString& MatchId) const;
	TSharedPtr<IHttpRequest> CreateRenewMatchRequest(const FString& MatchId) const;

	FString Address;

	friend TOptional<MatchServiceClient> CreateMatchServiceClient();
};

TOptional<MatchServiceClient> CreateMatchServiceClient();

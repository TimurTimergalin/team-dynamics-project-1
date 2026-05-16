#pragma once

#include "CoreMinimal.h"
#include "Interfaces/IHttpRequest.h"
#include "TagDuels/OnlineSubsystem/Contract/MatchHistory.h"

class MatchHistoryServiceV2Client
{
public:
    TFuture<TOptional<FMatchHistoryPage>> GetMatchHistory(int64 UserId, const FString& AuthToken, const FString& PageKey = FString()) const;
    TFuture<TOptional<int64>>             GetRating(int64 UserId, const FString& AuthToken) const;

private:
    explicit MatchHistoryServiceV2Client(const FString& Address);

    TSharedPtr<IHttpRequest> CreateGetMatchHistoryRequest(int64 UserId, const FString& PageKey, const FString& AuthToken) const;
    TSharedPtr<IHttpRequest> CreateGetRatingRequest(int64 UserId, const FString& AuthToken) const;

    FString Address;

    friend TOptional<MatchHistoryServiceV2Client> CreateMatchHistoryServiceV2Client();
};

TOptional<MatchHistoryServiceV2Client> CreateMatchHistoryServiceV2Client();

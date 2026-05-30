#pragma once

#include "CoreMinimal.h"
#include "Interfaces/IHttpRequest.h"

class RatingServiceClient
{
public:
	TFuture<TOptional<int64>> GetRating(int64 UserId) const;

private:
	explicit RatingServiceClient(const FString& Address);
	TSharedPtr<IHttpRequest> CreateGetRatingRequest(int64 UserId) const;

	FString Address;

	friend TOptional<RatingServiceClient> CreateRatingServiceClient();
};

TOptional<RatingServiceClient> CreateRatingServiceClient();

#pragma once

#include "CoreMinimal.h"
#include "TagDuels/OnlineSubsystem/Contract/MatchHistory.h"
#include "Match.generated.h"

USTRUCT(BlueprintType)
struct TAGDUELS_API FEndMatchResult
{
	GENERATED_BODY()

	UPROPERTY(BlueprintReadWrite)
	FString MatchId{};

	// Unset means draw
	TOptional<int64> WinnerId{};

	UPROPERTY(BlueprintReadWrite)
	TArray<FRoundData> Rounds{};
};

USTRUCT(BlueprintType)
struct TAGDUELS_API FEndMatchResponse
{
	GENERATED_BODY()

	UPROPERTY(BlueprintReadOnly)
	int64 NewRating1{};

	UPROPERTY(BlueprintReadOnly)
	int64 NewRating2{};
};

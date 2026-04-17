#pragma once

#include "CoreMinimal.h"
#include "MatchHistory.generated.h"

UENUM(BlueprintType)
enum class RoundKiller : uint8
{
	First = 0,
	Second = 1,
};

UENUM(BlueprintType)
enum class MatchResolution : uint8
{
	FirstWins = 0,
	SecondWins = 1,
	Draw = 2,
};

USTRUCT(BlueprintType)
struct TAGDUELS_API FRoundData
{
	GENERATED_BODY()

	UPROPERTY(BlueprintReadOnly)
	RoundKiller RoundKiller{};

	UPROPERTY(BlueprintReadOnly)
	FTimespan Duration{};
};

USTRUCT(BlueprintType)
struct TAGDUELS_API FMatchHistoryPlayerData
{
	GENERATED_BODY()
	UPROPERTY(BlueprintReadOnly)
	int64 Id{};

	UPROPERTY(BlueprintReadOnly)
	FString Name{};

	UPROPERTY(BlueprintReadOnly)
	int64 Rating{};
};

USTRUCT(BlueprintType)
struct TAGDUELS_API FMatchHistory
{
	GENERATED_BODY()
	TArray<FRoundData> Rounds{};

	UPROPERTY(BlueprintReadOnly)
	FMatchHistoryPlayerData Player1{};

	UPROPERTY(BlueprintReadOnly)
	FMatchHistoryPlayerData Player2{};

	UPROPERTY(BlueprintReadOnly)
	MatchResolution Resolution{};

	UPROPERTY(BlueprintReadOnly)
	FDateTime EndTime{};

	UPROPERTY(BlueprintReadOnly)
	FString MatchId{};
};

USTRUCT(BlueprintType)
struct TAGDUELS_API FMatchHistoryPage
{
	GENERATED_BODY()
	UPROPERTY(BlueprintReadOnly)
	FString NextPageKey{};

	UPROPERTY(BlueprintReadOnly)
	TArray<FMatchHistory> Matches{};
};

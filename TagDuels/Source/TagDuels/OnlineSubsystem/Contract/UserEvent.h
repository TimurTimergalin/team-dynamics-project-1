#pragma once

#include "UserEvent.generated.h"

UENUM(BlueprintType)
enum class EUserStatus : uint8
{
	Offline,
	Online,
	InGame,
};

USTRUCT(BlueprintType)
struct FStatusUpdatedEvent
{
	GENERATED_BODY()

	UPROPERTY(BlueprintReadOnly)
	int64 UserId{};

	UPROPERTY(BlueprintReadOnly)
	EUserStatus Status{};
};

USTRUCT(BlueprintType)
struct FChallengeReceivedEvent
{
	GENERATED_BODY()

	UPROPERTY(BlueprintReadOnly)
	int64 UserId{};

	UPROPERTY(BlueprintReadOnly)
	FString UserName{};

	FString MessageId{};
};

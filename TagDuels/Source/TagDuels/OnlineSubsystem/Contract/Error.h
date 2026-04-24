#pragma once

#include "CoreMinimal.h"
#include "Error.generated.h"

UENUM(BlueprintType)
enum class EOnlineErrorType : uint8
{
	NonCritical = 0,
	Critical     = 1,
};

USTRUCT(BlueprintType)
struct TAGDUELS_API FOnlineSubsystemError
{
	GENERATED_BODY()

	UPROPERTY(BlueprintReadOnly)
	FString Message{};

	UPROPERTY(BlueprintReadOnly)
	EOnlineErrorType Type = EOnlineErrorType::NonCritical;
};

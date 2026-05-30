#pragma once

#include "CoreMinimal.h"
#include "ServerAnnotations.generated.h"

USTRUCT(BlueprintType)
struct TAGDUELS_API FUserAnnotations
{
	GENERATED_BODY()

	UPROPERTY(BlueprintReadOnly)
	int64 PlayerId{};

	UPROPERTY(BlueprintReadOnly)
	FString Name{};

	UPROPERTY(BlueprintReadOnly)
	int64 Rating{};
};

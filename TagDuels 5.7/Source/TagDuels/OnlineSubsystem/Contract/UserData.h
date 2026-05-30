#pragma once

#include "CoreMinimal.h"
#include "UserData.generated.h"

USTRUCT(BlueprintType)
struct TAGDUELS_API FUserPlayerData
{
	GENERATED_BODY()
	UPROPERTY(BlueprintReadOnly)
	int64 Id{};

	UPROPERTY(BlueprintReadOnly)
	FString Name{};
};

USTRUCT(BlueprintType)
struct TAGDUELS_API FPlayersList
{
	GENERATED_BODY()
	UPROPERTY(BlueprintReadOnly)
	TArray<FUserPlayerData> Players{};

	UPROPERTY(BlueprintReadOnly)
	FString NextPageKey{};
};

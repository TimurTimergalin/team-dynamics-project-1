// Fill out your copyright notice in the Description page of Project Settings.

#pragma once

#include "CoreMinimal.h"
#include "Contract/MatchHistory.h"
#include "Kismet/BlueprintFunctionLibrary.h"
#include "OSSBlueprintFunctionLibrary.generated.h"

/**
 * 
 */
UCLASS()
class TAGDUELS_API UOSSBlueprintFunctionLibrary : public UBlueprintFunctionLibrary
{
	GENERATED_BODY()
public:
	UFUNCTION(BlueprintPure, Category="Round Data")
	static FRoundData MakeRoundData(RoundKiller Killer, FTimespan Duration)
	{
		return FRoundData{Killer, Duration};
	}
};

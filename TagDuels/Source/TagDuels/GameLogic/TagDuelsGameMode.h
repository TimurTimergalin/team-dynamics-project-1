// Fill out your copyright notice in the Description page of Project Settings.

#pragma once

#include "CoreMinimal.h"
#include "GameFramework/GameModeBase.h"
#include "TagDuelsGameMode.generated.h"

/**
 * 
 */
UCLASS()
class TAGDUELS_API ATagDuelsGameMode : public AGameModeBase
{
	GENERATED_BODY()
public:
	virtual void PreLogin(const FString& Options, const FString& Address, const FUniqueNetIdRepl& UniqueId, FString& ErrorMessage) override;

	UFUNCTION(BlueprintNativeEvent, BlueprintCallable, Category = "Login")
	void OnPreLogin(const FString& Options, const FString& Address, const FUniqueNetIdRepl& UniqueId);
};

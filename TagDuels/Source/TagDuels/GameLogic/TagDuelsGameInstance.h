// Fill out your copyright notice in the Description page of Project Settings.

#pragma once

#include "CoreMinimal.h"
#include "Engine/GameInstance.h"
#include "TagDuels/Enums/TagDuelsEnums.h"
#include "TagDuelsGameInstance.generated.h"

/**
 * 
 */
UCLASS()
class TAGDUELS_API UTagDuelsGameInstance : public UGameInstance
{
	GENERATED_BODY()

public:
	// Base Methods
	virtual void OnStart() override;

	UFUNCTION(BlueprintNativeEvent, BlueprintCallable, Category = "PlayerController")
	void BlStart();
	
	// Steam OSS
	UFUNCTION(BlueprintCallable, Category = "SteamAuth")
	FString GetSteamAuthToken();

	UFUNCTION(BlueprintCallable, Category = "SteamAuth")
	FString GetSteamID();

	// EpicGames OSS
	UFUNCTION(BlueprintCallable, Category = "EGSAuth")
	void LoginToEOS();

	// Online
	UFUNCTION(BlueprintCallable, Category = "Online")
	EOnlineSubsystemType GetActiveOnlineSubsystemType() const;

private:
	// EpicGames OSS
	void OnLoginComplete(int32 LocalUserNum, bool bWasSuccessful, const FUniqueNetId& UserId, const FString& Error);
};

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
	// Basic Methods
	virtual void OnStart() override;

	UFUNCTION(BlueprintNativeEvent, BlueprintCallable, Category = "PlayerController")
	void BlStart();
	
	// OSS
	UFUNCTION(BlueprintCallable, Category = "OSS")
	void LoginToSteam();
	
	UFUNCTION(BlueprintCallable, Category = "OSS")
	void LoginToEOS();

	UFUNCTION(BlueprintImplementableEvent, Category = "OSS")
	void OnSuccessfulLoginSteam(const FString& AccountID, const FString& AuthToken);
	
	UFUNCTION(BlueprintImplementableEvent, Category = "OSS")
	void OnSuccessfulLoginEOS(const FString& AccountID, const FString& AuthToken);
	
	UFUNCTION(BlueprintCallable, Category = "OSS")
	EOnlineSubsystemType GetActiveOnlineSubsystemType() const;

private:
	// OSS
	void OnPersistentLoginComplete(int32 LocalUserNum, bool bWasSuccessful, const FUniqueNetId& UserId, const FString& Error);

	void OnLoginComplete(int32 LocalUserNum, bool bWasSuccessful, const FUniqueNetId& UserId, const FString& Error);
    
	void HandleSuccessfulLogin(const FUniqueNetId& UserId, int32 LocalUserNum);

	FDelegateHandle PersistentLoginDelegateHandle;
	FDelegateHandle LoginDelegateHandle;
};

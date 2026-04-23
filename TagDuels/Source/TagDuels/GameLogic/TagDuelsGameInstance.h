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
	void Start();
	
	// OSS
	UFUNCTION(BlueprintCallable, Category = "OSS")
	void LoginToSteam();
	
	UFUNCTION(BlueprintCallable, Category = "OSS")
	void LoginToEOS();

	UFUNCTION(BlueprintImplementableEvent, Category = "OSS")
	void OnSuccessfulLoginSteam(const int64 AccountID, const FString& AuthToken);
	
	UFUNCTION(BlueprintImplementableEvent, Category = "OSS")
	void OnSuccessfulLoginEOS(const FString& AccountID, const FString& AuthToken);
	
	UFUNCTION(BlueprintCallable, Category = "OSS")
	EOnlineSubsystemType GetActiveOnlineSubsystemType() const;

	UFUNCTION(BlueprintImplementableEvent, Category = "OSS")
	void OnFailedToLogin(const FString& ErrorMessage);

private:
	// OSS
	void OnPersistentEOSLoginComplete(int32 LocalUserNum, bool bWasSuccessful, const FUniqueNetId& UserId, const FString& Error);

	void OnEOSLoginComplete(int32 LocalUserNum, bool bWasSuccessful, const FUniqueNetId& UserId, const FString& Error);
    
	void HandleSuccessfulEOSLogin(const FUniqueNetId& UserId, int32 LocalUserNum);

	FDelegateHandle PersistentLoginDelegateHandle;
	FDelegateHandle LoginDelegateHandle;
	FString DebugMessage;
};

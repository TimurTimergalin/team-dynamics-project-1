// Fill out your copyright notice in the Description page of Project Settings.

#pragma once

#include "CoreMinimal.h"
#include "GameFramework/GameUserSettings.h"
#include "TagDuelsGameUserSettings.generated.h"

/**
 * 
 */
UCLASS()
class TAGDUELS_API UTagDuelsGameUserSettings : public UGameUserSettings
{
	GENERATED_BODY()

public:
	UFUNCTION(BlueprintPure, Category=Settings)
	float GetMasterVolume();

	UFUNCTION(BlueprintCallable, Category=Settings)
	void SetMasterVolume(float Volume);

	UFUNCTION(BlueprintPure, Category=Settings)
	float GetMusicVolume();

	UFUNCTION(BlueprintCallable, Category=Settings)
	void SetMusicVolume(float Volume);

	UFUNCTION(BlueprintPure, Category=Settings)
	float GetSFXVolume();

	UFUNCTION(BlueprintCallable, Category=Settings)
	void SetSFXVolume(float Volume);
	
private:
	UPROPERTY(config)
	float MasterVolume = 1.0f;
	
	UPROPERTY(config)
	float MusicVolume = 1.0f;

	UPROPERTY(config)
	float SFXVolume = 1.0f;
};

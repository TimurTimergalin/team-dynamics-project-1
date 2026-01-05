// Fill out your copyright notice in the Description page of Project Settings.

#pragma once

#include "CoreMinimal.h"
#include "TagDuelsAttributeSetBase.h"
#include "AbilitySystemComponent.h"
#include "PreyAttributeSet.generated.h"

/**
 * 
 */
UCLASS()
class TAGDUELS_API UPreyAttributeSet : public UTagDuelsAttributeSetBase
{
	GENERATED_BODY()

public:
	UPreyAttributeSet();

	UPROPERTY(BlueprintReadOnly, Category = "Attributes", ReplicatedUsing=OnRep_PalletFastVaultDuration)
	FGameplayAttributeData PalletFastVaultDuration;
	ATTRIBUTE_ACCESSORS_BASIC(UPreyAttributeSet, PalletFastVaultDuration);

	UPROPERTY(BlueprintReadOnly, Category = "Attributes", ReplicatedUsing=OnRep_PalletSlowVaultDuration)
	FGameplayAttributeData PalletSlowVaultDuration;
	ATTRIBUTE_ACCESSORS_BASIC(UPreyAttributeSet, PalletSlowVaultDuration);
	
public:
	UFUNCTION()
	void OnRep_PalletFastVaultDuration (const FGameplayAttributeData& OldValue)
    {
    	GAMEPLAYATTRIBUTE_REPNOTIFY(UPreyAttributeSet, PalletFastVaultDuration, OldValue);
    }

	UFUNCTION()
	void OnRep_PalletSlowVaultDuration (const FGameplayAttributeData& OldValue)
    {
    	GAMEPLAYATTRIBUTE_REPNOTIFY(UPreyAttributeSet, PalletSlowVaultDuration, OldValue);
    }

	virtual void GetLifetimeReplicatedProps(TArray<FLifetimeProperty>& OutLifetimeProps) const override;
};

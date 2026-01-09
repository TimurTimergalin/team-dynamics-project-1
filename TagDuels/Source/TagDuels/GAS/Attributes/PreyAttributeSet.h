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

	UPROPERTY(BlueprintReadOnly, Category = "Attributes", ReplicatedUsing=OnRep_HP)
	FGameplayAttributeData HP;
	ATTRIBUTE_ACCESSORS_BASIC(UPreyAttributeSet, HP);

	UPROPERTY(BlueprintReadOnly, Category = "Attributes", ReplicatedUsing=OnRep_HitHaste)
	FGameplayAttributeData HitHaste;
	ATTRIBUTE_ACCESSORS_BASIC(UPreyAttributeSet, HitHaste);

	UPROPERTY(BlueprintReadOnly, Category = "Attributes", ReplicatedUsing=OnRep_HitHasteDuration)
	FGameplayAttributeData HitHasteDuration;
	ATTRIBUTE_ACCESSORS_BASIC(UPreyAttributeSet, HitHasteDuration);
	
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

	UFUNCTION()
	void OnRep_HP (const FGameplayAttributeData& OldValue)
	{
		GAMEPLAYATTRIBUTE_REPNOTIFY(UPreyAttributeSet, HP, OldValue);
	}

	UFUNCTION()
	void OnRep_HitHaste (const FGameplayAttributeData& OldValue)
	{
		GAMEPLAYATTRIBUTE_REPNOTIFY(UPreyAttributeSet, HitHaste, OldValue);
	}

	UFUNCTION()
	void OnRep_HitHasteDuration (const FGameplayAttributeData& OldValue)
	{
		GAMEPLAYATTRIBUTE_REPNOTIFY(UPreyAttributeSet, HitHasteDuration, OldValue);
	}

	virtual void GetLifetimeReplicatedProps(TArray<FLifetimeProperty>& OutLifetimeProps) const override;
};

// Fill out your copyright notice in the Description page of Project Settings.

#pragma once

#include "CoreMinimal.h"
#include "TagDuelsAttributeSetBase.h"
#include "AbilitySystemComponent.h"
#include "HunterAttributeSet.generated.h"

/**
 * 
 */
UCLASS()
class TAGDUELS_API UHunterAttributeSet : public UTagDuelsAttributeSetBase
{
	GENERATED_BODY()

public:
	UHunterAttributeSet();

	UPROPERTY(BlueprintReadOnly, Category = "Attributes", ReplicatedUsing=OnRep_PalletStunDuration)
	FGameplayAttributeData PalletStunDuration;
	ATTRIBUTE_ACCESSORS_BASIC(UHunterAttributeSet, PalletStunDuration);

	UPROPERTY(BlueprintReadOnly, Category = "Attributes", ReplicatedUsing=OnRep_PalletBreakDuration)
	FGameplayAttributeData PalletBreakDuration;
	ATTRIBUTE_ACCESSORS_BASIC(UHunterAttributeSet, PalletBreakDuration);

	UPROPERTY(BlueprintReadOnly, Category = "Attributes", ReplicatedUsing=OnRep_HunterVaultDuration)
	FGameplayAttributeData HunterVaultDuration;
	ATTRIBUTE_ACCESSORS_BASIC(UHunterAttributeSet, HunterVaultDuration);

	UPROPERTY(BlueprintReadOnly, Category = "Attributes", ReplicatedUsing=OnRep_AttackHaste)
	FGameplayAttributeData AttackHaste;
	ATTRIBUTE_ACCESSORS_BASIC(UHunterAttributeSet, AttackHaste);

	UPROPERTY(BlueprintReadOnly, Category = "Attributes", ReplicatedUsing=OnRep_AttackDuration)
	FGameplayAttributeData AttackDuration;
	ATTRIBUTE_ACCESSORS_BASIC(UHunterAttributeSet, AttackDuration);

	UPROPERTY(BlueprintReadOnly, Category = "Attributes", ReplicatedUsing=OnRep_AttackRecoveryHinder)
	FGameplayAttributeData AttackRecoveryHinder;
	ATTRIBUTE_ACCESSORS_BASIC(UHunterAttributeSet, AttackRecoveryHinder);

	UPROPERTY(BlueprintReadOnly, Category = "Attributes", ReplicatedUsing=OnRep_AttackRecoveryDuration)
	FGameplayAttributeData AttackRecoveryDuration;
	ATTRIBUTE_ACCESSORS_BASIC(UHunterAttributeSet, AttackRecoveryDuration);

	UPROPERTY(BlueprintReadOnly, Category = "Attributes", ReplicatedUsing=OnRep_AttackSuccessRecoveryDuration)
	FGameplayAttributeData AttackSuccessRecoveryDuration;
	ATTRIBUTE_ACCESSORS_BASIC(UHunterAttributeSet, AttackSuccessRecoveryDuration);

	UPROPERTY(BlueprintReadOnly, Category = "Attributes", ReplicatedUsing=OnRep_BloodlustLevel)
	FGameplayAttributeData BloodlustLevel;
	ATTRIBUTE_ACCESSORS_BASIC(UHunterAttributeSet, BloodlustLevel);
	
public:	
	UFUNCTION()
	void OnRep_PalletStunDuration (const FGameplayAttributeData& OldValue)
	{
		GAMEPLAYATTRIBUTE_REPNOTIFY(UHunterAttributeSet, PalletStunDuration, OldValue);
	}

	UFUNCTION()
	void OnRep_PalletBreakDuration (const FGameplayAttributeData& OldValue)
	{
		GAMEPLAYATTRIBUTE_REPNOTIFY(UHunterAttributeSet, PalletBreakDuration, OldValue);
	}

	UFUNCTION()
	void OnRep_HunterVaultDuration (const FGameplayAttributeData& OldValue)
	{
		GAMEPLAYATTRIBUTE_REPNOTIFY(UHunterAttributeSet, HunterVaultDuration, OldValue);
	}

	UFUNCTION()
	void OnRep_AttackHaste (const FGameplayAttributeData& OldValue)
	{
		GAMEPLAYATTRIBUTE_REPNOTIFY(UHunterAttributeSet, AttackHaste, OldValue);
	}

	UFUNCTION()
	void OnRep_AttackDuration (const FGameplayAttributeData& OldValue)
	{
		GAMEPLAYATTRIBUTE_REPNOTIFY(UHunterAttributeSet, AttackDuration, OldValue);
	}

	UFUNCTION()
	void OnRep_AttackRecoveryHinder (const FGameplayAttributeData& OldValue)
	{
		GAMEPLAYATTRIBUTE_REPNOTIFY(UHunterAttributeSet, AttackRecoveryHinder, OldValue);
	}

	UFUNCTION()
	void OnRep_AttackRecoveryDuration (const FGameplayAttributeData& OldValue)
	{
		GAMEPLAYATTRIBUTE_REPNOTIFY(UHunterAttributeSet, AttackRecoveryDuration, OldValue);
	}

	UFUNCTION()
	void OnRep_AttackSuccessRecoveryDuration (const FGameplayAttributeData& OldValue)
	{
		GAMEPLAYATTRIBUTE_REPNOTIFY(UHunterAttributeSet, AttackSuccessRecoveryDuration, OldValue);
	}

	UFUNCTION()
	void OnRep_BloodlustLevel (const FGameplayAttributeData& OldValue)
	{
		GAMEPLAYATTRIBUTE_REPNOTIFY(UHunterAttributeSet, BloodlustLevel, OldValue);
	}

	virtual void GetLifetimeReplicatedProps(TArray<FLifetimeProperty>& OutLifetimeProps) const override;
};

// Fill out your copyright notice in the Description page of Project Settings.

#pragma once

#include "CoreMinimal.h"
#include "AbilitySystemComponent.h"
#include "AttributeSet.h"
#include "BattleAttributeSet.generated.h"

/**
 * 
 */
UCLASS()
class NEXUS_API UBattleAttributeSet : public UAttributeSet
{
	GENERATED_BODY()
public:
	UBattleAttributeSet();
	
	// Battle Attributes
	UPROPERTY(BlueprintReadOnly, Category = "Attributes", ReplicatedUsing=OnRep_Haste)
	FGameplayAttributeData Haste;
	ATTRIBUTE_ACCESSORS_BASIC(UBattleAttributeSet, Haste);
	
	UPROPERTY(BlueprintReadOnly, Category = "Attributes", ReplicatedUsing=OnRep_Damage)
	FGameplayAttributeData Damage;
	ATTRIBUTE_ACCESSORS_BASIC(UBattleAttributeSet, Damage);

	UPROPERTY(BlueprintReadOnly, Category = "Attributes", ReplicatedUsing=OnRep_HealingAmount)
	FGameplayAttributeData HealingAmount;
	ATTRIBUTE_ACCESSORS_BASIC(UBattleAttributeSet, HealingAmount);

	UPROPERTY(BlueprintReadOnly, Category = "Attributes", ReplicatedUsing=OnRep_HealingTime)
	FGameplayAttributeData HealingTime;
	ATTRIBUTE_ACCESSORS_BASIC(UBattleAttributeSet, HealingTime);

public:
	UFUNCTION()
	void OnRep_Haste (const FGameplayAttributeData& OldValue)
	{
		GAMEPLAYATTRIBUTE_REPNOTIFY(UBattleAttributeSet, Haste, OldValue);
	}
	
	UFUNCTION()
	void OnRep_Damage (const FGameplayAttributeData& OldValue)
	{
		GAMEPLAYATTRIBUTE_REPNOTIFY(UBattleAttributeSet, Damage, OldValue);
	}

	UFUNCTION()
	void OnRep_HealingAmount (const FGameplayAttributeData& OldValue)
	{
		GAMEPLAYATTRIBUTE_REPNOTIFY(UBattleAttributeSet, HealingAmount, OldValue);
	}

	UFUNCTION()
	void OnRep_HealingTime (const FGameplayAttributeData& OldValue)
	{
		GAMEPLAYATTRIBUTE_REPNOTIFY(UBattleAttributeSet, HealingTime, OldValue);
	}

	virtual void GetLifetimeReplicatedProps(TArray<FLifetimeProperty>& OutLifetimeProps) const override;
};

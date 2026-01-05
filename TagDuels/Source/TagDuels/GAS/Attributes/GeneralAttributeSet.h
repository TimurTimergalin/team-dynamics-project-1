// Fill out your copyright notice in the Description page of Project Settings.

#pragma once

#include "CoreMinimal.h"
#include "TagDuelsAttributeSetBase.h"
#include "AbilitySystemComponent.h"
#include "GeneralAttributeSet.generated.h"

/**
 * 
 */
UCLASS()
class TAGDUELS_API UGeneralAttributeSet : public UTagDuelsAttributeSetBase
{
	GENERATED_BODY()
	
public:
	UGeneralAttributeSet();

	UPROPERTY(BlueprintReadOnly, Category = "Attributes", ReplicatedUsing=OnRep_Haste)
	FGameplayAttributeData Haste;
	ATTRIBUTE_ACCESSORS_BASIC(UGeneralAttributeSet, Haste);
	
	UPROPERTY(BlueprintReadOnly, Category = "Attributes", ReplicatedUsing=OnRep_Hinder)
	FGameplayAttributeData Hinder;
	ATTRIBUTE_ACCESSORS_BASIC(UGeneralAttributeSet, Hinder);
	
public:
	UFUNCTION()
	void OnRep_Haste (const FGameplayAttributeData& OldValue)
	{
		GAMEPLAYATTRIBUTE_REPNOTIFY(UGeneralAttributeSet, Haste, OldValue);
	}
	
	UFUNCTION()
	void OnRep_Hinder (const FGameplayAttributeData& OldValue)
	{
		GAMEPLAYATTRIBUTE_REPNOTIFY(UGeneralAttributeSet, Hinder, OldValue);
	}

	virtual void GetLifetimeReplicatedProps(TArray<FLifetimeProperty>& OutLifetimeProps) const override;
};


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
	
	UPROPERTY(BlueprintReadOnly, Category = "Attributes", ReplicatedUsing=OnRep_RageLevel)
	FGameplayAttributeData RageLevel;
	ATTRIBUTE_ACCESSORS_BASIC(UHunterAttributeSet, RageLevel);
	
public:
	UFUNCTION()
	void OnRep_RageLevel (const FGameplayAttributeData& OldValue)
	{
		GAMEPLAYATTRIBUTE_REPNOTIFY(UHunterAttributeSet, RageLevel, OldValue);
	}

	virtual void GetLifetimeReplicatedProps(TArray<FLifetimeProperty>& OutLifetimeProps) const override;
};

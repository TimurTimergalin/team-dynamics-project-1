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
	
	UPROPERTY(BlueprintReadOnly, Category = "Attributes", ReplicatedUsing=OnRep_General)
	FGameplayAttributeData General;
	ATTRIBUTE_ACCESSORS_BASIC(UGeneralAttributeSet, General);
	
public:
	UFUNCTION()
	void OnRep_General (const FGameplayAttributeData& OldValue)
	{
		GAMEPLAYATTRIBUTE_REPNOTIFY(UGeneralAttributeSet, General, OldValue);
	}

	virtual void GetLifetimeReplicatedProps(TArray<FLifetimeProperty>& OutLifetimeProps) const override;
};


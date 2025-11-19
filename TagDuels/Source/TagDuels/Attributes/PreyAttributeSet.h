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
	
	UPROPERTY(BlueprintReadOnly, Category = "Attributes", ReplicatedUsing=OnRep_Haste)
	FGameplayAttributeData Haste;
	ATTRIBUTE_ACCESSORS_BASIC(UPreyAttributeSet, Haste);
	
public:
    UFUNCTION()
    void OnRep_Haste (const FGameplayAttributeData& OldValue)
    {
    	GAMEPLAYATTRIBUTE_REPNOTIFY(UPreyAttributeSet, Haste, OldValue);
    }

	virtual void GetLifetimeReplicatedProps(TArray<FLifetimeProperty>& OutLifetimeProps) const override;
};

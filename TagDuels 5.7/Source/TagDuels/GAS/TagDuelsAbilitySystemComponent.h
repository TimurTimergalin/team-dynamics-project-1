// Fill out your copyright notice in the Description page of Project Settings.

#pragma once

#include "CoreMinimal.h"
#include "AbilitySystemComponent.h"
#include "TagDuelsAbilitySystemComponent.generated.h"


UCLASS(ClassGroup=(Custom), meta=(BlueprintSpawnableComponent))
class TAGDUELS_API UTagDuelsAbilitySystemComponent : public UAbilitySystemComponent
{
	GENERATED_BODY()

public:
	// Sets default values for this component's properties
	UTagDuelsAbilitySystemComponent();
	
	// /** Cancels the specified ability CDO. */
	// UFUNCTION(BlueprintCallable, Category="AbilitySystem")
	// void CancelAbility(UGameplayAbility* Ability);	
	//
	// /** Cancels the ability indicated by passed in spec handle. If handle is not found among reactivated abilities nothing happens. */
	// UFUNCTION(BlueprintCallable, Category="AbilitySystem")
	// void CancelAbilityHandle(const FGameplayAbilitySpecHandle& AbilityHandle);
	//
	// /** Cancel all abilities with the specified tags. Will not cancel the Ignore instance */
	// UFUNCTION(BlueprintCallable, Category="AbilitySystem")
	// void CancelAbilities(const FGameplayTagContainer* WithTags=nullptr, const FGameplayTagContainer* WithoutTags=nullptr, UGameplayAbility* Ignore=nullptr);
	//
	// /** Cancels all abilities regardless of tags. Will not cancel the ignore instance */
	// UFUNCTION(BlueprintCallable, Category="AbilitySystem")
	// void CancelAllAbilities(UGameplayAbility* Ignore=nullptr);

protected:
	// Called when the game starts
	virtual void BeginPlay() override;

public:
	// Called every frame
	virtual void TickComponent(float DeltaTime, ELevelTick TickType,
	                           FActorComponentTickFunction* ThisTickFunction) override;
};

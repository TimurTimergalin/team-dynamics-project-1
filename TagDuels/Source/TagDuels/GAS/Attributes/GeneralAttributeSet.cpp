// Fill out your copyright notice in the Description page of Project Settings.


#include "GeneralAttributeSet.h"
#include "Net/UnrealNetwork.h"
#include "GameplayEffectExtension.h"

UGeneralAttributeSet::UGeneralAttributeSet()
{
	Haste = 1.f;
	Hinder = 1.f;
}

void UGeneralAttributeSet::GetLifetimeReplicatedProps(TArray<FLifetimeProperty>& OutLifetimeProps) const
{
	Super::GetLifetimeReplicatedProps(OutLifetimeProps);
	DOREPLIFETIME_CONDITION_NOTIFY(UGeneralAttributeSet, Haste, COND_None, REPNOTIFY_Always);
	DOREPLIFETIME_CONDITION_NOTIFY(UGeneralAttributeSet, Hinder, COND_None, REPNOTIFY_Always);
}

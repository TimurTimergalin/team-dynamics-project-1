// Fill out your copyright notice in the Description page of Project Settings.


#include "HunterAttributeSet.h"
#include "Net/UnrealNetwork.h"
#include "GameplayEffectExtension.h"

UHunterAttributeSet::UHunterAttributeSet()
{
	RageLevel = 1.f;
	PalletStunDuration = 2.f;
	PalletBreakDuration = 2.34f;
}

void UHunterAttributeSet::GetLifetimeReplicatedProps(TArray<FLifetimeProperty>& OutLifetimeProps) const
{
	Super::GetLifetimeReplicatedProps(OutLifetimeProps);
	DOREPLIFETIME_CONDITION_NOTIFY(UHunterAttributeSet, RageLevel, COND_None, REPNOTIFY_Always);
	DOREPLIFETIME_CONDITION_NOTIFY(UHunterAttributeSet, PalletStunDuration, COND_None, REPNOTIFY_Always);
	DOREPLIFETIME_CONDITION_NOTIFY(UHunterAttributeSet, PalletBreakDuration, COND_None, REPNOTIFY_Always);
}

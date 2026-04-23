// Fill out your copyright notice in the Description page of Project Settings.


#include "PreyAttributeSet.h"
#include "Net/UnrealNetwork.h"
#include "GameplayEffectExtension.h"

UPreyAttributeSet::UPreyAttributeSet()
{
	PalletFastVaultDuration = 1.f; //og 1.1
	PalletSlowVaultDuration = 1.9f; //og = 2
	WindowFastVaultDuration = 0.5f;
	WindowMediumVaultDuration = 0.9f;
	WindowSlowVaultDuration = 1.5f;
	HP = 2.f;
	HitHaste = 1.5f;
	HitHasteDuration = 2.f;
}

void UPreyAttributeSet::GetLifetimeReplicatedProps(TArray<FLifetimeProperty>& OutLifetimeProps) const
{
	Super::GetLifetimeReplicatedProps(OutLifetimeProps);
	DOREPLIFETIME_CONDITION_NOTIFY(UPreyAttributeSet, PalletFastVaultDuration, COND_None, REPNOTIFY_Always);
	DOREPLIFETIME_CONDITION_NOTIFY(UPreyAttributeSet, PalletSlowVaultDuration, COND_None, REPNOTIFY_Always);
	DOREPLIFETIME_CONDITION_NOTIFY(UPreyAttributeSet, WindowFastVaultDuration, COND_None, REPNOTIFY_Always);
	DOREPLIFETIME_CONDITION_NOTIFY(UPreyAttributeSet, WindowMediumVaultDuration, COND_None, REPNOTIFY_Always);
	DOREPLIFETIME_CONDITION_NOTIFY(UPreyAttributeSet, WindowSlowVaultDuration, COND_None, REPNOTIFY_Always);
	DOREPLIFETIME_CONDITION_NOTIFY(UPreyAttributeSet, HP, COND_None, REPNOTIFY_Always);
	DOREPLIFETIME_CONDITION_NOTIFY(UPreyAttributeSet, HitHaste, COND_None, REPNOTIFY_Always);
	DOREPLIFETIME_CONDITION_NOTIFY(UPreyAttributeSet, HitHasteDuration, COND_None, REPNOTIFY_Always);
}

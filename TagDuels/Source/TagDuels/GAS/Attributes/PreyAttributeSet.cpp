// Fill out your copyright notice in the Description page of Project Settings.


#include "PreyAttributeSet.h"
#include "Net/UnrealNetwork.h"
#include "GameplayEffectExtension.h"

UPreyAttributeSet::UPreyAttributeSet()
{
	Haste = 100.f;
}

void UPreyAttributeSet::GetLifetimeReplicatedProps(TArray<FLifetimeProperty>& OutLifetimeProps) const
{
	Super::GetLifetimeReplicatedProps(OutLifetimeProps);
	DOREPLIFETIME_CONDITION_NOTIFY(UPreyAttributeSet, Haste, COND_None, REPNOTIFY_Always);
}

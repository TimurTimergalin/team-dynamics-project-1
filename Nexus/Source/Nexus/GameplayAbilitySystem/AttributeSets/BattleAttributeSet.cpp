// Fill out your copyright notice in the Description page of Project Settings.


#include "BattleAttributeSet.h"
#include "Net/UnrealNetwork.h"

UBattleAttributeSet::UBattleAttributeSet()
{
	Damage = 20.f;
	HealingAmount = 20.f;
	HealingTime = 2.f;
}

void UBattleAttributeSet::GetLifetimeReplicatedProps(TArray<FLifetimeProperty>& OutLifetimeProps) const
{
	Super::GetLifetimeReplicatedProps(OutLifetimeProps);
	DOREPLIFETIME_CONDITION_NOTIFY(UBattleAttributeSet, Damage, COND_None, REPNOTIFY_Always);
	DOREPLIFETIME_CONDITION_NOTIFY(UBattleAttributeSet, HealingAmount, COND_None, REPNOTIFY_Always);
	DOREPLIFETIME_CONDITION_NOTIFY(UBattleAttributeSet, HealingTime, COND_None, REPNOTIFY_Always);
}

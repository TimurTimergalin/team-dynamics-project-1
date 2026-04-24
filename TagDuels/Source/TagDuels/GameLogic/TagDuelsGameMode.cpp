// Fill out your copyright notice in the Description page of Project Settings.


#include "TagDuelsGameMode.h"

void ATagDuelsGameMode::PreLogin(const FString& Options, const FString& Address, const FUniqueNetIdRepl& UniqueId,
	FString& ErrorMessage)
{
	Super::PreLogin(Options, Address, UniqueId, ErrorMessage);
	OnPreLogin_Implementation(Options, Address, UniqueId, ErrorMessage);
}

void ATagDuelsGameMode::OnPreLogin_Implementation(const FString& Options, const FString& Address,
	const FUniqueNetIdRepl& UniqueId, FString& ErrorMessage)
{
}

// Fill out your copyright notice in the Description page of Project Settings.


#include "TagDuelsGameMode.h"

#include "GameFramework/GameSession.h"
#include "Net/OnlineEngineInterface.h"

void ATagDuelsGameMode::PreLogin(const FString& Options, const FString& Address, const FUniqueNetIdRepl& UniqueId,
                                 FString& ErrorMessage)
{
	// // Login unique id must match server expected unique id type OR No unique id could mean game doesn't use them
	// const bool bUniqueIdCheckOk = (!UniqueId.IsValid() || UOnlineEngineInterface::Get()->IsCompatibleUniqueNetId(UniqueId));
	// if (bUniqueIdCheckOk)
	// {
	// 	ErrorMessage = GameSession->ApproveLogin(Options);
	// }
	// else
	// {
	// 	ErrorMessage = TEXT("incompatible_unique_net_id");
	// }
	//
	// FGameModeEvents::GameModePreLoginEvent.Broadcast(this, UniqueId, ErrorMessage);
	Super::PreLogin(Options, Address, UniqueId, ErrorMessage);
	FString DebugMessage = FString::Printf(TEXT("Prelogin C++"));
	UE_LOG(LogTemp, Error, TEXT("%s"), *DebugMessage);
	if (GEngine) GEngine->AddOnScreenDebugMessage(-1,10.0f, FColor::Green, DebugMessage);
	OnPreLogin_Implementation(Options, Address, UniqueId);
}

void ATagDuelsGameMode::OnPreLogin_Implementation(const FString& Options, const FString& Address,
	const FUniqueNetIdRepl& UniqueId)
{
}

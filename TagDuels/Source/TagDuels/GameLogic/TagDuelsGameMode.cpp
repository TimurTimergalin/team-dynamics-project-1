// Fill out your copyright notice in the Description page of Project Settings.


#include "TagDuelsGameMode.h"
#include "GameFramework/GameSession.h"
#include "Net/OnlineEngineInterface.h"

void ATagDuelsGameMode::PreLogin(const FString& Options, const FString& Address, const FUniqueNetIdRepl& UniqueId,
                                 FString& ErrorMessage)
{
	if (UniqueId.IsValid())
	{
		if (UniqueId.IsV1())
		{
			ErrorMessage = GameSession->ApproveLogin(Options);
		}
		else
		{
			ErrorMessage = TEXT("incompatible_unique_net_id");
		}
	
		if (ErrorMessage.IsEmpty())
		{
			OnPreLogin(Options, Address, UniqueId, ErrorMessage);
		}
	}
	else
	{
		ErrorMessage = TEXT("invalid_unique_net_id");
	}
	UE_LOG(LogTemp, Error, TEXT("%s"), *ErrorMessage);
	if (GEngine) GEngine->AddOnScreenDebugMessage(-1,10.0f, FColor::Red, ErrorMessage);
	
	FGameModeEvents::GameModePreLoginEvent.Broadcast(this, UniqueId, ErrorMessage);
}

void ATagDuelsGameMode::OnPreLogin_Implementation(const FString& Options, const FString& Address,
	const FUniqueNetIdRepl& UniqueId, FString& ErrorMessage)
{
}

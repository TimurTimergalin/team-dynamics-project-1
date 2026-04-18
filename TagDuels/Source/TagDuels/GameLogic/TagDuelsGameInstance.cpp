// Fill out your copyright notice in the Description page of Project Settings.


#include "TagDuelsGameInstance.h"
#include "OnlineSubsystem.h"
#include "Interfaces/OnlineIdentityInterface.h"
#include "OnlineSubsystemUtils.h"

FString UTagDuelsGameInstance::GetSteamAuthToken()
{
	// 1. Get the Subsystem (it will auto-detect Steam if configured in .ini)
	IOnlineSubsystem* Subsystem = Online::GetSubsystem(GetWorld());
	if (!Subsystem) return FString();

	// 2. Get the Identity Interface
	IOnlineIdentityPtr Identity = Subsystem->GetIdentityInterface();
	if (!Identity.IsValid()) return FString();

	// 3. Get the Token for the local player (usually index 0)
	// For Steam, this string is the Hex-encoded Auth Session Ticket
	FString AuthToken = Identity->GetAuthToken(0);

	if (AuthToken.IsEmpty())
	{
		UE_LOG(LogTemp, Warning, TEXT("Failed to retrieve Auth Token. Is Steam running?"));
	}

	return AuthToken;
}

FString UTagDuelsGameInstance::GetSteamID()
{
	IOnlineSubsystem* Subsystem = Online::GetSubsystem(GetWorld());
	if (Subsystem)
	{
		IOnlineIdentityPtr Identity = Subsystem->GetIdentityInterface();
		if (Identity.IsValid())
		{
			// Get the Unique ID pointer for the local player (index 0)
			FUniqueNetIdPtr UserId = Identity->GetUniquePlayerId(0);

			if (UserId.IsValid())
			{
				// ToString() returns the SteamID64 as a string
				return UserId->ToString();
			}
		}
	}
	return FString(TEXT("InvalidID"));
}
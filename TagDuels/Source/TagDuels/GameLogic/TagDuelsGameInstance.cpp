// Fill out your copyright notice in the Description page of Project Settings.


#include "TagDuelsGameInstance.h"
#include "OnlineSubsystem.h"
#include "Interfaces/OnlineIdentityInterface.h"
#include "OnlineSubsystemUtils.h"

// Base Methods
void UTagDuelsGameInstance::OnStart()
{
	Super::OnStart();
	BlStart();
}
void UTagDuelsGameInstance::BlStart_Implementation()
{
}

// Steam OSS
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

// EpicGames OSS
void UTagDuelsGameInstance::LoginToEOS()
{
    IOnlineSubsystem* OnlineSubsystem = Online::GetSubsystem(GetWorld());
    if (!OnlineSubsystem)
    {
        UE_LOG(LogTemp, Error, TEXT("OnlineSubsystem not available"));
        return;
    }

    IOnlineIdentityPtr IdentityInterface = OnlineSubsystem->GetIdentityInterface();
    if (!IdentityInterface.IsValid())
    {
        UE_LOG(LogTemp, Error, TEXT("Identity interface not available"));
        return;
    }
	
	UE_LOG(LogTemp, Display, TEXT("AccountPortal Start"));
	// AccountPortal automatically launches the browser for user login
	FOnlineAccountCredentials Credentials;
	Credentials.Type = FString("accountportal");
	Credentials.Id = FString();
	Credentials.Token = FString();
    
    IdentityInterface->AddOnLoginCompleteDelegate_Handle(0, 
        FOnLoginCompleteDelegate::CreateUObject(this, &UTagDuelsGameInstance::OnLoginComplete));
	UE_LOG(LogTemp, Display, TEXT("Login"));
    IdentityInterface->Login(0, Credentials);
}

void UTagDuelsGameInstance::OnLoginComplete(int32 LocalUserNum, bool bWasSuccessful, const FUniqueNetId& UserId, const FString& Error)
{
	if (bWasSuccessful)
	{
		// Silent login succeeded!
		UE_LOG(LogTemp, Display, TEXT("Persistent Auth successful for Epic Account ID: %s"), *UserId.ToString());

		// Retrieve Auth Token and send to backend
		IOnlineSubsystem* OnlineSubsystem = Online::GetSubsystem(GetWorld());
		IOnlineIdentityPtr IdentityInterface = OnlineSubsystem->GetIdentityInterface();

		FString EpicAccountId = UserId.ToString();
		FString AuthToken = IdentityInterface->GetAuthToken(LocalUserNum);

		UE_LOG(LogTemp, Display, TEXT("Ready to send to backend - Epic ID: %s, Token: %s"), *EpicAccountId, *AuthToken);
		// In Game Print String
		if (GEngine)
		{
			FString DebugMessage = FString::Printf(TEXT("Ready to send to backend - Epic ID: %s, Token: %s"), *EpicAccountId, *AuthToken);
    
			GEngine->AddOnScreenDebugMessage(
				-1,                      // Unique key (-1 means it won't overwrite previous messages)
				10.0f,                   // Duration in seconds
				FColor::Green,           // Text color
				DebugMessage             // The actual string
			);
		}
	}
	else
	{
		UE_LOG(LogTemp, Warning, TEXT("Persistent Auth failed. Error: %s"), *Error);
	}
}

// Online
EOnlineSubsystemType UTagDuelsGameInstance::GetActiveOnlineSubsystemType() const
{
	FName SubsystemName = GetOnlinePlatformName(); // "Steam", "EOS", "NULL" и т.д.

	// IOnlineSubsystem* EpicOSS = IOnlineSubsystem::Get(TEXT("EOS"));
	//
	// IOnlineSubsystem* SteamOSS = IOnlineSubsystem::Get(TEXT("Steam"));

	if (SubsystemName == FName(TEXT("Steam")))
	{
		return EOnlineSubsystemType::Steam;
	}
	if (SubsystemName == FName(TEXT("EOS")))
	{
		return EOnlineSubsystemType::EOS;
	}
	if (SubsystemName == FName(TEXT("NULL")))
	{
		return EOnlineSubsystemType::Null;
	}
	if (SubsystemName == NAME_None)
	{
		return EOnlineSubsystemType::None;
	}
    
	// На случай неизвестных подсистем
	UE_LOG(LogTemp, Warning, TEXT("Unknown online subsystem: %s"), *SubsystemName.ToString());
	return EOnlineSubsystemType::None;
}





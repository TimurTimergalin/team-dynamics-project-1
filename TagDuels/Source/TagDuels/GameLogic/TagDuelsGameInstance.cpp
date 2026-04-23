// Fill out your copyright notice in the Description page of Project Settings.


#include "TagDuelsGameInstance.h"
#include "OnlineSubsystem.h"
#include "Interfaces/OnlineIdentityInterface.h"
#include "OnlineSubsystemUtils.h"
#include "TagDuels/OnlineSubsystem/utils.h"


// Basic Methods
void UTagDuelsGameInstance::OnStart()
{
	Super::OnStart();
	Start();
}
void UTagDuelsGameInstance::Start_Implementation()
{
}

// Steam
void UTagDuelsGameInstance::LoginToSteam()
{
	IOnlineSubsystem* OnlineSubsystem = Online::GetSubsystem(GetWorld());
	if (!OnlineSubsystem)
	{
		DebugMessage = FString::Printf(TEXT("OnlineSubsystem not available"));
		UE_LOG(LogTemp, Error, TEXT("%s"), *DebugMessage);
		OnFailedToLogin(DebugMessage);
		return;
	}
	
	IOnlineIdentityPtr IdentityInterface = OnlineSubsystem->GetIdentityInterface();
	if (!IdentityInterface.IsValid())
	{
		DebugMessage = FString::Printf(TEXT("Identity interface not available"));
		UE_LOG(LogTemp, Error, TEXT("%s"), *DebugMessage);
		OnFailedToLogin(DebugMessage);
		return;
	}
	
	FUniqueNetIdPtr UserId = IdentityInterface->GetUniquePlayerId(0);
	if (!UserId.IsValid())
	{
		DebugMessage = FString::Printf(TEXT("UserId not available"));
		UE_LOG(LogTemp, Error, TEXT("%s"), *DebugMessage);
		OnFailedToLogin(DebugMessage);
		return;
	}
	
	// Get ID and Auth Token
	int64 SteamId64 = StrToInt64(UserId->ToString()).Get(0);
	if (SteamId64 == 0)
	{
		DebugMessage = FString::Printf(TEXT("SteamId not available"));
		UE_LOG(LogTemp, Error, TEXT("%s"), *DebugMessage);
		OnFailedToLogin(DebugMessage);
		return;
	}
	FString AuthToken = IdentityInterface->GetAuthToken(0);
	
	// Print ID and Auth Token
	DebugMessage = FString::Printf(TEXT("Ready to send to backend - Steam ID64: %lld, Token: %s"), SteamId64, *AuthToken);
	UE_LOG(LogTemp, Display, TEXT("%s"), *DebugMessage);
	if (GEngine) GEngine->AddOnScreenDebugMessage(-1,10.0f, FColor::Green, DebugMessage);
	
	// Вызываем ивент в Blueprint GameInstance
	OnSuccessfulLoginSteam(SteamId64, AuthToken);
}

// EpicGames
void UTagDuelsGameInstance::LoginToEOS()
{
    IOnlineSubsystem* OnlineSubsystem = Online::GetSubsystem(GetWorld());
    if (!OnlineSubsystem)
    {
    	DebugMessage = FString::Printf(TEXT("OnlineSubsystem not available"));
    	UE_LOG(LogTemp, Display, TEXT("%s"), *DebugMessage);
    	OnFailedToLogin(DebugMessage);
        return;
    }

    IOnlineIdentityPtr IdentityInterface = OnlineSubsystem->GetIdentityInterface();
    if (!IdentityInterface.IsValid())
    {
    	DebugMessage = FString::Printf(TEXT("Identity interface not available"));
    	UE_LOG(LogTemp, Display, TEXT("%s"), *DebugMessage);
    	OnFailedToLogin(DebugMessage);
        return;
    }

    UE_LOG(LogTemp, Display, TEXT("Attempting Persistent Auth"));

    FOnlineAccountCredentials Credentials;
    Credentials.Type = FString("persistentauth");

    PersistentLoginDelegateHandle = IdentityInterface->AddOnLoginCompleteDelegate_Handle(0,
        FOnLoginCompleteDelegate::CreateUObject(this, &UTagDuelsGameInstance::OnPersistentEOSLoginComplete));

    if (!IdentityInterface->Login(0, Credentials))
    {
    	DebugMessage = FString::Printf(TEXT("Persistent Auth Login() call failed immediately"));
    	UE_LOG(LogTemp, Display, TEXT("%s"), *DebugMessage);
    	OnFailedToLogin(DebugMessage);
        IdentityInterface->ClearOnLoginCompleteDelegate_Handle(0, PersistentLoginDelegateHandle);
        PersistentLoginDelegateHandle.Reset();
    }
}

void UTagDuelsGameInstance::OnPersistentEOSLoginComplete(int32 LocalUserNum, bool bWasSuccessful, const FUniqueNetId& UserId, const FString& Error)
{
    UE_LOG(LogTemp, Warning, TEXT("Persistent Auth Result:"));
    UE_LOG(LogTemp, Warning, TEXT("LocalUserNum: %d, Success: %d, Error: %s"), LocalUserNum, bWasSuccessful, *Error);

    IOnlineSubsystem* OnlineSubsystem = Online::GetSubsystem(GetWorld());
    if (OnlineSubsystem)
    {
        IOnlineIdentityPtr IdentityInterface = OnlineSubsystem->GetIdentityInterface();
        if (IdentityInterface.IsValid())
        {
            IdentityInterface->ClearOnLoginCompleteDelegate_Handle(LocalUserNum, PersistentLoginDelegateHandle);
            PersistentLoginDelegateHandle.Reset();
        }
    }

    if (bWasSuccessful)
    {
        // Автоматический логин сработал
        UE_LOG(LogTemp, Display, TEXT("Silent login successful!"));
        HandleSuccessfulEOSLogin(UserId, LocalUserNum);
        return;
    }

    // Persistent Auth не сработал, запускаем AccountPortal
    UE_LOG(LogTemp, Warning, TEXT("Persistent Auth failed, falling back to AccountPortal. Error: %s"), *Error);

	IOnlineIdentityPtr IdentityInterface = OnlineSubsystem->GetIdentityInterface();
	if (!IdentityInterface.IsValid())
	{
		DebugMessage = FString::Printf(TEXT("Identity interface not available"));
		UE_LOG(LogTemp, Display, TEXT("%s"), *DebugMessage);
		OnFailedToLogin(DebugMessage);
		return;
	}

    FOnlineAccountCredentials PortalCredentials;
    PortalCredentials.Type = FString("accountportal");

    LoginDelegateHandle = IdentityInterface->AddOnLoginCompleteDelegate_Handle(0,
        FOnLoginCompleteDelegate::CreateUObject(this, &UTagDuelsGameInstance::OnEOSLoginComplete));

    if (!IdentityInterface->Login(0, PortalCredentials))
    {
    	DebugMessage = FString::Printf(TEXT("Login() call failed immediately"));
    	UE_LOG(LogTemp, Display, TEXT("%s"), *DebugMessage);
    	OnFailedToLogin(DebugMessage);
        IdentityInterface->ClearOnLoginCompleteDelegate_Handle(0, LoginDelegateHandle);
        LoginDelegateHandle.Reset();
    }
    else
    {
        UE_LOG(LogTemp, Display, TEXT("Login launched"));
    }
}

void UTagDuelsGameInstance::OnEOSLoginComplete(int32 LocalUserNum, bool bWasSuccessful, const FUniqueNetId& UserId, const FString& Error)
{
	UE_LOG(LogTemp, Warning, TEXT("Portal Auth Result:"));
    UE_LOG(LogTemp, Warning, TEXT("LocalUserNum: %d, Success: %d, Error: %s"), LocalUserNum, bWasSuccessful, *Error);

    IOnlineSubsystem* OnlineSubsystem = Online::GetSubsystem(GetWorld());
    if (OnlineSubsystem)
    {
        IOnlineIdentityPtr IdentityInterface = OnlineSubsystem->GetIdentityInterface();
        if (IdentityInterface.IsValid())
        {
            IdentityInterface->ClearOnLoginCompleteDelegate_Handle(LocalUserNum, LoginDelegateHandle);
            LoginDelegateHandle.Reset();
        }
    }

    if (bWasSuccessful)
    {
        UE_LOG(LogTemp, Display, TEXT("Login successful!"));
        HandleSuccessfulEOSLogin(UserId, LocalUserNum);
    }
    else
    {
        UE_LOG(LogTemp, Error, TEXT("Login failed: %s"), *Error);
    }
}

void UTagDuelsGameInstance::HandleSuccessfulEOSLogin(const FUniqueNetId& UserId, int32 LocalUserNum)
{
	IOnlineSubsystem* OnlineSubsystem = Online::GetSubsystem(GetWorld());
	if (!OnlineSubsystem)
	{
		DebugMessage = FString::Printf(TEXT("OnlineSubsystem not available"));
		UE_LOG(LogTemp, Display, TEXT("%s"), *DebugMessage);
		OnFailedToLogin(DebugMessage);
		return;
	}

	IOnlineIdentityPtr IdentityInterface = OnlineSubsystem->GetIdentityInterface();
	if (!IdentityInterface.IsValid())
	{
		DebugMessage = FString::Printf(TEXT("Identity interface not available"));
		UE_LOG(LogTemp, Display, TEXT("%s"), *DebugMessage);
		OnFailedToLogin(DebugMessage);
		return;
	}


	FString AccountId = UserId.ToString();
	FString AuthToken = IdentityInterface->GetAuthToken(LocalUserNum);

	// Print ID and Auth Token
	DebugMessage = FString::Printf(TEXT("Ready to send to backend - Epic ID: %s, Token: %s"), *AccountId, *AuthToken);
	UE_LOG(LogTemp, Display, TEXT("%s"), *DebugMessage);
	if (GEngine) GEngine->AddOnScreenDebugMessage(-1,10.0f, FColor::Green, DebugMessage);

	// Вызываем ивент в Blueprint GameInstance
	OnSuccessfulLoginEOS(AccountId, AuthToken);
}

// OSS
EOnlineSubsystemType UTagDuelsGameInstance::GetActiveOnlineSubsystemType() const
{
	FName SubsystemName = GetOnlinePlatformName();

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

// Fill out your copyright notice in the Description page of Project Settings.


#include "TagDuelsGameInstance.h"
#include "OnlineSubsystem.h"
#include "Interfaces/OnlineIdentityInterface.h"
#include "OnlineSubsystemUtils.h"
// #include "EOSSubsystem.h"
// #include "OnlineSubsystemEOS.h"
// #include "eos_auth.h"

// Basic Methods
void UTagDuelsGameInstance::OnStart()
{
	Super::OnStart();
	BlStart();
}
void UTagDuelsGameInstance::BlStart_Implementation()
{
}

// OSS
void UTagDuelsGameInstance::LoginToSteam()
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

	FUniqueNetIdPtr UserId = IdentityInterface->GetUniquePlayerId(0);
	if (!UserId.IsValid())
	{
		UE_LOG(LogTemp, Error, TEXT("UserId not available"));
		return;
	}

	// Get ID and Auth Token
	FString AccountId = UserId->ToString();
	FString AuthToken = OnlineSubsystem->GetIdentityInterface()->GetAuthToken(0);

	// Print ID and Auth Token
	UE_LOG(LogTemp, Display, TEXT("Ready to send to backend - Steam ID: %s, Token: %s"), *AccountId, *AuthToken);
	FString DebugMessage = FString::Printf(TEXT("Ready to send to backend - ID: %s, Token: %s"), *AccountId, *AuthToken);
	GEngine->AddOnScreenDebugMessage(-1,10.0f, FColor::Green, DebugMessage);

	// Вызываем Blueprint событие
	OnSuccessfulLoginSteam(AccountId, AuthToken);
}

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

    UE_LOG(LogTemp, Display, TEXT("Attempting Persistent Auth (silent login)..."));

    FOnlineAccountCredentials Credentials;
    Credentials.Type = FString("persistentauth");

    PersistentLoginDelegateHandle = IdentityInterface->AddOnLoginCompleteDelegate_Handle(0,
        FOnLoginCompleteDelegate::CreateUObject(this, &UTagDuelsGameInstance::OnPersistentLoginComplete));

    if (!IdentityInterface->Login(0, Credentials))
    {
        UE_LOG(LogTemp, Error, TEXT("Persistent Auth Login() call failed immediately"));
        IdentityInterface->ClearOnLoginCompleteDelegate_Handle(0, PersistentLoginDelegateHandle);
        PersistentLoginDelegateHandle.Reset();
        // Fallback не вызываем здесь, так как Login() вернул false
    }
}

void UTagDuelsGameInstance::OnPersistentLoginComplete(int32 LocalUserNum, bool bWasSuccessful, const FUniqueNetId& UserId, const FString& Error)
{
    UE_LOG(LogTemp, Warning, TEXT("=== Persistent Auth Result ==="));
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
        // Автоматический вход сработал!
        UE_LOG(LogTemp, Display, TEXT("Silent login successful!"));
        HandleSuccessfulLogin(UserId, LocalUserNum);
        return;
    }

    // Persistent Auth не сработал – запускаем AccountPortal
    UE_LOG(LogTemp, Warning, TEXT("Persistent Auth failed, falling back to AccountPortal. Error: %s"), *Error);

    IOnlineIdentityPtr IdentityInterface = Online::GetSubsystem(GetWorld())->GetIdentityInterface();
    if (!IdentityInterface.IsValid())
    {
        UE_LOG(LogTemp, Error, TEXT("Identity interface lost before fallback"));
        return;
    }

    FOnlineAccountCredentials PortalCredentials;
    PortalCredentials.Type = FString("accountportal");

    LoginDelegateHandle = IdentityInterface->AddOnLoginCompleteDelegate_Handle(0,
        FOnLoginCompleteDelegate::CreateUObject(this, &UTagDuelsGameInstance::OnLoginComplete));

    if (!IdentityInterface->Login(0, PortalCredentials))
    {
        UE_LOG(LogTemp, Error, TEXT("Login() call failed immediately"));
        IdentityInterface->ClearOnLoginCompleteDelegate_Handle(0, LoginDelegateHandle);
        LoginDelegateHandle.Reset();
    }
    else
    {
        UE_LOG(LogTemp, Display, TEXT("Login launched"));
    }
}

void UTagDuelsGameInstance::OnLoginComplete(int32 LocalUserNum, bool bWasSuccessful, const FUniqueNetId& UserId, const FString& Error)
{
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
        HandleSuccessfulLogin(UserId, LocalUserNum);
    }
    else
    {
        UE_LOG(LogTemp, Error, TEXT("Login failed: %s"), *Error);
    }
}

void UTagDuelsGameInstance::HandleSuccessfulLogin(const FUniqueNetId& UserId, int32 LocalUserNum)
{
    IOnlineSubsystem* OnlineSubsystem = Online::GetSubsystem(GetWorld());
    if (!OnlineSubsystem) return;

    IOnlineIdentityPtr IdentityInterface = OnlineSubsystem->GetIdentityInterface();
    if (!IdentityInterface.IsValid()) return;

	FString AccountId = UserId.ToString();
	FString AuthToken = OnlineSubsystem->GetIdentityInterface()->GetAuthToken(LocalUserNum);

	// Print ID and Auth Token
	UE_LOG(LogTemp, Display, TEXT("Ready to send to backend - Epic ID: %s, Token: %s"), *AccountId, *AuthToken);
	FString DebugMessage = FString::Printf(TEXT("Ready to send to backend - ID: %s, Token: %s"), *AccountId, *AuthToken);
	GEngine->AddOnScreenDebugMessage(-1,10.0f, FColor::Green, DebugMessage);

	// Вызываем Blueprint событие
	OnSuccessfulLoginEOS(AccountId, AuthToken);
}

EOnlineSubsystemType UTagDuelsGameInstance::GetActiveOnlineSubsystemType() const
{
	FName SubsystemName = GetOnlinePlatformName(); // "Steam", "EOS", "NULL" и т.д.

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

void UTagDuelsGameInstance::DeletePersistentAuthEOS()
{
    // Получаем подсистему EOS
    IOnlineSubsystem* OnlineSubsystem = Online::GetSubsystem(GetWorld());
    if (!OnlineSubsystem || OnlineSubsystem->GetSubsystemName() != FName(TEXT("EOS")))
    {
        UE_LOG(LogTemp, Error, TEXT("EOS subsystem not available"));
        return;
    }

	IOnlineIdentityPtr IdentityInterface = OnlineSubsystem->GetIdentityInterface();
	if (IdentityInterface.IsValid())
	{
		// Завершаем сессию текущего пользователя
		IdentityInterface->Logout(0);
		UE_LOG(LogTemp, Display, TEXT("User logged out."));
	}

    // FOnlineSubsystemEOS* EOSSubsystem = static_cast<FOnlineSubsystemEOS*>(OnlineSubsystem);
    // EOS_HAuth AuthHandle = EOSSubsystem->GetEOSAuthHandle();
    //
    // if (!AuthHandle)
    // {
    //     UE_LOG(LogTemp, Error, TEXT("EOS Auth handle not available"));
    //     return;
    // }
    //
    // // Настройки для удаления токена
    // EOS_Auth_DeletePersistentAuthOptions DeleteOptions = {};
    // DeleteOptions.ApiVersion = EOS_AUTH_DELETEPERSISTENTAUTH_API_LATEST;
    // DeleteOptions.RefreshToken = nullptr;
    //
    // // Асинхронно удаляем токен
    // EOS_Auth_DeletePersistentAuth(AuthHandle, &DeleteOptions, nullptr, 
    //     [](const EOS_Auth_DeletePersistentAuthCallbackInfo* Data)
    //     {
    //         if (Data->ResultCode == EOS_EResult::EOS_Success)
    //         {
    //             UE_LOG(LogTemp, Display, TEXT("Persistent auth token successfully deleted"));
    //         }
    //         else
    //         {
    //             UE_LOG(LogTemp, Error, TEXT("Failed to delete persistent auth token. Error: %s"), 
    //                 ANSI_TO_TCHAR(EOS_EResult_ToString(Data->ResultCode)));
    //         }
    //     });
    //
    // // При желании можно также сразу разлогиниться
    // // IdentityInterface->Logout(0);
}
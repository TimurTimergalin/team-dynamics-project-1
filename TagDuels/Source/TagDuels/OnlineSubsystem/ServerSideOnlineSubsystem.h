// Fill out your copyright notice in the Description page of Project Settings.

#pragma once

#include "CoreMinimal.h"
#include "Delegates.h"
#include "Clients/MatchServiceClient.h"
#include "Contract/MatchHistory.h"
#include "Contract/ServerAnnotations.h"
#include "Subsystems/GameInstanceSubsystem.h"
#include "AgonesSubsystem.h"
#include "ServerSideOnlineSubsystem.generated.h"

UCLASS()
class TAGDUELS_API UServerSideOnlineSubsystem : public UGameInstanceSubsystem
{
	GENERATED_BODY()
public:
	virtual void Initialize(FSubsystemCollectionBase&) override;
	virtual bool ShouldCreateSubsystem(UObject* Outer) const override;
	
	UFUNCTION(BlueprintCallable, Category = "OnlineSubsystem")
	bool Ready(FReadyDelegate OnResponse, FAgonesErrorDelegate OnError);

	UFUNCTION(BlueprintCallable, Category = "OnlineSubsystem")
	bool Shutdown(FShutdownDelegate OnResponse, FAgonesErrorDelegate OnError);

	UFUNCTION(BlueprintCallable, Category = "OnlineSubsystem")
	bool GetPlayer1(FUserAnnotations& PlayerData);

	UFUNCTION(BlueprintCallable, Category = "OnlineSubsystem")
	bool GetPlayer2(FUserAnnotations& PlayerData);

	UFUNCTION(BlueprintCallable, Category = "OnlineSubsystem")
	bool DrawMatch(const TArray<FRoundData>& Rounds, FOnMatchEnd OnResponse, FOnErroneousResponse OnError);

	UFUNCTION(BlueprintCallable, Category = "OnlineSubsystem")
	bool EndMatch(int64 WinnerId, const TArray<FRoundData>& Rounds, FOnMatchEnd OnResponse, FOnErroneousResponse OnError);

	UFUNCTION(BlueprintCallable, Category = "OnlineSubsystem")
	bool CancelMatch(FOnEmptyResponse OnResponse, FOnErroneousResponse OnError);

	UFUNCTION(BlueprintCallable, Category = "OnlineSubsystem")
	bool RenewMatch(FOnEmptyResponse OnResponse, FOnErroneousResponse OnError);

	UFUNCTION(BlueprintCallable, Category = "OnlineSubsystem")
	bool ValidateConnection(const FString& Options, FString& ErrorMessage);
private:
	UFUNCTION()
	void OnAgonesUpdated(const FGameServerResponse& Response);
	
	TOptional<MatchServiceClient> MsClient;
	
	struct FMatchData
	{
		FUserAnnotations Player1{};
		FUserAnnotations Player2{};
		FString MatchId{};
	};
	TOptional<FMatchData> MatchData{};
	bool bReady = false;
	TWeakObjectPtr<UAgonesSubsystem> AgonesSubsystem;
};

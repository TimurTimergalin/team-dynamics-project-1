// Fill out your copyright notice in the Description page of Project Settings.

#pragma once

#include "CoreMinimal.h"
#include "Contract/MatchHistory.h"
#include "Contract/UserData.h"
#include "Subsystems/GameInstanceSubsystem.h"
#include "ClientSideOnlineSubsystem.generated.h"

DECLARE_DYNAMIC_DELEGATE(FOnEmptyResponse);

DECLARE_DYNAMIC_DELEGATE_OneParam(FOnErroneousResponse, FString, ErrorMessage);

DECLARE_DYNAMIC_DELEGATE_OneParam(FOnInt64Response, int64, Response);

DECLARE_DYNAMIC_DELEGATE_OneParam(FOnMatchHistoryResponse, const FMatchHistoryPage&, Response);

DECLARE_DYNAMIC_DELEGATE_OneParam(FOnUserListResponse, const FPlayersList&, Response);

DECLARE_DYNAMIC_DELEGATE_OneParam(FOnMatch, FString, Address);

UCLASS()
class TAGDUELS_API UClientSideOnlineSubsystem : public UGameInstanceSubsystem
{
	GENERATED_BODY()

public:
	UFUNCTION(BlueprintCallable, Category = "OnlineSubsystem")
	void SteamAuthorize(FString AuthToken, int64 SteamId, FOnEmptyResponse OnResponse);

	UFUNCTION(BlueprintCallable, Category = "OnlineSubsystem")
	bool GetPlayerData(FUserPlayerData& PlayerDataOut);

	UFUNCTION(BlueprintCallable, Category = "OnlineSubsystem")
	bool GetMatchHistoryPage(FString PageToken, FOnMatchHistoryResponse OnResponse, FOnErroneousResponse OnError);

	UFUNCTION(BlueprintCallable, Category = "OnlineSubsystem")
	bool GetFriendsList(FString PageToken, FOnUserListResponse OnResponse, FOnInt64Response OnError);

	UFUNCTION(BlueprintCallable, Category = "OnlineSubsystem")
	bool GetOutgoingRequests(FString PageToken, FOnUserListResponse OnResponse, FOnInt64Response OnError);

	UFUNCTION(BlueprintCallable, Category = "OnlineSubsystem")
	bool GetIncomingRequests(FString PageToken, FOnUserListResponse OnResponse, FOnInt64Response OnError);

	UFUNCTION(BlueprintCallable, Category = "OnlineSubsystem")
	bool GetRating(FOnInt64Response OnResponse, FOnInt64Response OnError);

	UFUNCTION(BlueprintCallable, Category = "OnlineSubsystem")
	void SubscribeToMatchStart(FOnMatch Callback);

	UFUNCTION(BlueprintCallable, Category = "OnlineSubsystem")
	void ClearMatchStartCallback();

	UFUNCTION(BlueprintCallable, Category = "OnlineSubsystem")
	void SendOrAcceptRequest(int64 OtherUserId, FOnEmptyResponse OnResponse, FOnErroneousResponse OnError);

	UFUNCTION(BlueprintCallable, Category = "OnlineSubsystem")
	void DeclineOrDeleteFriend(int64 OtherUserId, FOnEmptyResponse OnResponse, FOnErroneousResponse OnError);

	UFUNCTION(BlueprintCallable, Category = "OnlineSubsystem")
	void ChallengeUser(int64 OtherUserId, FOnEmptyResponse OnResponse, FOnEmptyResponse OnError);

	UFUNCTION(BlueprintCallable, Category = "OnlineSubsystem")
	void CancelChallenge(FOnEmptyResponse OnResponse, FOnEmptyResponse OnError);

	UFUNCTION(BlueprintCallable, Category = "OnlineSubsystem")
	void StartMatchmaking(FOnEmptyResponse OnResponse, FOnEmptyResponse OnError);

	UFUNCTION(BlueprintCallable, Category = "OnlineSubsystem")
	void SubscribeToMMEventError(FOnErroneousResponse OnError);

	UFUNCTION(BlueprintCallable, Category = "OnlineSubsystem")
	void ClearMMEventErrorCallback();

	UFUNCTION(BlueprintCallable, Category = "OnlineSubsystem")
	void SubscribeToUserEventError(FOnErroneousResponse OnError);

	UFUNCTION(BlueprintCallable, Category = "OnlineSubsystem")
	void ClearUserEventErrorCallback();

	UFUNCTION(BlueprintCallable, Category = "OnlineSubsystem")
	void ConnectToMMEvent(FOnEmptyResponse OnResponse, FOnErroneousResponse OnError);

	UFUNCTION(BlueprintCallable, Category = "OnlineSubsystem")
	void ConnectToUserEvent(FOnEmptyResponse OnResponse, FOnErroneousResponse OnError);

private:
	TOptional<FUserPlayerData> PlayerData{};
	FOnMatch OnMatchCallback{};
	TOptional<int64> ChallengedUserId{};
	bool InMatchmaking = false;
	FOnErroneousResponse OnMMEventErrorCallback{};
	FOnErroneousResponse OnUserEventErrorCallback{};
};

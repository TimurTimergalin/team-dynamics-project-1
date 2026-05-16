// Fill out your copyright notice in the Description page of Project Settings.

#pragma once

#include "CoreMinimal.h"
#include "Delegates.h"
#include "Clients/AuthServiceClient.h"
#include "Clients/MatchHistoryServiceV2Client.h"
#include "Clients/MMEventClient.h"
#include "Clients/UserEventClient.h"
#include "Clients/UserServiceClient.h"
#include "Contract/UserData.h"
#include "Subsystems/GameInstanceSubsystem.h"
#include "ClientSideOnlineSubsystem.generated.h"

UCLASS()
class TAGDUELS_API UClientSideOnlineSubsystem : public UGameInstanceSubsystem
{
	GENERATED_BODY()

public:
	virtual void Initialize(FSubsystemCollectionBase&) override;
	
	UFUNCTION(BlueprintCallable, Category = "OnlineSubsystem")
	bool SteamAuthorize(FString AuthToken, int64 SteamId, FOnEmptyResponse OnResponse);

	UFUNCTION(BlueprintCallable, Category = "OnlineSubsystem")
	bool EgsAuthorize(FString AuthToken, FString Id, FOnEmptyResponse OnResponse);

	UFUNCTION(BlueprintCallable, Category = "OnlineSubsystem")
	bool GetPlayerData(FUserPlayerData& PlayerDataOut);

	UFUNCTION(BlueprintCallable, Category = "OnlineSubsystem")
	bool GetMatchHistoryPage(FString PageToken, FOnMatchHistoryResponse OnResponse, FOnErroneousResponse OnError);

	UFUNCTION(BlueprintCallable, Category = "OnlineSubsystem")
	bool GetFriendsList(FString PageToken, FOnUserListResponse OnResponse, FOnErroneousResponse OnError);

	UFUNCTION(BlueprintCallable, Category = "OnlineSubsystem")
	bool GetOutgoingRequests(FString PageToken, FOnUserListResponse OnResponse, FOnErroneousResponse OnError);

	UFUNCTION(BlueprintCallable, Category = "OnlineSubsystem")
	bool GetIncomingRequests(FString PageToken, FOnUserListResponse OnResponse, FOnErroneousResponse OnError);

	UFUNCTION(BlueprintCallable, Category = "OnlineSubsystem")
	bool GetRating(FOnInt64Response OnResponse, FOnErroneousResponse OnError);

	UFUNCTION(BlueprintCallable, Category = "OnlineSubsystem")
	bool GetRatingById(int64 OtherUserId, FOnInt64Response OnResponse, FOnErroneousResponse OnError);

	UFUNCTION(BlueprintCallable, Category = "OnlineSubsystem")
	void SendOrAcceptRequest(int64 OtherUserId, FOnEmptyResponse OnResponse, FOnErroneousResponse OnError);

	UFUNCTION(BlueprintCallable, Category = "OnlineSubsystem")
	void DeclineOrDeleteFriend(int64 OtherUserId, FOnEmptyResponse OnResponse, FOnErroneousResponse OnError);

	UFUNCTION(BlueprintCallable, Category = "OnlineSubsystem")
	bool ConnectToMMEvent(FOnMatch OnResponse, FOnErroneousResponse OnError);

	UFUNCTION(BlueprintCallable, Category = "OnlineSubsystem")
	bool StartMatchMaking();

	UFUNCTION(BlueprintCallable, Category = "OnlineSubsystem")
	bool CancelMatchMaking();

	UFUNCTION(BlueprintCallable, Category = "OnlineSubsystem")
	bool DisconnectFromMMEvent();

	UFUNCTION(BlueprintCallable, Category = "OnlineSubsystem")
	bool ConnectToUserEvent(
		FOnStatusUpdated OnStatusChanged,
		FOnChallengeReceived OnChallengeReceived,
		FOnMatch OnMatchStarted,
		FOnEmptyResponse OnChallengeDeclined,
		FOnChallengeCancelled OnChallengeCancelled,
		FOnErroneousResponse OnError
	);

	UFUNCTION(BlueprintCallable, Category = "OnlineSubsystem")
	bool DisconnectFromUserEvent();

	UFUNCTION(BlueprintCallable, Category = "OnlineSubsystem")
	bool SubscribeToUsers(TArray<int64> UserIds);

	UFUNCTION(BlueprintCallable, Category = "OnlineSubsystem")
	bool ChallengeUserEvent(int64 UserId);

	UFUNCTION(BlueprintCallable, Category = "OnlineSubsystem")
	bool CancelChallengeUserEvent();

	UFUNCTION(BlueprintCallable, Category = "OnlineSubsystem")
	bool AcceptChallengeUserEvent(FChallengeReceivedEvent Challenge);

	UFUNCTION(BlueprintCallable, Category = "OnlineSubsystem")
	bool DeclineChallengeUserEvent(FChallengeReceivedEvent Challenge);

	UFUNCTION(BlueprintCallable, Category = "OnlineSubsystem")
	bool NotifyBusy();

	UFUNCTION(BlueprintCallable, Category = "OnlineSubsystem")
	bool NotifyFree();

	UFUNCTION(BlueprintCallable, Category = "OnlineSubsystemDebug")
	void SetPlayerData(const FString& Name, int64 PlayerId);

private:
	TFuture<bool> UpdateTokens(FOnErroneousResponse OnError);

	TOptional<UserServiceClient> UsClient{};
	TOptional<AuthServiceClient> AsClient{};
	TOptional<MatchHistoryServiceV2Client> MhsV2Client{};
	TOptional<MMEventClient> MmeClient{};
	TOptional<UserEventClient> UeClient{};
	
	TOptional<FUserPlayerData> PlayerData{};
	TOptional<FTokenWithExp> AccessToken{};
	TOptional<FTokenWithExp> RefreshToken{};
	TOptional<int64> ChallengedUserId{};
	bool InMatchmaking = false;
	FOnErroneousResponse OnUserEventErrorCallback{};
};

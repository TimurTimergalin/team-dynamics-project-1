// Fill out your copyright notice in the Description page of Project Settings.


#include "ClientSideOnlineSubsystem.h"

void UClientSideOnlineSubsystem::SteamAuthorize(FString AuthToken, int64 SteamId, FOnEmptyResponse OnResponse)
{
}

bool UClientSideOnlineSubsystem::GetPlayerData(FUserPlayerData& PlayerDataOut)
{
	return true;
}

bool UClientSideOnlineSubsystem::GetMatchHistoryPage(FString PageToken, FOnMatchHistoryResponse OnResponse,
                                                     FOnErroneousResponse OnError)
{
	return false;
}

bool UClientSideOnlineSubsystem::GetFriendsList(FString PageToken, FOnUserListResponse OnResponse,
                                                FOnInt64Response OnError)
{
	return false;
}

bool UClientSideOnlineSubsystem::GetOutgoingRequests(FString PageToken, FOnUserListResponse OnResponse,
                                                     FOnInt64Response OnError)
{
	return false;
}

bool UClientSideOnlineSubsystem::GetIncomingRequests(FString PageToken, FOnUserListResponse OnResponse,
                                                     FOnInt64Response OnError)
{
	return false;
}

bool UClientSideOnlineSubsystem::GetRating(FOnInt64Response OnResponse, FOnInt64Response OnError)
{
	return false;
}

void UClientSideOnlineSubsystem::SubscribeToMatchStart(FOnMatch Callback)
{
}

void UClientSideOnlineSubsystem::ClearMatchStartCallback()
{
}

void UClientSideOnlineSubsystem::SendOrAcceptRequest(int64 OtherUserId, FOnEmptyResponse OnResponse,
                                                     FOnErroneousResponse OnError)
{
}

void UClientSideOnlineSubsystem::DeclineOrDeleteFriend(int64 OtherUserId, FOnEmptyResponse OnResponse,
                                                       FOnErroneousResponse OnError)
{
}

void UClientSideOnlineSubsystem::ChallengeUser(int64 OtherUserId, FOnEmptyResponse OnResponse, FOnEmptyResponse OnError)
{
}

void UClientSideOnlineSubsystem::CancelChallenge(FOnEmptyResponse OnResponse, FOnEmptyResponse OnError)
{
}

void UClientSideOnlineSubsystem::StartMatchmaking(FOnEmptyResponse OnResponse, FOnEmptyResponse OnError)
{
}

void UClientSideOnlineSubsystem::SubscribeToMMEventError(FOnErroneousResponse OnError)
{
}

void UClientSideOnlineSubsystem::ClearMMEventErrorCallback()
{
}

void UClientSideOnlineSubsystem::SubscribeToUserEventError(FOnErroneousResponse OnError)
{
}

void UClientSideOnlineSubsystem::ClearUserEventErrorCallback()
{
}

void UClientSideOnlineSubsystem::ConnectToMMEvent(FOnEmptyResponse OnResponse, FOnErroneousResponse OnError)
{
}

void UClientSideOnlineSubsystem::ConnectToUserEvent(FOnEmptyResponse OnResponse, FOnErroneousResponse OnError)
{
}

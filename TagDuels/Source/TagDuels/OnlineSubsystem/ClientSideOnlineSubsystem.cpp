#include "ClientSideOnlineSubsystem.h"

#include "utils.h"

void UClientSideOnlineSubsystem::Initialize(FSubsystemCollectionBase& SubsystemCollectionBase)
{
	Super::Initialize(SubsystemCollectionBase);
	UsClient = CreateUserServiceClient();
	MhsClient = CreateMatchHistoryServiceClient();
	RsClient = CreateRatingServiceClient();
}

bool UClientSideOnlineSubsystem::SteamAuthorize(FString /*AuthToken*/, int64 SteamId, FOnEmptyResponse OnResponse)
{
	if (!UsClient)
	{
		return false;
	}
	WithRetry<FUserPlayerData>([this, SteamId]()
	{
		return UsClient->GetSelfData(SteamId);
	}, 3).Next([this, OnResponse](TOptional<FUserPlayerData> PlayerData) {
		this->PlayerData = PlayerData;
		return OnGameThread(&decltype(OnResponse)::ExecuteIfBound, OnResponse);
	}).Next(FutureJoin).Next([](bool bWasBound)
	{
		if (!bWasBound)
		{
			UE_LOG(LogTemp, Error, TEXT("OnResponse for SteamAuthorize was not set"));
		}
	});
	return true;
}

bool UClientSideOnlineSubsystem::EgsAuthorize(FString AuthToken, int64 Id, FOnEmptyResponse OnResponse)
{
	return false;
}

bool UClientSideOnlineSubsystem::GetPlayerData(FUserPlayerData& PlayerDataOut)
{
	if (PlayerData.IsSet())
	{
		PlayerDataOut = PlayerData.GetValue();
		return true;
	}
	return false;
}

bool UClientSideOnlineSubsystem::GetMatchHistoryPage(FString PageToken, FOnMatchHistoryResponse OnResponse,
                                                     FOnErroneousResponse OnError)
{
	if (!MhsClient)
	{
		return false;
	}
	if (!PlayerData.IsSet())
	{
		return false;
	}

	WithRetry<FMatchHistoryPage>([this, PageToken]()
	{
		return MhsClient->GetMatchHistory(PlayerData.GetValue().Id, PageToken);
	}, 3).Next([this, OnResponse, OnError](TOptional<FMatchHistoryPage> Page)
	{
		if (!Page.IsSet())
		{
			return OnGameThread(&decltype(OnError)::ExecuteIfBound, OnError, "Failed to fetch match history after 3 attempts");
		}
		return OnGameThread(&decltype(OnResponse)::ExecuteIfBound, OnResponse, Page.GetValue());
	}).Next(FutureJoin).Next([](bool bWasBound)
	{
		if (!bWasBound)
		{
			UE_LOG(LogTemp, Error, TEXT("OnResponse or OnError for GetMatchHistoryPage was not set"));
		}
	});
	return true;
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

bool UClientSideOnlineSubsystem::GetRating(FOnInt64Response OnResponse, FOnErroneousResponse OnError)
{
	if (!RsClient)
	{
		return false;
	}
	if (!PlayerData.IsSet())
	{
		return false;
	}
	WithRetry<int64>([this]()
	{
		return RsClient->GetRating(PlayerData.GetValue().Id);
	}, 3).Next([OnResponse, OnError](TOptional<int64> Rating)
	{
		if (!Rating.IsSet())
		{
			return OnGameThread(&decltype(OnError)::ExecuteIfBound, OnError, FString(TEXT("Failed to fetch rating after 3 attempts")));
		}
		return OnGameThread(&decltype(OnResponse)::ExecuteIfBound, OnResponse, Rating.GetValue());
	}).Next(FutureJoin).Next([](bool bWasBound)
	{
		if (!bWasBound)
		{
			UE_LOG(LogTemp, Error, TEXT("OnResponse or OnError for GetRating was not set"));
		}
	});
	return true;
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

bool UClientSideOnlineSubsystem::ConnectToMMEvent(FOnMatch OnResponse, FOnErroneousResponse OnError)
{
	if (!PlayerData.IsSet())
	{
		return false;
	}
	if (MmeClient.IsSet() && MmeClient->IsConnected())
	{
		return false;
	}
	MmeClient = CreateMMEventClient(PlayerData->Id, OnResponse, OnError);
	return MmeClient.IsSet();
}

bool UClientSideOnlineSubsystem::StartMatchMaking()
{
	if (!MmeClient.IsSet())
	{
		return false;
	}
	if (!MmeClient->IsConnected())
	{
		return false;
	}
	return MmeClient->StartMatchmaking();;
}

bool UClientSideOnlineSubsystem::CancelMatchMaking()
{
	if (!MmeClient.IsSet())
	{
		return false;
	}
	if (!MmeClient->IsConnected())
	{
		return false;
	}
	return MmeClient->CancelMatchmaking();
}

bool UClientSideOnlineSubsystem::DisconnectFromMMEvent()
{
	if (!MmeClient.IsSet())
	{
		return false;
	}
	MmeClient->Close();
	return true;
}

void UClientSideOnlineSubsystem::ConnectToUserEvent(FOnEmptyResponse OnResponse, FOnErroneousResponse OnError)
{
}

#include "ClientSideOnlineSubsystem.h"

#include "utils.h"

void UClientSideOnlineSubsystem::Initialize(FSubsystemCollectionBase& SubsystemCollectionBase)
{
	Super::Initialize(SubsystemCollectionBase);
	UsClient = CreateUserServiceClient();
	AsClient = CreateAuthServiceClient();
	MhsV2Client = CreateMatchHistoryServiceV2Client();
	MmeClient = CreateMMEventClient();
	UeClient = CreateUserEventClient();
}

bool UClientSideOnlineSubsystem::SteamAuthorize(FString AuthToken, int64 SteamId, FOnEmptyResponse OnResponse)
{
	if (!AsClient || !UsClient)
	{
		return false;
	}
	AsClient->AuthExternal(SteamId, AuthToken).Next([this, OnResponse](TOptional<FAuthExternalResponse> AuthResp)
	{
		if (!AuthResp.IsSet())
		{
			UE_LOG(LogTemp, Error, TEXT("SteamAuthorize: AuthExternal failed"));
			return OnGameThread(&decltype(OnResponse)::ExecuteIfBound, OnResponse);
		}
		AccessToken = AuthResp->Access;
		RefreshToken = AuthResp->Refresh;
		const int64 UserId = AuthResp->UserId;
		return UsClient->GetUserData(UserId, AuthResp->Access.Token).Next([this, OnResponse](TOptional<FUserPlayerData> UserData)
		{
			this->PlayerData = UserData;
			return OnGameThread(&decltype(OnResponse)::ExecuteIfBound, OnResponse);
		}).Next(FutureJoin).Next([](bool bWasBound)
		{
			if (!bWasBound)
			{
				UE_LOG(LogTemp, Error, TEXT("OnResponse for SteamAuthorize was not set"));
			}
			return true;
		});
	});
	return true;
}

TFuture<bool> UClientSideOnlineSubsystem::UpdateTokens(FOnErroneousResponse OnError)
{
	if (!AsClient)
	{
		return MakeFulfilledPromise<bool>(false).GetFuture();
	}

	const FDateTime Now = FDateTime::UtcNow();

	if (!RefreshToken.IsSet() || RefreshToken->Expiry <= Now)
	{
		OnGameThread(&FOnErroneousResponse::ExecuteIfBound, OnError,
			FOnlineSubsystemError{TEXT("Your session is expired, please reconnect"), EOnlineErrorType::Critical});
		return MakeFulfilledPromise<bool>(false).GetFuture();
	}

	if (AccessToken.IsSet() && AccessToken->Expiry > Now)
	{
		return MakeFulfilledPromise<bool>(true).GetFuture();
	}

	return AsClient->Refresh(RefreshToken->Token).Next([this, OnError](TOptional<FRefreshResponse> Resp)
	{
		if (!Resp.IsSet())
		{
			UE_LOG(LogTemp, Error, TEXT("UpdateTokens: Refresh failed"));
			return false;
		}
		AccessToken = Resp->Access;
		RefreshToken = Resp->Refresh;
		return true;
	});
}

bool UClientSideOnlineSubsystem::EgsAuthorize(FString AuthToken, FString Id, FOnEmptyResponse OnResponse)
{
	if (!AsClient || !UsClient)
	{
		return false;
	}
	AsClient->AuthExternal(Id, AuthToken).Next([this, OnResponse](TOptional<FAuthExternalResponse> AuthResp)
	{
		if (!AuthResp.IsSet())
		{
			UE_LOG(LogTemp, Error, TEXT("EgsAuthorize: AuthExternal failed"));
			return OnGameThread(&decltype(OnResponse)::ExecuteIfBound, OnResponse);
		}
		AccessToken = AuthResp->Access;
		RefreshToken = AuthResp->Refresh;
		const int64 UserId = AuthResp->UserId;
		return UsClient->GetUserData(UserId, AuthResp->Access.Token).Next([this, OnResponse](TOptional<FUserPlayerData> UserData)
		{
			this->PlayerData = UserData;
			return OnGameThread(&decltype(OnResponse)::ExecuteIfBound, OnResponse);
		}).Next(FutureJoin).Next([](bool bWasBound)
		{
			if (!bWasBound)
			{
				UE_LOG(LogTemp, Error, TEXT("OnResponse for EgsAuthorize was not set"));
			}
			return true;
		});
	});
	return true;
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
	if (!MhsV2Client || !PlayerData.IsSet())
	{
		return false;
	}
	UpdateTokens(OnError).Next([this, PageToken, OnResponse, OnError](bool bOk)
	{
		if (!bOk) return;
		WithRetry<FMatchHistoryPage>([this, PageToken]()
		{
			return MhsV2Client->GetMatchHistory(PlayerData.GetValue().Id, AccessToken->Token, PageToken);
		}, 3).Next([this, OnResponse, OnError](TOptional<FMatchHistoryPage> Page)
		{
			if (!Page.IsSet())
			{
				return OnGameThread(&decltype(OnError)::ExecuteIfBound, OnError, FOnlineSubsystemError{TEXT("Failed to fetch match history after 3 attempts"), EOnlineErrorType::NonCritical});
			}
			return OnGameThread(&decltype(OnResponse)::ExecuteIfBound, OnResponse, Page.GetValue());
		}).Next(FutureJoin).Next([](bool bWasBound)
		{
			if (!bWasBound)
			{
				UE_LOG(LogTemp, Error, TEXT("OnResponse or OnError for GetMatchHistoryPage was not set"));
			}
		});
	});
	return true;
}

bool UClientSideOnlineSubsystem::GetFriendsList(FString PageToken, FOnUserListResponse OnResponse,
                                                FOnErroneousResponse OnError)
{
	if (!UsClient || !PlayerData.IsSet())
	{
		return false;
	}
	UpdateTokens(OnError).Next([this, PageToken, OnResponse, OnError](bool bOk)
	{
		if (!bOk) return;
		WithRetry<FPlayersList>([this, PageToken]()
		{
			return UsClient->GetFriends(PlayerData.GetValue().Id, PageToken, AccessToken->Token);
		}, 3).Next([OnResponse, OnError](TOptional<FPlayersList> Result)
		{
			if (!Result.IsSet())
			{
				return OnGameThread(&decltype(OnError)::ExecuteIfBound, OnError, FOnlineSubsystemError{TEXT("Failed to fetch friends list after 3 attempts"), EOnlineErrorType::NonCritical});
			}
			return OnGameThread(&decltype(OnResponse)::ExecuteIfBound, OnResponse, Result.GetValue());
		});
	});
	return true;
}

bool UClientSideOnlineSubsystem::GetOutgoingRequests(FString PageToken, FOnUserListResponse OnResponse,
                                                     FOnErroneousResponse OnError)
{
	if (!UsClient || !PlayerData.IsSet())
	{
		return false;
	}
	UpdateTokens(OnError).Next([this, PageToken, OnResponse, OnError](bool bOk)
	{
		if (!bOk) return;
		WithRetry<FPlayersList>([this, PageToken]()
		{
			return UsClient->GetOutgoingRequests(PlayerData.GetValue().Id, PageToken, AccessToken->Token);
		}, 3).Next([OnResponse, OnError](TOptional<FPlayersList> Result)
		{
			if (!Result.IsSet())
			{
				return OnGameThread(&decltype(OnError)::ExecuteIfBound, OnError, FOnlineSubsystemError{TEXT("Failed to fetch outgoing requests after 3 attempts"), EOnlineErrorType::NonCritical});
			}
			return OnGameThread(&decltype(OnResponse)::ExecuteIfBound, OnResponse, Result.GetValue());
		});
	});
	return true;
}

bool UClientSideOnlineSubsystem::GetIncomingRequests(FString PageToken, FOnUserListResponse OnResponse,
                                                     FOnErroneousResponse OnError)
{
	if (!UsClient || !PlayerData.IsSet())
	{
		return false;
	}
	UpdateTokens(OnError).Next([this, PageToken, OnResponse, OnError](bool bOk)
	{
		if (!bOk) return;
		WithRetry<FPlayersList>([this, PageToken]()
		{
			return UsClient->GetIncomingRequests(PlayerData.GetValue().Id, PageToken, AccessToken->Token);
		}, 3).Next([OnResponse, OnError](TOptional<FPlayersList> Result)
		{
			if (!Result.IsSet())
			{
				return OnGameThread(&decltype(OnError)::ExecuteIfBound, OnError, FOnlineSubsystemError{TEXT("Failed to fetch incoming requests after 3 attempts"), EOnlineErrorType::NonCritical});
			}
			return OnGameThread(&decltype(OnResponse)::ExecuteIfBound, OnResponse, Result.GetValue());
		});
	});
	return true;
}

bool UClientSideOnlineSubsystem::GetRating(FOnInt64Response OnResponse, FOnErroneousResponse OnError)
{
	if (!MhsV2Client || !PlayerData.IsSet())
	{
		return false;
	}
	UpdateTokens(OnError).Next([this, OnResponse, OnError](bool bOk)
	{
		if (!bOk) return;
		WithRetry<int64>([this]()
		{
			return MhsV2Client->GetRating(PlayerData.GetValue().Id, AccessToken->Token);
		}, 3).Next([OnResponse, OnError](TOptional<int64> Rating)
		{
			if (!Rating.IsSet())
			{
				return OnGameThread(&decltype(OnError)::ExecuteIfBound, OnError, FOnlineSubsystemError{TEXT("Failed to fetch rating after 3 attempts"), EOnlineErrorType::NonCritical});
			}
			return OnGameThread(&decltype(OnResponse)::ExecuteIfBound, OnResponse, Rating.GetValue());
		}).Next(FutureJoin).Next([](bool bWasBound)
		{
			if (!bWasBound)
			{
				UE_LOG(LogTemp, Error, TEXT("OnResponse or OnError for GetRating was not set"));
			}
		});
	});
	return true;
}

bool UClientSideOnlineSubsystem::GetRatingById(int64 OtherUserId, FOnInt64Response OnResponse, FOnErroneousResponse OnError)
{
	if (!MhsV2Client || !PlayerData.IsSet())
	{
		return false;
	}
	UpdateTokens(OnError).Next([this, OtherUserId, OnResponse, OnError](bool bOk)
	{
		if (!bOk) return;
		WithRetry<int64>([this, OtherUserId]()
		{
			return MhsV2Client->GetRating(OtherUserId, AccessToken->Token);
		}, 3).Next([OnResponse, OnError](TOptional<int64> Rating)
		{
			if (!Rating.IsSet())
			{
				return OnGameThread(&decltype(OnError)::ExecuteIfBound, OnError, FOnlineSubsystemError{TEXT("Failed to fetch rating after 3 attempts"), EOnlineErrorType::NonCritical});
			}
			return OnGameThread(&decltype(OnResponse)::ExecuteIfBound, OnResponse, Rating.GetValue());
		}).Next(FutureJoin).Next([](bool bWasBound)
		{
			if (!bWasBound)
			{
				UE_LOG(LogTemp, Error, TEXT("OnResponse or OnError for GetRating was not set"));
			}
		});
	});
	return true;
}

void UClientSideOnlineSubsystem::SendOrAcceptRequest(int64 OtherUserId, FOnEmptyResponse OnResponse,
                                                     FOnErroneousResponse OnError)
{
	if (!UsClient || !PlayerData.IsSet())
	{
		return;
	}
	UpdateTokens(OnError).Next([this, OtherUserId, OnResponse, OnError](bool bOk)
	{
		if (!bOk) return;
		UsClient->AddFriend(PlayerData.GetValue().Id, OtherUserId, AccessToken->Token).Next([OnResponse, OnError](bool bSuccess)
		{
			if (!bSuccess)
			{
				return OnGameThread(&decltype(OnError)::ExecuteIfBound, OnError, FOnlineSubsystemError{TEXT("Failed to send/accept friend request"), EOnlineErrorType::NonCritical});
			}
			return OnGameThread(&decltype(OnResponse)::ExecuteIfBound, OnResponse);
		});
	});
}

void UClientSideOnlineSubsystem::DeclineOrDeleteFriend(int64 OtherUserId, FOnEmptyResponse OnResponse,
                                                       FOnErroneousResponse OnError)
{
	if (!UsClient || !PlayerData.IsSet())
	{
		return;
	}
	UpdateTokens(OnError).Next([this, OtherUserId, OnResponse, OnError](bool bOk)
	{
		if (!bOk) return;
		UsClient->RemoveFriend(PlayerData.GetValue().Id, OtherUserId, AccessToken->Token).Next([OnResponse, OnError](bool bSuccess)
		{
			if (!bSuccess)
			{
				return OnGameThread(&decltype(OnError)::ExecuteIfBound, OnError, FOnlineSubsystemError{TEXT("Failed to decline/remove friend"), EOnlineErrorType::NonCritical});
			}
			return OnGameThread(&decltype(OnResponse)::ExecuteIfBound, OnResponse);
		});
	});
}

bool UClientSideOnlineSubsystem::ConnectToMMEvent(FOnMatch OnResponse, FOnErroneousResponse OnError)
{
	if (!PlayerData.IsSet() || !MmeClient.IsSet() || MmeClient->IsConnected())
	{
		return false;
	}
	UpdateTokens(OnError).Next([this, OnResponse, OnError](bool bOk)
	{
		if (!bOk) return;
		MmeClient->EstablishConnection(MakeShared<FOnMatch>(OnResponse), MakeShared<FOnErroneousResponse>(OnError), PlayerData->Id, AccessToken->Token);
	});
	return true;
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
	return MmeClient->Close();
}

bool UClientSideOnlineSubsystem::ConnectToUserEvent(
	FOnStatusUpdated OnStatusChanged,
	FOnChallengeReceived OnChallengeReceived,
	FOnMatch OnMatchStarted,
	FOnEmptyResponse OnChallengeDeclined,
	FOnChallengeCancelled OnChallengeCancelled,
	FOnErroneousResponse OnError)
{
	if (!PlayerData.IsSet() || !UeClient.IsSet() || UeClient->IsConnected())
	{
		return false;
	}
	UpdateTokens(OnError).Next([this,
		OnStatusChanged, OnChallengeReceived, OnMatchStarted, OnChallengeDeclined, OnChallengeCancelled, OnError](bool bOk)
	{
		if (!bOk) return;
		UeClient->EstablishConnection(
			MakeShared<FOnStatusUpdated>(OnStatusChanged),
			MakeShared<FOnChallengeReceived>(OnChallengeReceived),
			MakeShared<FOnMatch>(OnMatchStarted),
			MakeShared<FOnEmptyResponse>(OnChallengeDeclined),
			MakeShared<FOnChallengeCancelled>(OnChallengeCancelled),
			MakeShared<FOnErroneousResponse>(OnError),
			PlayerData->Id,
			AccessToken->Token
		);
	});
	return true;
}

bool UClientSideOnlineSubsystem::DisconnectFromUserEvent()
{
	if (!UeClient.IsSet())
	{
		return false;
	}
	return UeClient->Close();
}

bool UClientSideOnlineSubsystem::SubscribeToUsers(TArray<int64> UserIds)
{
	if (!UeClient.IsSet())
	{
		return false;
	}
	return UeClient->Subscribe(UserIds);
}

bool UClientSideOnlineSubsystem::ChallengeUserEvent(int64 UserId)
{
	if (!UeClient.IsSet())
	{
		return false;
	}
	return UeClient->Challenge(UserId);
}

bool UClientSideOnlineSubsystem::CancelChallengeUserEvent()
{
	if (!UeClient.IsSet())
	{
		return false;
	}
	return UeClient->CancelChallenge();
}

bool UClientSideOnlineSubsystem::AcceptChallengeUserEvent(FChallengeReceivedEvent Challenge)
{
	if (!UeClient.IsSet())
	{
		return false;
	}
	return UeClient->AcceptChallenge(Challenge);
}

bool UClientSideOnlineSubsystem::DeclineChallengeUserEvent(FChallengeReceivedEvent Challenge)
{
	if (!UeClient.IsSet())
	{
		return false;
	}
	return UeClient->DeclineChallenge(Challenge);
}

bool UClientSideOnlineSubsystem::NotifyBusy()
{
	if (!UeClient.IsSet())
	{
		return false;
	}
	return UeClient->NotifyBusy();
}

bool UClientSideOnlineSubsystem::NotifyFree()
{
	if (!UeClient.IsSet())
	{
		return false;
	}
	return UeClient->NotifyFree();
}

void UClientSideOnlineSubsystem::SetPlayerData(const FString& Name, int64 PlayerId)
{
	FUserPlayerData Pd;
	Pd.Id = PlayerId;
	Pd.Name = Name;
	PlayerData = Pd;
}

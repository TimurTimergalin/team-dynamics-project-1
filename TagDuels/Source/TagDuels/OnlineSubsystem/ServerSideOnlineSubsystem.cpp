#include "ServerSideOnlineSubsystem.h"

#include "utils.h"
#include "Kismet/GameplayStatics.h"

void UServerSideOnlineSubsystem::Initialize(FSubsystemCollectionBase& SubsystemCollectionBase)
{
	Super::Initialize(SubsystemCollectionBase);
	MsClient = CreateMatchServiceClient();
}

bool UServerSideOnlineSubsystem::ShouldCreateSubsystem(UObject* Outer) const
{
	return UE_SERVER;
}

bool UServerSideOnlineSubsystem::ValidateConnection(const FString& Options, FString& ErrorMessage, int64& PlayerID)
{
	if (!bReady)
	{
		ErrorMessage = "Server unavailable";
		return false;
	}
	if (!MatchData)
	{
		ErrorMessage = "MatchData missing";
		return false;
	}
	FString PlayerId = UGameplayStatics::ParseOption(Options, TEXT("player_id"));
	if (PlayerId == "")
	{
		ErrorMessage = "No Player Id";
		return false;
	}
	TOptional<int64> PlayerIdInt = StrToInt64(PlayerId);
	if (!PlayerIdInt)
	{
		ErrorMessage = "Invalid Player Id";
		return false;
	}
	if (PlayerIdInt.GetValue() != MatchData.GetValue().Player1.PlayerId && PlayerIdInt.GetValue() != MatchData.GetValue().Player2.PlayerId)
	{
		ErrorMessage = "Unknown Player Id";
		return false;
	}
	PlayerID = PlayerIdInt.GetValue();
	return true;
}

void UServerSideOnlineSubsystem::OnAgonesUpdated(const FGameServerResponse& Response)
{
	bReady = Response.Status.State == "Allocated";
	const TMap<FString, FString>& Annotations = Response.ObjectMeta.Annotations;
	FMatchData Result{};

	if (const FString* P1Id = Annotations.Find("player1_id"); P1Id)
	{
		if (auto P1IdInt = StrToInt64(*P1Id); P1IdInt.IsSet())
		{
			Result.Player1.PlayerId = P1IdInt.GetValue();
		}
		else { return; }
	}
	else { return; }

	if (const FString* P1Name = Annotations.Find("player1_name"); P1Name)
	{
		Result.Player1.Name = *P1Name;
	}
	else { return; }

	if (const FString* P1Rating = Annotations.Find("player1_rating"); P1Rating)
	{
		if (auto P1RatingInt = StrToInt64(*P1Rating); P1RatingInt.IsSet())
		{
			Result.Player1.Rating = P1RatingInt.GetValue();
		}
		else { return; }
	}
	else { return; }

	if (const FString* P2Id = Annotations.Find("player2_id"); P2Id)
	{
		if (auto P2IdInt = StrToInt64(*P2Id); P2IdInt.IsSet())
		{
			Result.Player2.PlayerId = P2IdInt.GetValue();
		}
		else { return; }
	}
	else { return; }

	if (const FString* P2Name = Annotations.Find("player2_name"); P2Name)
	{
		Result.Player2.Name = *P2Name;
	}
	else { return; }

	if (const FString* P2Rating = Annotations.Find("player2_rating"); P2Rating)
	{
		if (auto P2RatingInt = StrToInt64(*P2Rating); P2RatingInt.IsSet())
		{
			Result.Player2.Rating = P2RatingInt.GetValue();
		}
		else { return; }
	}
	else { return; }

	if (const FString* MatchId = Annotations.Find("match_id"); MatchId)
	{
		Result.MatchId = *MatchId;
	}
	else { return; }

	MatchData = MoveTemp(Result);
}

bool UServerSideOnlineSubsystem::Ready(FReadyDelegate OnResponse, FAgonesErrorDelegate OnError)
{
	UAgonesSubsystem* AgonesSDK = UAgonesSubsystem::Get(this);
	if (!AgonesSDK)
	{
		return false;
	}
	FGameServerDelegate Delegate;
	Delegate.BindDynamic(this, &UServerSideOnlineSubsystem::OnAgonesUpdated);
	AgonesSDK->WatchGameServer(Delegate);
	AgonesSDK->Ready(OnResponse, OnError);
	return true;
}

bool UServerSideOnlineSubsystem::Shutdown(FShutdownDelegate OnResponse, FAgonesErrorDelegate OnError)
{
	UAgonesSubsystem* AgonesSDK = UAgonesSubsystem::Get(this);
	if (!AgonesSDK)
	{
		return false;
	}
	AgonesSDK->Shutdown(OnResponse, OnError);
	return true;
}

bool UServerSideOnlineSubsystem::GetPlayer1(FUserAnnotations& PlayerData)
{
	if (MatchData.IsSet())
	{
		PlayerData = MatchData.GetValue().Player1;
		return true;
	}
	return false;
}

bool UServerSideOnlineSubsystem::GetPlayer2(FUserAnnotations& PlayerData)
{
	if (MatchData.IsSet())
	{
		PlayerData = MatchData.GetValue().Player2;
		return true;
	}
	return false;
}

bool UServerSideOnlineSubsystem::GetPlayerData(int64 PlayerId, FUserAnnotations& PlayerData)
{
	if (!MatchData.IsSet())
	{
		return false;
	}
	if (PlayerId == MatchData.GetValue().Player1.PlayerId)
	{
		PlayerData = MatchData.GetValue().Player1;
		return true;
	}
	if (PlayerId == MatchData.GetValue().Player2.PlayerId)
	{
		PlayerData = MatchData.GetValue().Player2;
		return true;
	}
	return false;
}

bool UServerSideOnlineSubsystem::DrawMatch(const TArray<FRoundData>& Rounds, FOnMatchEnd OnResponse, FOnErroneousResponse OnError)
{
	if (!MsClient.IsSet() || !MatchData.IsSet())
	{
		return false;
	}
	FEndMatchResult Request;
	Request.MatchId = MatchData.GetValue().MatchId;
	Request.Rounds = Rounds;
	WithRetry<FEndMatchResponse>([this, Request]()
	{
		return MsClient->EndMatch(Request);
	}, 3).Next([OnResponse, OnError](TOptional<FEndMatchResponse> Response)
	{
		if (!Response.IsSet())
		{
			return OnGameThread(&decltype(OnError)::ExecuteIfBound, OnError, FOnlineSubsystemError{TEXT("EndMatch failed after 3 attempts"), EOnlineErrorType::NonCritical});
		}
		return OnGameThread(&decltype(OnResponse)::ExecuteIfBound, OnResponse, Response.GetValue());
	}).Next(FutureJoin).Next([](bool bWasBound)
	{
		if (!bWasBound)
		{
			UE_LOG(LogTemp, Error, TEXT("OnResponse or OnError for DrawMatch was not set"));
		}
	});
	return true;
}

bool UServerSideOnlineSubsystem::EndMatch(int64 WinnerId, const TArray<FRoundData>& Rounds, FOnMatchEnd OnResponse, FOnErroneousResponse OnError)
{
	if (!MsClient.IsSet() || !MatchData.IsSet())
	{
		return false;
	}
	FEndMatchResult Request;
	Request.MatchId = MatchData.GetValue().MatchId;
	Request.WinnerId = WinnerId;
	Request.Rounds = Rounds;
	WithRetry<FEndMatchResponse>([this, Request]()
	{
		return MsClient->EndMatch(Request);
	}, 3).Next([OnResponse, OnError](TOptional<FEndMatchResponse> Response)
	{
		if (!Response.IsSet())
		{
			return OnGameThread(&decltype(OnError)::ExecuteIfBound, OnError, FOnlineSubsystemError{TEXT("EndMatch failed after 3 attempts"), EOnlineErrorType::NonCritical});
		}
		return OnGameThread(&decltype(OnResponse)::ExecuteIfBound, OnResponse, Response.GetValue());
	}).Next(FutureJoin).Next([](bool bWasBound)
	{
		if (!bWasBound)
		{
			UE_LOG(LogTemp, Error, TEXT("OnResponse or OnError for EndMatch was not set"));
		}
	});
	return true;
}

bool UServerSideOnlineSubsystem::CancelMatch(FOnEmptyResponse OnResponse, FOnErroneousResponse OnError)
{
	if (!MsClient.IsSet() || !MatchData.IsSet())
	{
		return false;
	}
	FString MatchId = MatchData.GetValue().MatchId;
	WithRetry<bool>([this, MatchId]()
	{
		return MsClient->CancelMatch(MatchId);
	}, 3).Next([OnResponse, OnError](TOptional<bool> Result)
	{
		if (!Result.IsSet() || !Result.GetValue())
		{
			return OnGameThread(&decltype(OnError)::ExecuteIfBound, OnError, FOnlineSubsystemError{TEXT("CancelMatch failed after 3 attempts"), EOnlineErrorType::NonCritical});
		}
		return OnGameThread(&decltype(OnResponse)::ExecuteIfBound, OnResponse);
	}).Next(FutureJoin).Next([](bool bWasBound)
	{
		if (!bWasBound)
		{
			UE_LOG(LogTemp, Error, TEXT("OnResponse or OnError for CancelMatch was not set"));
		}
	});
	return true;
}

bool UServerSideOnlineSubsystem::RenewMatch(FOnEmptyResponse OnResponse, FOnErroneousResponse OnError)
{
	if (!MsClient.IsSet() || !MatchData.IsSet())
	{
		return false;
	}
	FString MatchId = MatchData.GetValue().MatchId;
	WithRetry<bool>([this, MatchId]()
	{
		return MsClient->RenewMatch(MatchId);
	}, 3).Next([OnResponse, OnError](TOptional<bool> Result)
	{
		if (!Result.IsSet() || !Result.GetValue())
		{
			return OnGameThread(&decltype(OnError)::ExecuteIfBound, OnError, FOnlineSubsystemError{TEXT("RenewMatch failed after 3 attempts"), EOnlineErrorType::NonCritical});
		}
		return OnGameThread(&decltype(OnResponse)::ExecuteIfBound, OnResponse);
	}).Next(FutureJoin).Next([](bool bWasBound)
	{
		if (!bWasBound)
		{
			UE_LOG(LogTemp, Error, TEXT("OnResponse or OnError for RenewMatch was not set"));
		}
	});
	return true;
}

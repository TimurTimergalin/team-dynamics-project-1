#include "UserEventClient.h"

#include "Dom/JsonObject.h"
#include "Serialization/JsonReader.h"
#include "Serialization/JsonSerializer.h"
#include "TagDuels/OnlineSubsystem/utils.h"

namespace
{
	enum class EventType : uint8
	{
		StatusChanged,
		ChallengeReceived,
		ChallengeAccepted,
		ChallengeDeclined,
		ChallengeCancelled,
		MatchStarted,
		Error,
	};

	struct Event
	{
		EventType Type;
		// StatusChanged
		TOptional<int64> UserId;
		TOptional<FString> UserName;
		TOptional<FString> MessageId;
		TOptional<FString> Address;
		TOptional<FString> ErrorMessage;
		TOptional<EUserStatus> Status;
	};

	TOptional<EUserStatus> ParseStatus(const FString& StatusStr)
	{
		if (StatusStr == TEXT("Online"))  return EUserStatus::Online;
		if (StatusStr == TEXT("Offline")) return EUserStatus::Offline;
		if (StatusStr == TEXT("InGame"))  return EUserStatus::InGame;
		return {};
	}

	TOptional<Event> ParseEvent(const TSharedPtr<FJsonObject>& JsonObject)
	{
		FString TypeStr;
		if (!JsonObject->TryGetStringField(TEXT("type"), TypeStr))
		{
			UE_LOG(LogTemp, Error, TEXT("UserEventClient: event missing 'type' field"));
			return {};
		}

		Event Result;
		const TSharedPtr<FJsonObject>* PayloadObject;
		bool HasPayload = JsonObject->TryGetObjectField(TEXT("payload"), PayloadObject);

		if (TypeStr == TEXT("StatusChanged"))
		{
			Result.Type = EventType::StatusChanged;
			if (HasPayload)
			{
				double Id;
				if ((*PayloadObject)->TryGetNumberField(TEXT("userId"), Id))
					Result.UserId = static_cast<int64>(Id);
				FString StatusStr;
				if ((*PayloadObject)->TryGetStringField(TEXT("newStatus"), StatusStr))
					Result.Status = ParseStatus(StatusStr);
			}
		}
		else if (TypeStr == TEXT("ChallengeReceived"))
		{
			Result.Type = EventType::ChallengeReceived;
			if (HasPayload)
			{
				double Id;
				if ((*PayloadObject)->TryGetNumberField(TEXT("userId"), Id))
					Result.UserId = static_cast<int64>(Id);
				FString Name;
				if ((*PayloadObject)->TryGetStringField(TEXT("userName"), Name))
					Result.UserName = MoveTemp(Name);
				FString MsgId;
				if ((*PayloadObject)->TryGetStringField(TEXT("messageId"), MsgId))
					Result.MessageId = MoveTemp(MsgId);
			}
		}
		else if (TypeStr == TEXT("ChallengeAccepted") || TypeStr == TEXT("MatchStarted"))
		{
			Result.Type = TypeStr == TEXT("ChallengeAccepted") ? EventType::ChallengeAccepted : EventType::MatchStarted;
			if (HasPayload)
			{
				FString Addr;
				if ((*PayloadObject)->TryGetStringField(TEXT("address"), Addr))
					Result.Address = MoveTemp(Addr);
			}
		}
		else if (TypeStr == TEXT("ChallengeDeclined"))
		{
			Result.Type = EventType::ChallengeDeclined;
		}
		else if (TypeStr == TEXT("ChallengeCancelled"))
		{
			Result.Type = EventType::ChallengeCancelled;
			if (HasPayload)
			{
				double Id;
				if ((*PayloadObject)->TryGetNumberField(TEXT("userId"), Id))
					Result.UserId = static_cast<int64>(Id);
			}
		}
		else if (TypeStr == TEXT("Error"))
		{
			Result.Type = EventType::Error;
			if (HasPayload)
			{
				FString Msg;
				if ((*PayloadObject)->TryGetStringField(TEXT("message"), Msg))
					Result.ErrorMessage = MoveTemp(Msg);
			}
		}
		else
		{
			UE_LOG(LogTemp, Error, TEXT("UserEventClient: unknown event type '%s'"), *TypeStr);
			return {};
		}

		return Result;
	}

	void HandleEvent(
		const Event& Ev,
		int64 UserId,
		TSharedRef<FOnStatusUpdated> OnStatusChanged,
		TSharedRef<FOnChallengeReceived> OnChallengeReceived,
		TSharedRef<FOnMatch> OnMatchStarted,
		TSharedRef<FOnEmptyResponse> OnChallengeDeclined,
		TSharedRef<FOnChallengeCancelled> OnChallengeCancelled,
		TSharedRef<FOnErroneousResponse> OnError
	)
	{
		switch (Ev.Type)
		{
		case EventType::StatusChanged:
			if (Ev.UserId.IsSet() && Ev.Status.IsSet())
			{
				FStatusUpdatedEvent StatusEvent;
				StatusEvent.UserId = Ev.UserId.GetValue();
				StatusEvent.Status = Ev.Status.GetValue();
				OnStatusChanged->ExecuteIfBound(StatusEvent);
			}
			break;
		case EventType::ChallengeReceived:
			if (Ev.UserId.IsSet() && Ev.UserName.IsSet() && Ev.MessageId.IsSet())
			{
				FChallengeReceivedEvent ChallengeEvent;
				ChallengeEvent.UserId = Ev.UserId.GetValue();
				ChallengeEvent.UserName = Ev.UserName.GetValue();
				ChallengeEvent.MessageId = Ev.MessageId.GetValue();
				OnChallengeReceived->ExecuteIfBound(ChallengeEvent);
			}
			break;
		case EventType::ChallengeAccepted:
		case EventType::MatchStarted:
			if (Ev.Address.IsSet())
			{
				OnMatchStarted->ExecuteIfBound(Ev.Address.GetValue() + FString::Printf(TEXT("?player_id=%lld"), UserId));
			}
			break;
		case EventType::ChallengeDeclined:
			OnChallengeDeclined->ExecuteIfBound();
			break;
		case EventType::ChallengeCancelled:
			if (Ev.UserId.IsSet())
			{
				FChallengeCancelledEvent CancelledEvent;
				CancelledEvent.UserId = Ev.UserId.GetValue();
				OnChallengeCancelled->ExecuteIfBound(CancelledEvent);
			}
			break;
		case EventType::Error:
			OnError->ExecuteIfBound(FOnlineSubsystemError{Ev.ErrorMessage.Get(TEXT("Unknown error")), EOnlineErrorType::NonCritical});
			break;
		}
	}
}

void UserEventClient::EstablishConnection(TSharedRef<FOnStatusUpdated> OnStatusChanged,
                                          TSharedRef<FOnChallengeReceived> OnChallengeReceived, TSharedRef<FOnMatch> OnMatchStarted,
                                          TSharedRef<FOnEmptyResponse> OnChallengeDeclined, TSharedRef<FOnChallengeCancelled> OnChallengeCancelled,
                                          TSharedRef<FOnErroneousResponse> OnError, int64 UserId_, const FString& AuthToken)
{
	UserId = UserId_;
	Url = FString::Printf(TEXT("ws://%s/events?playerId=%lld&authToken=%s"), *Address, UserId, *AuthToken);
	EstablishConnection(OnStatusChanged, OnChallengeReceived, OnMatchStarted, OnChallengeDeclined, OnChallengeCancelled, OnError);
}

void UserEventClient::Retry(const FString& Error, EOnlineErrorType Type, TSharedRef<FOnStatusUpdated> OnStatusChanged,
                            TSharedRef<FOnChallengeReceived> OnChallengeReceived, TSharedRef<FOnMatch> OnMatchStarted,
                            TSharedRef<FOnEmptyResponse> OnChallengeDeclined, TSharedRef<FOnChallengeCancelled> OnChallengeCancelled,
                            TSharedRef<FOnErroneousResponse> OnError)
{
	if (ConnectionRetries <= 0)
	{
		Connection = nullptr;
		AsyncTask(ENamedThreads::GameThread, [OnError, Error, Type]()
		{
			OnError->ExecuteIfBound(FOnlineSubsystemError{Error, Type});
		});
	} else
	{
		--ConnectionRetries;
		EstablishConnection(OnStatusChanged, OnChallengeReceived, OnMatchStarted, OnChallengeDeclined, OnChallengeCancelled, OnError);
	}
}

UserEventClient::UserEventClient(const FString& Address_) : Address(Address_)
{
}

bool UserEventClient::Close()
{
	if (Connection == nullptr)
	{
		return false;
	}
	if (!Connection->IsConnected())
	{
		return false;
	}
	Resolved = true;
	Connection->Close();
	return true;
}

bool UserEventClient::IsConnected()
{
	return Connection != nullptr && Connection->IsConnected();
}

void UserEventClient::EstablishConnection(TSharedRef<FOnStatusUpdated> OnStatusChanged,
                                          TSharedRef<FOnChallengeReceived> OnChallengeReceived, TSharedRef<FOnMatch> OnMatchStarted,
                                          TSharedRef<FOnEmptyResponse> OnChallengeDeclined, TSharedRef<FOnChallengeCancelled> OnChallengeCancelled,
                                          TSharedRef<FOnErroneousResponse> OnError)
{
	Connection = FWebSocketsModule::Get().CreateWebSocket(Url, TEXT(""));
	Connection->OnConnected().AddLambda([this]()
	{
		ConnectionRetries = InitialConnectionRetries;
	});
	Connection->OnConnectionError().AddLambda([this, OnStatusChanged, OnChallengeReceived, OnMatchStarted, OnChallengeDeclined, OnChallengeCancelled, OnError](const FString& Error)
	{
		Retry(Error, EOnlineErrorType::Critical, OnStatusChanged, OnChallengeReceived, OnMatchStarted, OnChallengeDeclined, OnChallengeCancelled, OnError);
	});
	Connection->OnClosed().AddLambda([this, OnStatusChanged, OnChallengeReceived, OnMatchStarted, OnChallengeDeclined, OnChallengeCancelled, OnError](int32 /*StatusCode*/, const FString& Reason, bool /*WasClean*/)
	{
		if (!Resolved)
		{
			Retry(TEXT("Connection closed unexpectedly: ") + Reason, EOnlineErrorType::Critical, OnStatusChanged, OnChallengeReceived, OnMatchStarted, OnChallengeDeclined, OnChallengeCancelled, OnError);
		} else {
			Resolved = false;
		}
	});
	Connection->OnMessage().AddLambda([this, OnStatusChanged, OnChallengeReceived, OnMatchStarted, OnChallengeDeclined, OnChallengeCancelled, OnError](const FString& Message)
	{
		TSharedPtr<FJsonObject> JsonObject;
		TSharedRef<TJsonReader<>> Reader = TJsonReaderFactory<>::Create(Message);
		if (!FJsonSerializer::Deserialize(Reader, JsonObject) || !JsonObject.IsValid())
		{
			UE_LOG(LogTemp, Error, TEXT("UserEventClient: failed to parse message: %s"), *Message);
			return;
		}

		const TArray<TSharedPtr<FJsonValue>>* EventsArray;
		if (!JsonObject->TryGetArrayField(TEXT("events"), EventsArray))
		{
			UE_LOG(LogTemp, Error, TEXT("UserEventClient: missing 'events' field"));
			return;
		}

		TArray<Event> Events;
		for (const TSharedPtr<FJsonValue>& EventValue : *EventsArray)
		{
			const TSharedPtr<FJsonObject>* EventObject;
			if (!EventValue->TryGetObject(EventObject))
			{
				continue;
			}
			TOptional<Event> Ev = ParseEvent(*EventObject);
			if (!Ev.IsSet())
			{
				continue;
			}
			Events.Add(MoveTemp(Ev.GetValue()));
		}
		if (Events.IsEmpty())
		{
			return;
		}
		AsyncTask(ENamedThreads::GameThread, [UserId = this->UserId, Events = MoveTemp(Events), OnStatusChanged, OnChallengeReceived, OnMatchStarted, OnChallengeDeclined, OnChallengeCancelled, OnError]()
		{
			for (const Event& Ev : Events)
			{
				HandleEvent(Ev, UserId, OnStatusChanged, OnChallengeReceived, OnMatchStarted, OnChallengeDeclined, OnChallengeCancelled, OnError);
			}
		});
	});
	AsyncTask(ENamedThreads::GameThread, [this]()
	{
		if (!Connection)
		{
			UE_LOG(LogTemp, Error, TEXT("UserEventClient: Connection is nullptr"));
			return;
		}
		Connection->Connect();
	});
}

bool UserEventClient::SendRequest(const FString& JsonStr)
{
	if (Connection == nullptr || !Connection->IsConnected())
	{
		return false;
	}
	Connection->Send(JsonStr);
	return true;
}

bool UserEventClient::Subscribe(TArray<int64> UserIds)
{
	FString IdsJson = TEXT("[");
	for (int32 i = 0; i < UserIds.Num(); i++)
	{
		if (i > 0) IdsJson += TEXT(",");
		IdsJson += FString::Printf(TEXT("%lld"), UserIds[i]);
	}
	IdsJson += TEXT("]");
	return SendRequest(FString::Printf(TEXT("{\"type\":\"Subscribe\",\"payload\":{\"users\":%s}}"), *IdsJson));
}

bool UserEventClient::Challenge(int64 TargetUserId)
{
	return SendRequest(FString::Printf(TEXT("{\"type\":\"Challenge\",\"payload\":{\"userId\":%lld}}"), TargetUserId));
}

bool UserEventClient::CancelChallenge()
{
	return SendRequest(TEXT("{\"type\":\"CancelChallenge\",\"payload\":{}}"));
}

bool UserEventClient::AcceptChallenge(FChallengeReceivedEvent Challenge)
{
	return SendRequest(FString::Printf(TEXT("{\"type\":\"AcceptChallenge\",\"payload\":{\"messageId\":\"%s\",\"userId\":%lld}}"),
		*Challenge.MessageId, Challenge.UserId));
}

bool UserEventClient::DeclineChallenge(FChallengeReceivedEvent Challenge)
{
	return SendRequest(FString::Printf(TEXT("{\"type\":\"DeclineChallenge\",\"payload\":{\"messageId\":\"%s\",\"userId\":%lld}}"),
		*Challenge.MessageId, Challenge.UserId));
}

bool UserEventClient::NotifyBusy()
{
	return SendRequest(TEXT("{\"type\":\"NotifyBusy\",\"payload\":{}}"));
}

bool UserEventClient::NotifyFree()
{
	return SendRequest(TEXT("{\"type\":\"NotifyFree\",\"payload\":{}}"));
}

TOptional<UserEventClient> CreateUserEventClient()
{
	FString Address;
	if (!GConfig)
	{
		UE_LOG(LogTemp, Error, TEXT("CreateUserEventClient: GConfig absent"));
		return {};
	}
	if (!GConfig->GetString(TEXT("OnlineSubsystemAddresses"), TEXT("UserEvent"), Address, GGameIni))
	{
		UE_LOG(LogTemp, Error, TEXT("CreateUserEventClient: OnlineSubsystemAddresses:UserEvent not set"));
		return {};
	}
	return UserEventClient(Address);
}

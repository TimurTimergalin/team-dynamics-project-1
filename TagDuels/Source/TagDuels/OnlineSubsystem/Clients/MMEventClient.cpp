#include "MMEventClient.h"

#include "HttpModule.h"

namespace
{
	enum class ResponseType : uint8{
		Registered,
		Unregistered,
		Match,
		Error,
	};

	struct Response
	{
		ResponseType Type;
		TOptional<FString> GameServerAddress;
		TOptional<FString> ErrorMessage;
	};

	TOptional<Response> ParseResponse(const FString& JsonString)
	{
		TSharedPtr<FJsonObject> JsonObject;
		TSharedRef<TJsonReader<>> Reader = TJsonReaderFactory<>::Create(JsonString);
		if (!FJsonSerializer::Deserialize(Reader, JsonObject) || !JsonObject.IsValid())
		{
			UE_LOG(LogTemp, Error, TEXT("ParseResponse: failed to parse JSON: %s"), *JsonString);
			return {};
		}

		FString TypeStr;
		if (!JsonObject->TryGetStringField(TEXT("type"), TypeStr))
		{
			UE_LOG(LogTemp, Error, TEXT("ParseResponse: missing 'type' field"));
			return {};
		}

		ResponseType Type;
		if (TypeStr == TEXT("Registered"))
			Type = ResponseType::Registered;
		else if (TypeStr == TEXT("Unregistered"))
			Type = ResponseType::Unregistered;
		else if (TypeStr == TEXT("Match"))
			Type = ResponseType::Match;
		else if (TypeStr == TEXT("Error"))
			Type = ResponseType::Error;
		else
		{
			UE_LOG(LogTemp, Error, TEXT("ParseResponse: unknown response type '%s'"), *TypeStr);
			return {};
		}

		Response Result;
		Result.Type = Type;

		FString Address;
		if (JsonObject->TryGetStringField(TEXT("address"), Address))
		{
			Result.GameServerAddress = MoveTemp(Address);
		}

		FString ErrorMsg;
		if (JsonObject->TryGetStringField(TEXT("errorMessage"), ErrorMsg))
		{
			Result.ErrorMessage = MoveTemp(ErrorMsg);
		}

		return Result;
	}
}

MMEventClient::MMEventClient(const FString& Address): Address(Address)
{
	Fleet = "Moscow"; // Пока хардкод
}

void MMEventClient::ExecuteOnError(const FString& Error, EOnlineErrorType Type, TSharedRef<FOnErroneousResponse> OnError)
{
	AsyncTask(ENamedThreads::GameThread, [this, Error, Type, OnError]()
	{
		if (!OnError.Get().ExecuteIfBound(FOnlineSubsystemError{Error, Type}))
		{
			UE_LOG(LogTemp, Error, TEXT("OnError for MMEvent not set"));
		}
	});
}

void MMEventClient::ExecuteOnMatch(const FString& GameServerAddress, TSharedRef<FOnMatch> OnMatchCallback)
{
	AsyncTask(ENamedThreads::GameThread, [this, GameServerAddress, OnMatchCallback]()
	{
		if (!OnMatchCallback.Get().ExecuteIfBound(GameServerAddress + FString::Printf(TEXT("?player_id=%lld"), UserId)))
		{
			UE_LOG(LogTemp, Error, TEXT("OnMatchCallback for MMEvent not set"));
		}
	});
}

void MMEventClient::Retry(const FString& Error, EOnlineErrorType Type, TSharedRef<FOnMatch> OnResponse, TSharedRef<FOnErroneousResponse> OnError)
{
	if (ConnectionRetries <= 0)
	{
		Connection = nullptr;
		Resolved = false;
		ExecuteOnError(Error, Type, OnError);
	}
	else
	{
		--ConnectionRetries;
		EstablishConnection(OnResponse, OnError);
	}
}

void MMEventClient::EstablishConnection(TSharedRef<FOnMatch> OnResponse, TSharedRef<FOnErroneousResponse> OnError, int64 UserId_)
{
	UserId = UserId_;
	Url = FString::Printf(TEXT("ws://%s/events?playerId=%lld&fleet=%s"), *Address, UserId, *Fleet);
	EstablishConnection(OnResponse, OnError);
}

void MMEventClient::EstablishConnection(TSharedRef<FOnMatch> OnResponse, TSharedRef<FOnErroneousResponse> OnError)
{
	Connection = FWebSocketsModule::Get().CreateWebSocket(Url, "");
	Connection->OnConnected().AddLambda([this]()
	{
		ConnectionRetries = InitialConnectionRetries;
	});
	Connection->OnConnectionError().AddLambda([this, OnResponse, OnError](const FString& Error)
	{
		Retry(Error, EOnlineErrorType::Critical, OnResponse, OnError);
	});
	Connection->OnClosed().AddLambda([this, OnResponse, OnError](int32 /*StatusCode*/, const FString& Reason, bool /*WasClean*/)
	{
		if (!Resolved)
		{
			Retry("Connection closed before receiving match: " + Reason, EOnlineErrorType::Critical, OnResponse, OnError);
		} else {
			Resolved = false;
		}
	});
	Connection->OnMessage().AddLambda([this, OnResponse, OnError](const FString& message)
	{
		TOptional<Response> RespOpt = ParseResponse(message);
		if (!RespOpt)
		{
			Resolved = true;
			ExecuteOnError(FString::Printf(TEXT("Cannot parse response: %s"), *message), EOnlineErrorType::NonCritical, OnError);
			return;	
		}
		const Response& Resp = RespOpt.GetValue();
		if (Resp.GameServerAddress.IsSet())
		{
			Resolved = true;
			ExecuteOnMatch(Resp.GameServerAddress.GetValue(), OnResponse);
			return;
		}
		if (Resp.ErrorMessage.IsSet())
		{
			Resolved = true;
			ExecuteOnError(Resp.ErrorMessage.GetValue(), EOnlineErrorType::NonCritical, OnError);
			return;
		}
		if (Resp.Type == ResponseType::Unregistered)
		{
			Resolved = false;
		}
	});
	AsyncTask(ENamedThreads::GameThread, [this]()
	{
		if (!Connection) {
			UE_LOG(LogTemp, Error, TEXT("MME Connection is nullptr - that is Wierd!!!"));
			return;
		}
		Connection->Connect();
	});
}

bool MMEventClient::Close()
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

bool MMEventClient::StartMatchmaking()
{
	if (Connection == nullptr)
	{
		return false;
	}
	if (!Connection->IsConnected())
	{
		return false;
	}
	if (Resolved)
	{
		return false;
	}
	Connection->Send("{\"type\": \"Register\"}");
	return true;
}

bool MMEventClient::CancelMatchmaking()
{
	if (Connection == nullptr)
	{
		return false;
	}
	if (!Connection->IsConnected())
	{
		return false;
	}
	Connection->Send("{\"type\": \"Unregister\"}");
	return true;
}

bool MMEventClient::IsConnected()
{
	return Connection != nullptr && Connection->IsConnected();
}

TOptional<MMEventClient> CreateMMEventClient()
{
	FString Address;
	if (!GConfig)
	{
		UE_LOG(LogTemp, Error, TEXT("CreateRatingServiceClient: GConfig absent"));
		return {};
	}
	if (!GConfig->GetString(TEXT("OnlineSubsystemAddresses"), TEXT("MMEvent"), Address, GGameIni))
	{
		UE_LOG(LogTemp, Error, TEXT("CreateRatingServiceClient: OnlineSubsystemAddresses:MMEvent not set"));
		return {};
	}
	return MMEventClient(Address);
}

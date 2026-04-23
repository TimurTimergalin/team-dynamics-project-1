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

MMEventClient::MMEventClient(const FString& Address, int64 UserId, FOnMatch OnMatchCallback,
                             FOnErroneousResponse OnError): OnMatchCallback(OnMatchCallback), OnError(OnError), UserId(UserId)
{
	FString Fleet = "Moscow"; // Пока хардкод
	Url = FString::Printf(TEXT("ws://%s/events?playerId=%lld&fleet=%s"), *Address, UserId, *Fleet);
	EstablishConnection();
}

void MMEventClient::ExecuteOnError(const FString& Error)
{
	AsyncTask(ENamedThreads::GameThread, [OnError = this->OnError, Error]()
	{
		if (!OnError.ExecuteIfBound(Error))
		{
			UE_LOG(LogTemp, Error, TEXT("OnError for MMEvent not set"));
		}
	});
}

void MMEventClient::ExecuteOnMatch(const FString& GameServerAddress)
{
	AsyncTask(ENamedThreads::GameThread, [OnMatchCallback = this->OnMatchCallback, UserId = this->UserId, GameServerAddress]()
	{
		if (!OnMatchCallback.ExecuteIfBound(GameServerAddress + FString::Printf(TEXT("player_id=%lld"), UserId)))
		{
			UE_LOG(LogTemp, Error, TEXT("OnMatchCallback for MMEvent not set"));
		}
	});
}

void MMEventClient::Retry(const FString& Error)
{
	if (ConnectionRetries <= 0)
	{
		ExecuteOnError(Error);
	}
	else
	{
		--ConnectionRetries;
		EstablishConnection();
	}
}

void MMEventClient::EstablishConnection()
{
	Connection = FWebSocketsModule::Get().CreateWebSocket(Url, "");
	Connection->OnConnected().AddLambda([this]()
	{
		ConnectionRetries = InitialConnectionRetries;
	});
	Connection->OnConnectionError().AddLambda([this](const FString& Error)
	{
		Retry(Error);
	});
	Connection->OnClosed().AddLambda([this](int32 /*StatusCode*/, const FString& Reason, bool /*WasClean*/)
	{
		if (!Resolved)
		{
			Retry("Connection closed before receiving match: " + Reason);
		}
		Connection = nullptr;
	});
	Connection->OnMessage().AddLambda([this](const FString& message)
	{
		TOptional<Response> RespOpt = ParseResponse(message);
		if (!RespOpt)
		{
			Resolved = true;
			ExecuteOnError(FString::Printf(TEXT("Cannot parse response: %s"), *message));
			return;
		}
		const Response& Resp = RespOpt.GetValue();
		if (Resp.GameServerAddress.IsSet())
		{
			Resolved = true;
			ExecuteOnMatch(Resp.GameServerAddress.GetValue());
			return;
		}
		if (Resp.ErrorMessage.IsSet())
		{
			Resolved = true;
			ExecuteOnError(Resp.ErrorMessage.GetValue());
			return;
		}
		if (Resp.Type == ResponseType::Unregistered)
		{
			Resolved = false;
		}
	});
	Connection->Connect();
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
		return true;
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

TOptional<MMEventClient> CreateMMEventClient(int64 UserId, FOnMatch OnMatchCallback, FOnErroneousResponse OnError)
{
	FString Address;
	if (!GConfig)
	{
		UE_LOG(LogTemp, Error, TEXT("CreateRatingServiceClient: GConfig absent"));
		return {};
	}
	if (!GConfig->GetString(TEXT("/Script/OnlineServices.Addresses"), TEXT("MMEvent"), Address, GGameIni))
	{
		UE_LOG(LogTemp, Error, TEXT("CreateRatingServiceClient: /Script/OnlineServices.Addresses:MMEvent not set"));
		return {};
	}
	return MMEventClient(Address, UserId, OnMatchCallback, OnError);
}

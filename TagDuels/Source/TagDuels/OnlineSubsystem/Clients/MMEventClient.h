#pragma once

#include "WebSocketsModule.h"
#include "IWebSocket.h"
#include "TagDuels/OnlineSubsystem/Delegates.h"

class MMEventClient
{
public:
	void EstablishConnection();
	bool Close();
	bool StartMatchmaking();
	bool CancelMatchmaking();
	bool IsConnected();
private:
	void ExecuteOnError(const FString& Error, EOnlineErrorType Type);
	void ExecuteOnMatch(const FString& GameServerAddress);
	void Retry(const FString& Error,EOnlineErrorType Type);
	MMEventClient(const FString& Address, int64 UserId, FOnMatch OnMatchCallback, FOnErroneousResponse OnError);
	
	static constexpr int InitialConnectionRetries = 3; 
	
	FOnMatch OnMatchCallback;
	FOnErroneousResponse OnError;
	FString Url;
	bool Resolved{};
	int64 UserId{};
	TSharedPtr<IWebSocket> Connection;
	int ConnectionRetries = InitialConnectionRetries;
	TSharedRef<bool> Alive{MakeShared<bool>(true)};
	friend TOptional<MMEventClient> CreateMMEventClient(int64 UserId, FOnMatch OnMatchCallback, FOnErroneousResponse OnError);
};

TOptional<MMEventClient> CreateMMEventClient(int64 UserId, FOnMatch OnMatchCallback, FOnErroneousResponse OnError);

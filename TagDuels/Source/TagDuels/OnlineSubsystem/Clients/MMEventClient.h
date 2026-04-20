#pragma once

#include "WebSocketsModule.h"
#include "IWebSocket.h"
#include "TagDuels/OnlineSubsystem/Delegates.h"

class MMEventClient
{
public:
	void EstablishConnection();
	void ExecuteOnError(const FString& Error);
	void ExecuteOnMatch(const FString& GameServerAddress);
	void Retry(const FString& Error);
	bool Close();
	bool StartMatchmaking();
	bool CancelMatchmaking();
	bool IsConnected();
private:
	MMEventClient(const FString& Address, int64 UserId, FOnMatch OnMatchCallback, FOnErroneousResponse OnError);
	
	static constexpr int InitialConnectionRetries = 3; 
	
	FOnMatch OnMatchCallback;
	FOnErroneousResponse OnError;
	FString Url;
	bool Resolved{};
	TSharedPtr<IWebSocket> Connection;
	int ConnectionRetries = InitialConnectionRetries;
	friend TOptional<MMEventClient> CreateMMEventClient(int64 UserId, FOnMatch OnMatchCallback, FOnErroneousResponse OnError);
};

TOptional<MMEventClient> CreateMMEventClient(int64 UserId, FOnMatch OnMatchCallback, FOnErroneousResponse OnError);

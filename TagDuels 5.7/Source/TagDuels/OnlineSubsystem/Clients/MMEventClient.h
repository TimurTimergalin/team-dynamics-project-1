#pragma once

#include "WebSocketsModule.h"
#include "IWebSocket.h"
#include "TagDuels/OnlineSubsystem/Delegates.h"

class MMEventClient
{
public:
	void EstablishConnection(TSharedRef<FOnMatch> OnResponse, TSharedRef<FOnErroneousResponse> OnError, int64 UserId);
	bool Close();
	bool StartMatchmaking();
	bool CancelMatchmaking();
	bool IsConnected();
private:
	void EstablishConnection(TSharedRef<FOnMatch> OnResponse, TSharedRef<FOnErroneousResponse> OnError);
	void ExecuteOnError(const FString& Error, EOnlineErrorType Type, TSharedRef<FOnErroneousResponse> OnError);
	void ExecuteOnMatch(const FString& GameServerAddress, TSharedRef<FOnMatch> OnMatchCallback);
	void Retry(const FString& Error, EOnlineErrorType Type, TSharedRef<FOnMatch> OnResponse, TSharedRef<FOnErroneousResponse> OnError);
	explicit MMEventClient(const FString& Address);
	
	static constexpr int InitialConnectionRetries = 3; 
	
	FString Url;
	bool Resolved{};
	int64 UserId{};
	FString Address;
	TSharedPtr<IWebSocket> Connection;
	int ConnectionRetries = InitialConnectionRetries;
	FString Fleet;
	friend TOptional<MMEventClient> CreateMMEventClient();
};

TOptional<MMEventClient> CreateMMEventClient();

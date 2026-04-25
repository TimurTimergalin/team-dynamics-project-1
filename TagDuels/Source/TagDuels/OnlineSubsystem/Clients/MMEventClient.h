#pragma once

#include "WebSocketsModule.h"
#include "IWebSocket.h"
#include "TagDuels/OnlineSubsystem/Delegates.h"

class MMEventClient
{
public:
	void EstablishConnection(TSharedRef<FOnMatch> OnResponse, TSharedRef<FOnErroneousResponse> OnError, int64 UserId);
	bool Close() const;
	bool StartMatchmaking() const;
	bool CancelMatchmaking() const;
	bool IsConnected() const;
private:
	void EstablishConnection(TSharedRef<FOnMatch> OnResponse, TSharedRef<FOnErroneousResponse> OnError);
	void ExecuteOnError(const FString& Error, EOnlineErrorType Type, TSharedRef<FOnErroneousResponse> OnError) const;
	void ExecuteOnMatch(const FString& GameServerAddress, TSharedRef<FOnMatch> OnMatchCallback) const;
	void Retry(const FString& Error, EOnlineErrorType Type, TSharedRef<FOnMatch> OnResponse, TSharedRef<FOnErroneousResponse> OnError);
	MMEventClient(const FString& Address);
	
	static constexpr int InitialConnectionRetries = 3; 
	
	FString Url;
	bool Resolved{};
	int64 UserId{};
	FString Address;
	TSharedPtr<IWebSocket> Connection;
	int ConnectionRetries = InitialConnectionRetries;
	FString ConnectionId;
	FString Fleet;
	friend TOptional<MMEventClient> CreateMMEventClient();
};

TOptional<MMEventClient> CreateMMEventClient();

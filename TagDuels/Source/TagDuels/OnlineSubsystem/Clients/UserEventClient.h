#pragma once

#include "WebSocketsModule.h"
#include "IWebSocket.h"
#include "TagDuels/OnlineSubsystem/Delegates.h"

class UserEventClient
{
public:
	void EstablishConnection(
		TSharedRef<FOnStatusUpdated> OnStatusChanged,
		TSharedRef<FOnChallengeReceived> OnChallengeReceived,
		TSharedRef<FOnMatch> OnMatchStarted,
		TSharedRef<FOnEmptyResponse> OnChallengeDeclined,
		TSharedRef<FOnEmptyResponse> OnChallengeCancelled,
		TSharedRef<FOnErroneousResponse> OnError,
		int64 UserId_
	);
	
	bool Subscribe(TArray<int64> UserIds);
	bool Challenge(int64 UserId);
	bool CancelChallenge();
	bool AcceptChallenge(FChallengeReceivedEvent Challenge);
	bool DeclineChallenge(FChallengeReceivedEvent Challenge);
	bool NotifyBusy();
	bool NotifyFree();
	bool Close();
	bool IsConnected();
private:
	static constexpr int InitialConnectionRetries = 5;

	bool SendRequest(const FString& JsonStr);
	void EstablishConnection(
		TSharedRef<FOnStatusUpdated> OnStatusChanged,
		TSharedRef<FOnChallengeReceived> OnChallengeReceived,
		TSharedRef<FOnMatch> OnMatchStarted,
		TSharedRef<FOnEmptyResponse> OnChallengeDeclined,
		TSharedRef<FOnEmptyResponse> OnChallengeCancelled,
		TSharedRef<FOnErroneousResponse> OnError
	);
	void Retry(
		const FString& Error,
		EOnlineErrorType Type,
		TSharedRef<FOnStatusUpdated> OnStatusChanged,
		TSharedRef<FOnChallengeReceived> OnChallengeReceived,
		TSharedRef<FOnMatch> OnMatchStarted,
		TSharedRef<FOnEmptyResponse> OnChallengeDeclined,
		TSharedRef<FOnEmptyResponse> OnChallengeCancelled,
		TSharedRef<FOnErroneousResponse> OnError
	);
	explicit UserEventClient(const FString& Address_);

	FString Url;
	FString Address;
	int64 UserId{};
	TSharedPtr<IWebSocket> Connection;
	int ConnectionRetries = InitialConnectionRetries;
	bool Resolved = false;
	friend TOptional<UserEventClient> CreateUserEventClient();
};

TOptional<UserEventClient> CreateUserEventClient();

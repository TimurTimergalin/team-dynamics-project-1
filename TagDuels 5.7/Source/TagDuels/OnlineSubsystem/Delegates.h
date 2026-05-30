#pragma once
#include "Contract/Error.h"
#include "Contract/Match.h"
#include "Contract/MatchHistory.h"
#include "Contract/UserData.h"
#include "Contract/UserEvent.h"
#include "Delegates.generated.h"

USTRUCT()
struct FMyDelegateDummy
{
	GENERATED_BODY()
};

DECLARE_DYNAMIC_DELEGATE(FOnEmptyResponse);

DECLARE_DYNAMIC_DELEGATE_OneParam(FOnErroneousResponse, FOnlineSubsystemError, Error);

DECLARE_DYNAMIC_DELEGATE_OneParam(FOnInt64Response, int64, Response);

DECLARE_DYNAMIC_DELEGATE_OneParam(FOnMatchHistoryResponse, const FMatchHistoryPage&, Response);

DECLARE_DYNAMIC_DELEGATE_OneParam(FOnUserListResponse, const FPlayersList&, Response);

DECLARE_DYNAMIC_DELEGATE_OneParam(FOnMatch, FString, Address);

DECLARE_DYNAMIC_DELEGATE_OneParam(FOnMatchEnd, FEndMatchResponse, NewRatings);

DECLARE_DYNAMIC_DELEGATE_OneParam(FOnStatusUpdated, FStatusUpdatedEvent, Event);

DECLARE_DYNAMIC_DELEGATE_OneParam(FOnChallengeReceived, FChallengeReceivedEvent, Event);

DECLARE_DYNAMIC_DELEGATE_OneParam(FOnChallengeCancelled, FChallengeCancelledEvent, Event);

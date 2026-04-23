#pragma once
#include "Contract/Match.h"
#include "Contract/MatchHistory.h"
#include "Contract/UserData.h"
#include "Delegates.generated.h"

USTRUCT()
struct FMyDelegateDummy
{
	GENERATED_BODY()
};

DECLARE_DYNAMIC_DELEGATE(FOnEmptyResponse);

DECLARE_DYNAMIC_DELEGATE_OneParam(FOnErroneousResponse, FString, ErrorMessage);

DECLARE_DYNAMIC_DELEGATE_OneParam(FOnInt64Response, int64, Response);

DECLARE_DYNAMIC_DELEGATE_OneParam(FOnMatchHistoryResponse, const FMatchHistoryPage&, Response);

DECLARE_DYNAMIC_DELEGATE_OneParam(FOnUserListResponse, const FPlayersList&, Response);

DECLARE_DYNAMIC_DELEGATE_OneParam(FOnMatch, FString, Address);

DECLARE_DYNAMIC_DELEGATE_OneParam(FOnMatchEnd, FEndMatchResponse, NewRatings);

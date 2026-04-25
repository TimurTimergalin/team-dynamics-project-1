#pragma once

#include "CoreMinimal.h"
#include "AgonesSubsystem.h"

class TAgonesClient
{
public:
	TFuture<TOptional<FAgonesError>> Ready(UAgonesSubsystem& AgonesSDK);
	TFuture<TOptional<FAgonesError>> Shutdown(UAgonesSubsystem& AgonesSDK);
	TFuture<TOptional<FAgonesError>> NewMatchId(UAgonesSubsystem& AgonesSDK, const FString& MatchId);
};


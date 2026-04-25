#include "TAgonesClient.h"

#include "AgonesPromise.h"

TFuture<TOptional<FAgonesError>> TAgonesClient::Ready(UAgonesSubsystem& AgonesSDK)
{
	auto Promise = MakeShared<TPromise<TOptional<FAgonesError>>>();
	auto Future = Promise.Get().GetFuture();
	UAgonesPromise* PromiseObj = NewObject<UAgonesPromise>();
	PromiseObj->Promise = Promise;
	FReadyDelegate SuccessDelegate;
	SuccessDelegate.BindDynamic(PromiseObj, &UAgonesPromise::OnSuccess);
	FAgonesErrorDelegate AgonesErrorDelegate;
	AgonesErrorDelegate.BindDynamic(PromiseObj, &UAgonesPromise::OnError);
	AgonesSDK.Ready(SuccessDelegate, AgonesErrorDelegate);
	return Future;
}

TFuture<TOptional<FAgonesError>> TAgonesClient::Shutdown(UAgonesSubsystem& AgonesSDK)
{
	auto Promise = MakeShared<TPromise<TOptional<FAgonesError>>>();
	auto Future = Promise.Get().GetFuture();
	UAgonesPromise* PromiseObj = NewObject<UAgonesPromise>();
	PromiseObj->Promise = Promise;
	FShutdownDelegate SuccessDelegate;
	SuccessDelegate.BindDynamic(PromiseObj, &UAgonesPromise::OnSuccess);
	FAgonesErrorDelegate AgonesErrorDelegate;
	AgonesErrorDelegate.BindDynamic(PromiseObj, &UAgonesPromise::OnError);
	AgonesSDK.Shutdown(SuccessDelegate, AgonesErrorDelegate);
	return Future;
}

TFuture<TOptional<FAgonesError>> TAgonesClient::NewMatchId(UAgonesSubsystem& AgonesSDK, const FString& MatchId)
{
	auto Promise = MakeShared<TPromise<TOptional<FAgonesError>>>();
	auto Future = Promise.Get().GetFuture();
	UAgonesPromise* PromiseObj = NewObject<UAgonesPromise>();
	PromiseObj->Promise = Promise;
	FSetLabelDelegate SuccessDelegate;
	SuccessDelegate.BindDynamic(PromiseObj, &UAgonesPromise::OnSuccess);
	FAgonesErrorDelegate AgonesErrorDelegate;
	AgonesErrorDelegate.BindDynamic(PromiseObj, &UAgonesPromise::OnError);
	AgonesSDK.SetLabel("match-id", MatchId, SuccessDelegate, AgonesErrorDelegate);
	return Future;
}

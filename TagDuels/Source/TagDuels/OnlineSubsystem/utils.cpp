#include "utils.h"
#include "HttpModule.h"

TFuture<FHttpResponsePtr> MakeHttpRequest(TSharedPtr<IHttpRequest> Request)
{
	TSharedPtr<TPromise<FHttpResponsePtr>> Promise = MakeShared<TPromise<FHttpResponsePtr>>();
	TFuture<FHttpResponsePtr> future = Promise->GetFuture();

	Request->OnProcessRequestComplete().BindLambda([Promise](FHttpRequestPtr, FHttpResponsePtr response, bool bWasSuccessful) mutable
	{
		Promise->SetValue(bWasSuccessful && response.IsValid() ? response : nullptr);
	});
	Request->ProcessRequest();
	return future;
}
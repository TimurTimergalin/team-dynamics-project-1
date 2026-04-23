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

TOptional<int64> StrToInt64(const FString& InString)
{
	const TCHAR *Start = *InString;
	TCHAR* End = nullptr;
	int64 Value = FCString::Strtoi64(Start, &End, 10);
	if (*End != '\0')
	{
		return {};
	}
	return Value;
}

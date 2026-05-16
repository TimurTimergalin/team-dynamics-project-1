#include "utils.h"
#include "HttpModule.h"
#include "Interfaces/IHttpResponse.h"


TFuture<FHttpResponsePtr> MakeHttpRequest(TSharedPtr<IHttpRequest> Request)
{
	TSharedPtr<AutoFullfill<FHttpResponsePtr>> Promise = MakeShared<AutoFullfill<FHttpResponsePtr>>();
	TFuture<FHttpResponsePtr> future = Promise->Promise.GetFuture();

	Request->OnProcessRequestComplete().BindLambda([Promise](FHttpRequestPtr, FHttpResponsePtr Response, bool bWasSuccessful) mutable
	{
		Promise->Promise.SetValue(bWasSuccessful && Response.IsValid() ? Response : nullptr);
		if (bWasSuccessful && Response.IsValid()) {
			int32 StatusCode = Response->GetResponseCode();
			UE_LOG(LogTemp, Warning, TEXT("HTTP Status: %d"), StatusCode);

			// Log response body (as string)
			FString ResponseBody = Response->GetContentAsString();
			UE_LOG(LogTemp, Warning, TEXT("Response Body: %s"), *ResponseBody);
		}
		else
		{
			UE_LOG(LogTemp, Error, TEXT("HTTP request failed or invalid response"));
		}
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

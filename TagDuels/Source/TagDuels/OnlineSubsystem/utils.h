#pragma once

#include "CoreMinimal.h"
#include "Interfaces/IHttpRequest.h"
#include "Async/Future.h"

TFuture<FHttpResponsePtr> MakeHttpRequest(TSharedPtr<IHttpRequest> Request);

template <class T, class F>
TFuture<TOptional<T>> WithRetry(F Func, int Attempts)
{
	return Async(EAsyncExecution::ThreadPool, [Func, Attempts]()
	{
		for (int i = 0; i < Attempts; i++)
		{
			TOptional<T> Res = Func().Get();
			if (Res.IsSet())
			{
				return Res;
			}
		}
		return TOptional<T>();
	});
}

TOptional<int64> StrToInt64(const FString& InString);

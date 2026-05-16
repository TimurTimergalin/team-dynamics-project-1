#pragma once

#include "CoreMinimal.h"
#include "Interfaces/IHttpRequest.h"
#include "Async/Future.h"
#include <type_traits>
#include <functional>

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

template <typename F, typename... Args>
auto OnGameThread(F&& Callable, Args&&... ArgsToForward)
    -> TFuture<typename TInvokeResult<F, Args...>::Type>
{
	using ResultType = typename TInvokeResult<F, Args...>::Type;
	TPromise<ResultType> Promise;
	TFuture<ResultType> Future = Promise.GetFuture();

	AsyncTask(ENamedThreads::GameThread,
		[Callable = std::forward<F>(Callable),
		 ArgsTuple = std::make_tuple(std::forward<Args>(ArgsToForward)...),
		 Promise = std::move(Promise)]() mutable
		{
			if constexpr (std::is_void_v<ResultType>)
			{
				std::apply(Callable, ArgsTuple);
				Promise.SetValue();
			}
			else
			{
				ResultType Result = std::apply(Callable, ArgsTuple);
				Promise.SetValue(std::move(Result));
			}
		});

	return Future;
}

struct TFutureJoin {
	template<class T>
	T operator()(TFuture<T>&& future) {
		return future.Get();
	}
};

static constexpr TFutureJoin FutureJoin{};

template<class T>
struct AutoFullfill {
	TPromise<T> Promise{};
	~AutoFullfill() {
		T zero{};
		try {
			Promise.SetValue(zero);
		} catch (...) {
		}
	};
};

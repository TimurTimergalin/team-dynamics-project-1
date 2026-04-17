#pragma once

#include "CoreMinimal.h"
#include "Interfaces/IHttpRequest.h"
#include "interfaces/IHttpResponse.h"
#include "Async/Future.h"

TFuture<FHttpResponsePtr> MakeHttpRequest(TSharedPtr<IHttpRequest> Request);

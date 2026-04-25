// Fill out your copyright notice in the Description page of Project Settings.

#pragma once
#include "AgonesSubsystem.h"

#include "CoreMinimal.h"
#include "UObject/Object.h"
#include "AgonesPromise.generated.h"

/**
 * 
 */
UCLASS()
class TAGDUELS_API UAgonesPromise : public UObject
{
	GENERATED_BODY()
public:
	UFUNCTION()
	void OnSuccess(const FEmptyResponse& Resp);
	UFUNCTION()
	void OnError(const FAgonesError& Err);
	TSharedRef<TPromise<TOptional<FAgonesError>>> Promise;
};

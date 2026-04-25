// Fill out your copyright notice in the Description page of Project Settings.


#include "AgonesPromise.h"

void UAgonesPromise::OnSuccess(const FEmptyResponse&)
{
	Promise->SetValue({});
}

void UAgonesPromise::OnError(const FAgonesError& Err)
{
	Promise->SetValue(Err);
}

// Fill out your copyright notice in the Description page of Project Settings.

#pragma once

#include "CoreMinimal.h"
#include "GameFramework/PlayerController.h"
#include "TagDuelsPlayerController.generated.h"

/**
 * 
 */
UCLASS()
class TAGDUELS_API ATagDuelsPlayerController : public APlayerController
{
	GENERATED_BODY()
	
public:
    virtual void SetViewTarget(class AActor* NewViewTarget, FViewTargetTransitionParams TransitionParams = FViewTargetTransitionParams()) override;
};

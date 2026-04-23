// Fill out your copyright notice in the Description page of Project Settings.


#include "TagDuelsPlayerController.h"
#include "Camera/PlayerCameraManager.h"

void ATagDuelsPlayerController::SetViewTarget(AActor* NewViewTarget, FViewTargetTransitionParams TransitionParams)
{
	Super::SetViewTarget(NewViewTarget, TransitionParams);

	if (PlayerCameraManager)
	{
		//  PlayerCameraManager->SetManualCameraFade(1.0f, FColor::Black, false);
	}
}

void ATagDuelsPlayerController::BeginPlayingState()
{
	Super::BeginPlayingState();
	OnBeginPlayingState();
}

void ATagDuelsPlayerController::OnBeginPlayingState_Implementation()
{
}    

// Fill out your copyright notice in the Description page of Project Settings.


#include "TagDuelsPlayerController.h"
#include "Camera/PlayerCameraManager.h"

void ATagDuelsPlayerController::SetViewTarget(AActor* NewViewTarget, FViewTargetTransitionParams TransitionParams)
{
	Super::SetViewTarget(NewViewTarget, TransitionParams);

	if (PlayerCameraManager)
	{
		// (Не работает) Черный экран, пока грузится уровень
		//  PlayerCameraManager->SetManualCameraFade(1.0f, FColor::Black, false);
	}
}

// Fill out your copyright notice in the Description page of Project Settings.

#pragma once

#include "CoreMinimal.h"
#include "TagDuelsEnums.generated.h"

UENUM(BlueprintType)
enum class EOnlineSubsystemType : uint8
{
	None            UMETA(DisplayName = "None"),
	Null            UMETA(DisplayName = "NULL (LAN)"),
	Steam           UMETA(DisplayName = "Steam"),
	EOS             UMETA(DisplayName = "Epic Online Services"),
	// Добавьте другие по необходимости
};

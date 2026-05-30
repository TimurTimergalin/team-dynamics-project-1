// Fill out your copyright notice in the Description page of Project Settings.


#include "TagDuelsGameUserSettings.h"

UTagDuelsGameUserSettings* UTagDuelsGameUserSettings::GetTagDuelsGameUserSettings()
{
	// Используем стандартный метод GameUserSettings, который уже возвращает правильный объект
	UGameUserSettings* BaseSettings = UGameUserSettings::GetGameUserSettings();
	if (!BaseSettings)
	{
		UE_LOG(LogTemp, Error, TEXT("GetCustomSettings: GameUserSettings is null!"));
		return nullptr;
	}

	// Приводим к нашему типу
	UTagDuelsGameUserSettings* CustomSettings = Cast<UTagDuelsGameUserSettings>(BaseSettings);
	if (!CustomSettings)
	{
		// Это может произойти, если в Project Settings не выбран наш класс
		UE_LOG(LogTemp, Error, TEXT("GetCustomSettings: Failed to cast to UCustom_Settings. Did you set Game User Settings Class in Project Settings?"));
	}

	return CustomSettings;
}

float UTagDuelsGameUserSettings::GetMasterVolume() const
{
	return MasterVolume;
}

void UTagDuelsGameUserSettings::SetMasterVolume(float Volume)
{
	MasterVolume = Volume;
}

float UTagDuelsGameUserSettings::GetMusicVolume() const
{
	return MusicVolume;
}

void UTagDuelsGameUserSettings::SetMusicVolume(float Volume)
{
	MusicVolume = Volume;
}

float UTagDuelsGameUserSettings::GetSFXVolume() const
{
	return SFXVolume;
}

void UTagDuelsGameUserSettings::SetSFXVolume(float Volume)
{
	SFXVolume = Volume;
}

// Fill out your copyright notice in the Description page of Project Settings.


#include "TagDuelsGameUserSettings.h"

float UTagDuelsGameUserSettings::GetMasterVolume()
{
	return MasterVolume;
}

void UTagDuelsGameUserSettings::SetMasterVolume(float Volume)
{
	MasterVolume = Volume;
}

float UTagDuelsGameUserSettings::GetMusicVolume()
{
	return MusicVolume;
}

void UTagDuelsGameUserSettings::SetMusicVolume(float Volume)
{
	MusicVolume = Volume;
}

float UTagDuelsGameUserSettings::GetSFXVolume()
{
	return SFXVolume;
}

void UTagDuelsGameUserSettings::SetSFXVolume(float Volume)
{
	SFXVolume = Volume;
}

// Fill out your copyright notice in the Description page of Project Settings.

using UnrealBuildTool;
using System.Collections.Generic;

public class ThirdPerson_BP_UE5_7ClientTarget : TargetRules
{
	public ThirdPerson_BP_UE5_7ClientTarget(TargetInfo Target) : base(Target)
	{
		Type = TargetType.Client;
		DefaultBuildSettings = BuildSettingsVersion.V6;

		ExtraModuleNames.AddRange( new string[] { "ThirdPerson_BP_UE5_7" } );
	}
}

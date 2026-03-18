// Fill out your copyright notice in the Description page of Project Settings.

using UnrealBuildTool;
using System.Collections.Generic;

public class ThirdPerson_BP_UE5_7EditorTarget : TargetRules
{
	public ThirdPerson_BP_UE5_7EditorTarget(TargetInfo Target) : base(Target)
	{
		Type = TargetType.Editor;
		DefaultBuildSettings = BuildSettingsVersion.V6;

		ExtraModuleNames.AddRange( new string[] { "ThirdPerson_BP_UE5_7" } );
	}
}

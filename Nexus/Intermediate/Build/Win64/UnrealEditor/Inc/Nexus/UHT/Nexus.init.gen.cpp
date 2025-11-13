// Copyright Epic Games, Inc. All Rights Reserved.
/*===========================================================================
	Generated code exported from UnrealHeaderTool.
	DO NOT modify this manually! Edit the corresponding .h files instead!
===========================================================================*/

#include "UObject/GeneratedCppIncludes.h"
PRAGMA_DISABLE_DEPRECATION_WARNINGS
void EmptyLinkFunctionForGeneratedCodeNexus_init() {}
	static FPackageRegistrationInfo Z_Registration_Info_UPackage__Script_Nexus;
	FORCENOINLINE UPackage* Z_Construct_UPackage__Script_Nexus()
	{
		if (!Z_Registration_Info_UPackage__Script_Nexus.OuterSingleton)
		{
			static const UECodeGen_Private::FPackageParams PackageParams = {
				"/Script/Nexus",
				nullptr,
				0,
				PKG_CompiledIn | 0x00000000,
				0x08A02449,
				0x85A44A08,
				METADATA_PARAMS(0, nullptr)
			};
			UECodeGen_Private::ConstructUPackage(Z_Registration_Info_UPackage__Script_Nexus.OuterSingleton, PackageParams);
		}
		return Z_Registration_Info_UPackage__Script_Nexus.OuterSingleton;
	}
	static FRegisterCompiledInInfo Z_CompiledInDeferPackage_UPackage__Script_Nexus(Z_Construct_UPackage__Script_Nexus, TEXT("/Script/Nexus"), Z_Registration_Info_UPackage__Script_Nexus, CONSTRUCT_RELOAD_VERSION_INFO(FPackageReloadVersionInfo, 0x08A02449, 0x85A44A08));
PRAGMA_ENABLE_DEPRECATION_WARNINGS

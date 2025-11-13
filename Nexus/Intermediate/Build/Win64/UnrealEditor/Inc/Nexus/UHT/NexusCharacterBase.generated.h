// Copyright Epic Games, Inc. All Rights Reserved.
/*===========================================================================
	Generated code exported from UnrealHeaderTool.
	DO NOT modify this manually! Edit the corresponding .h files instead!
===========================================================================*/

// IWYU pragma: private, include "GameplayAbilitySystem/Character/NexusCharacterBase.h"

#ifdef NEXUS_NexusCharacterBase_generated_h
#error "NexusCharacterBase.generated.h already included, missing '#pragma once' in NexusCharacterBase.h"
#endif
#define NEXUS_NexusCharacterBase_generated_h

#include "UObject/ObjectMacros.h"
#include "UObject/ScriptMacros.h"

PRAGMA_DISABLE_DEPRECATION_WARNINGS

// ********** Begin Class ANexusCharacterBase ******************************************************
NEXUS_API UClass* Z_Construct_UClass_ANexusCharacterBase_NoRegister();

#define FID_UnrealProjects_team_dynamics_project_1_Nexus_Source_Nexus_GameplayAbilitySystem_Character_NexusCharacterBase_h_14_INCLASS_NO_PURE_DECLS \
private: \
	static void StaticRegisterNativesANexusCharacterBase(); \
	friend struct Z_Construct_UClass_ANexusCharacterBase_Statics; \
	static UClass* GetPrivateStaticClass(); \
	friend NEXUS_API UClass* Z_Construct_UClass_ANexusCharacterBase_NoRegister(); \
public: \
	DECLARE_CLASS2(ANexusCharacterBase, ACharacter, COMPILED_IN_FLAGS(0 | CLASS_Config), CASTCLASS_None, TEXT("/Script/Nexus"), Z_Construct_UClass_ANexusCharacterBase_NoRegister) \
	DECLARE_SERIALIZER(ANexusCharacterBase) \
	virtual UObject* _getUObject() const override { return const_cast<ANexusCharacterBase*>(this); }


#define FID_UnrealProjects_team_dynamics_project_1_Nexus_Source_Nexus_GameplayAbilitySystem_Character_NexusCharacterBase_h_14_ENHANCED_CONSTRUCTORS \
	/** Deleted move- and copy-constructors, should never be used */ \
	ANexusCharacterBase(ANexusCharacterBase&&) = delete; \
	ANexusCharacterBase(const ANexusCharacterBase&) = delete; \
	DECLARE_VTABLE_PTR_HELPER_CTOR(NO_API, ANexusCharacterBase); \
	DEFINE_VTABLE_PTR_HELPER_CTOR_CALLER(ANexusCharacterBase); \
	DEFINE_DEFAULT_CONSTRUCTOR_CALL(ANexusCharacterBase) \
	NO_API virtual ~ANexusCharacterBase();


#define FID_UnrealProjects_team_dynamics_project_1_Nexus_Source_Nexus_GameplayAbilitySystem_Character_NexusCharacterBase_h_11_PROLOG
#define FID_UnrealProjects_team_dynamics_project_1_Nexus_Source_Nexus_GameplayAbilitySystem_Character_NexusCharacterBase_h_14_GENERATED_BODY \
PRAGMA_DISABLE_DEPRECATION_WARNINGS \
public: \
	FID_UnrealProjects_team_dynamics_project_1_Nexus_Source_Nexus_GameplayAbilitySystem_Character_NexusCharacterBase_h_14_INCLASS_NO_PURE_DECLS \
	FID_UnrealProjects_team_dynamics_project_1_Nexus_Source_Nexus_GameplayAbilitySystem_Character_NexusCharacterBase_h_14_ENHANCED_CONSTRUCTORS \
private: \
PRAGMA_ENABLE_DEPRECATION_WARNINGS


class ANexusCharacterBase;

// ********** End Class ANexusCharacterBase ********************************************************

#undef CURRENT_FILE_ID
#define CURRENT_FILE_ID FID_UnrealProjects_team_dynamics_project_1_Nexus_Source_Nexus_GameplayAbilitySystem_Character_NexusCharacterBase_h

PRAGMA_ENABLE_DEPRECATION_WARNINGS

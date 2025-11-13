// Copyright Epic Games, Inc. All Rights Reserved.
/*===========================================================================
	Generated code exported from UnrealHeaderTool.
	DO NOT modify this manually! Edit the corresponding .h files instead!
===========================================================================*/

// IWYU pragma: private, include "GameplayAbilitySystem/AttributeSets/BattleAttributeSet.h"

#ifdef NEXUS_BattleAttributeSet_generated_h
#error "BattleAttributeSet.generated.h already included, missing '#pragma once' in BattleAttributeSet.h"
#endif
#define NEXUS_BattleAttributeSet_generated_h

#include "UObject/ObjectMacros.h"
#include "UObject/ScriptMacros.h"
#include "Net/Core/PushModel/PushModelMacros.h"

PRAGMA_DISABLE_DEPRECATION_WARNINGS

struct FGameplayAttributeData;

// ********** Begin Class UBattleAttributeSet ******************************************************
#define FID_UnrealProjects_team_dynamics_project_1_Nexus_Source_Nexus_GameplayAbilitySystem_AttributeSets_BattleAttributeSet_h_16_RPC_WRAPPERS_NO_PURE_DECLS \
	DECLARE_FUNCTION(execOnRep_HealingTime); \
	DECLARE_FUNCTION(execOnRep_HealingAmount); \
	DECLARE_FUNCTION(execOnRep_Damage);


NEXUS_API UClass* Z_Construct_UClass_UBattleAttributeSet_NoRegister();

#define FID_UnrealProjects_team_dynamics_project_1_Nexus_Source_Nexus_GameplayAbilitySystem_AttributeSets_BattleAttributeSet_h_16_INCLASS_NO_PURE_DECLS \
private: \
	static void StaticRegisterNativesUBattleAttributeSet(); \
	friend struct Z_Construct_UClass_UBattleAttributeSet_Statics; \
	static UClass* GetPrivateStaticClass(); \
	friend NEXUS_API UClass* Z_Construct_UClass_UBattleAttributeSet_NoRegister(); \
public: \
	DECLARE_CLASS2(UBattleAttributeSet, UAttributeSet, COMPILED_IN_FLAGS(0), CASTCLASS_None, TEXT("/Script/Nexus"), Z_Construct_UClass_UBattleAttributeSet_NoRegister) \
	DECLARE_SERIALIZER(UBattleAttributeSet) \
	enum class ENetFields_Private : uint16 \
	{ \
		NETFIELD_REP_START=(uint16)((int32)Super::ENetFields_Private::NETFIELD_REP_END + (int32)1), \
		Damage=NETFIELD_REP_START, \
		HealingAmount, \
		HealingTime, \
		NETFIELD_REP_END=HealingTime	}; \
	DECLARE_VALIDATE_GENERATED_REP_ENUMS(NO_API) \
private: \
	REPLICATED_BASE_CLASS(UBattleAttributeSet) \
public:


#define FID_UnrealProjects_team_dynamics_project_1_Nexus_Source_Nexus_GameplayAbilitySystem_AttributeSets_BattleAttributeSet_h_16_ENHANCED_CONSTRUCTORS \
	/** Deleted move- and copy-constructors, should never be used */ \
	UBattleAttributeSet(UBattleAttributeSet&&) = delete; \
	UBattleAttributeSet(const UBattleAttributeSet&) = delete; \
	DECLARE_VTABLE_PTR_HELPER_CTOR(NO_API, UBattleAttributeSet); \
	DEFINE_VTABLE_PTR_HELPER_CTOR_CALLER(UBattleAttributeSet); \
	DEFINE_DEFAULT_CONSTRUCTOR_CALL(UBattleAttributeSet) \
	NO_API virtual ~UBattleAttributeSet();


#define FID_UnrealProjects_team_dynamics_project_1_Nexus_Source_Nexus_GameplayAbilitySystem_AttributeSets_BattleAttributeSet_h_13_PROLOG
#define FID_UnrealProjects_team_dynamics_project_1_Nexus_Source_Nexus_GameplayAbilitySystem_AttributeSets_BattleAttributeSet_h_16_GENERATED_BODY \
PRAGMA_DISABLE_DEPRECATION_WARNINGS \
public: \
	FID_UnrealProjects_team_dynamics_project_1_Nexus_Source_Nexus_GameplayAbilitySystem_AttributeSets_BattleAttributeSet_h_16_RPC_WRAPPERS_NO_PURE_DECLS \
	FID_UnrealProjects_team_dynamics_project_1_Nexus_Source_Nexus_GameplayAbilitySystem_AttributeSets_BattleAttributeSet_h_16_INCLASS_NO_PURE_DECLS \
	FID_UnrealProjects_team_dynamics_project_1_Nexus_Source_Nexus_GameplayAbilitySystem_AttributeSets_BattleAttributeSet_h_16_ENHANCED_CONSTRUCTORS \
private: \
PRAGMA_ENABLE_DEPRECATION_WARNINGS


class UBattleAttributeSet;

// ********** End Class UBattleAttributeSet ********************************************************

#undef CURRENT_FILE_ID
#define CURRENT_FILE_ID FID_UnrealProjects_team_dynamics_project_1_Nexus_Source_Nexus_GameplayAbilitySystem_AttributeSets_BattleAttributeSet_h

PRAGMA_ENABLE_DEPRECATION_WARNINGS

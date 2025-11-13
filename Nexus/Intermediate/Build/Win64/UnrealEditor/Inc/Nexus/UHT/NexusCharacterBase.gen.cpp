// Copyright Epic Games, Inc. All Rights Reserved.
/*===========================================================================
	Generated code exported from UnrealHeaderTool.
	DO NOT modify this manually! Edit the corresponding .h files instead!
===========================================================================*/

#include "UObject/GeneratedCppIncludes.h"
#include "Nexus/GameplayAbilitySystem/Character/NexusCharacterBase.h"

PRAGMA_DISABLE_DEPRECATION_WARNINGS

void EmptyLinkFunctionForGeneratedCodeNexusCharacterBase() {}

// ********** Begin Cross Module References ********************************************************
ENGINE_API UClass* Z_Construct_UClass_ACharacter();
GAMEPLAYABILITIES_API UClass* Z_Construct_UClass_UAbilitySystemComponent_NoRegister();
GAMEPLAYABILITIES_API UClass* Z_Construct_UClass_UAbilitySystemInterface_NoRegister();
GAMEPLAYABILITIES_API UEnum* Z_Construct_UEnum_GameplayAbilities_EGameplayEffectReplicationMode();
NEXUS_API UClass* Z_Construct_UClass_ANexusCharacterBase();
NEXUS_API UClass* Z_Construct_UClass_ANexusCharacterBase_NoRegister();
NEXUS_API UClass* Z_Construct_UClass_UBasicAttributeSet_NoRegister();
NEXUS_API UClass* Z_Construct_UClass_UBattleAttributeSet_NoRegister();
UPackage* Z_Construct_UPackage__Script_Nexus();
// ********** End Cross Module References **********************************************************

// ********** Begin Class ANexusCharacterBase ******************************************************
void ANexusCharacterBase::StaticRegisterNativesANexusCharacterBase()
{
}
FClassRegistrationInfo Z_Registration_Info_UClass_ANexusCharacterBase;
UClass* ANexusCharacterBase::GetPrivateStaticClass()
{
	using TClass = ANexusCharacterBase;
	if (!Z_Registration_Info_UClass_ANexusCharacterBase.InnerSingleton)
	{
		GetPrivateStaticClassBody(
			StaticPackage(),
			TEXT("NexusCharacterBase"),
			Z_Registration_Info_UClass_ANexusCharacterBase.InnerSingleton,
			StaticRegisterNativesANexusCharacterBase,
			sizeof(TClass),
			alignof(TClass),
			TClass::StaticClassFlags,
			TClass::StaticClassCastFlags(),
			TClass::StaticConfigName(),
			(UClass::ClassConstructorType)InternalConstructor<TClass>,
			(UClass::ClassVTableHelperCtorCallerType)InternalVTableHelperCtorCaller<TClass>,
			UOBJECT_CPPCLASS_STATICFUNCTIONS_FORCLASS(TClass),
			&TClass::Super::StaticClass,
			&TClass::WithinClass::StaticClass
		);
	}
	return Z_Registration_Info_UClass_ANexusCharacterBase.InnerSingleton;
}
UClass* Z_Construct_UClass_ANexusCharacterBase_NoRegister()
{
	return ANexusCharacterBase::GetPrivateStaticClass();
}
struct Z_Construct_UClass_ANexusCharacterBase_Statics
{
#if WITH_METADATA
	static constexpr UECodeGen_Private::FMetaDataPairParam Class_MetaDataParams[] = {
		{ "HideCategories", "Navigation" },
		{ "IncludePath", "GameplayAbilitySystem/Character/NexusCharacterBase.h" },
		{ "ModuleRelativePath", "GameplayAbilitySystem/Character/NexusCharacterBase.h" },
	};
	static constexpr UECodeGen_Private::FMetaDataPairParam NewProp_AbilitySystemComponent_MetaData[] = {
		{ "Category", "AbilitySystem" },
#if !UE_BUILD_SHIPPING
		{ "Comment", "// Ability System Component\n" },
#endif
		{ "EditInline", "true" },
		{ "ModuleRelativePath", "GameplayAbilitySystem/Character/NexusCharacterBase.h" },
#if !UE_BUILD_SHIPPING
		{ "ToolTip", "Ability System Component" },
#endif
	};
	static constexpr UECodeGen_Private::FMetaDataPairParam NewProp_BasicAttributeSet_MetaData[] = {
		{ "Category", "AbilitySystem" },
		{ "EditInline", "true" },
		{ "ModuleRelativePath", "GameplayAbilitySystem/Character/NexusCharacterBase.h" },
	};
	static constexpr UECodeGen_Private::FMetaDataPairParam NewProp_BattleAttributeSet_MetaData[] = {
		{ "Category", "AbilitySystem" },
		{ "EditInline", "true" },
		{ "ModuleRelativePath", "GameplayAbilitySystem/Character/NexusCharacterBase.h" },
	};
	static constexpr UECodeGen_Private::FMetaDataPairParam NewProp_AscReplicationMode_MetaData[] = {
		{ "Category", "AbilitySystem" },
		{ "ModuleRelativePath", "GameplayAbilitySystem/Character/NexusCharacterBase.h" },
	};
#endif // WITH_METADATA
	static const UECodeGen_Private::FObjectPropertyParams NewProp_AbilitySystemComponent;
	static const UECodeGen_Private::FObjectPropertyParams NewProp_BasicAttributeSet;
	static const UECodeGen_Private::FObjectPropertyParams NewProp_BattleAttributeSet;
	static const UECodeGen_Private::FBytePropertyParams NewProp_AscReplicationMode_Underlying;
	static const UECodeGen_Private::FEnumPropertyParams NewProp_AscReplicationMode;
	static const UECodeGen_Private::FPropertyParamsBase* const PropPointers[];
	static UObject* (*const DependentSingletons[])();
	static const UECodeGen_Private::FImplementedInterfaceParams InterfaceParams[];
	static constexpr FCppClassTypeInfoStatic StaticCppClassTypeInfo = {
		TCppClassTypeTraits<ANexusCharacterBase>::IsAbstract,
	};
	static const UECodeGen_Private::FClassParams ClassParams;
};
const UECodeGen_Private::FObjectPropertyParams Z_Construct_UClass_ANexusCharacterBase_Statics::NewProp_AbilitySystemComponent = { "AbilitySystemComponent", nullptr, (EPropertyFlags)0x00100000000a001d, UECodeGen_Private::EPropertyGenFlags::Object, RF_Public|RF_Transient|RF_MarkAsNative, nullptr, nullptr, 1, STRUCT_OFFSET(ANexusCharacterBase, AbilitySystemComponent), Z_Construct_UClass_UAbilitySystemComponent_NoRegister, METADATA_PARAMS(UE_ARRAY_COUNT(NewProp_AbilitySystemComponent_MetaData), NewProp_AbilitySystemComponent_MetaData) };
const UECodeGen_Private::FObjectPropertyParams Z_Construct_UClass_ANexusCharacterBase_Statics::NewProp_BasicAttributeSet = { "BasicAttributeSet", nullptr, (EPropertyFlags)0x00100000000a001d, UECodeGen_Private::EPropertyGenFlags::Object, RF_Public|RF_Transient|RF_MarkAsNative, nullptr, nullptr, 1, STRUCT_OFFSET(ANexusCharacterBase, BasicAttributeSet), Z_Construct_UClass_UBasicAttributeSet_NoRegister, METADATA_PARAMS(UE_ARRAY_COUNT(NewProp_BasicAttributeSet_MetaData), NewProp_BasicAttributeSet_MetaData) };
const UECodeGen_Private::FObjectPropertyParams Z_Construct_UClass_ANexusCharacterBase_Statics::NewProp_BattleAttributeSet = { "BattleAttributeSet", nullptr, (EPropertyFlags)0x00100000000a001d, UECodeGen_Private::EPropertyGenFlags::Object, RF_Public|RF_Transient|RF_MarkAsNative, nullptr, nullptr, 1, STRUCT_OFFSET(ANexusCharacterBase, BattleAttributeSet), Z_Construct_UClass_UBattleAttributeSet_NoRegister, METADATA_PARAMS(UE_ARRAY_COUNT(NewProp_BattleAttributeSet_MetaData), NewProp_BattleAttributeSet_MetaData) };
const UECodeGen_Private::FBytePropertyParams Z_Construct_UClass_ANexusCharacterBase_Statics::NewProp_AscReplicationMode_Underlying = { "UnderlyingType", nullptr, (EPropertyFlags)0x0000000000000000, UECodeGen_Private::EPropertyGenFlags::Byte, RF_Public|RF_Transient|RF_MarkAsNative, nullptr, nullptr, 1, 0, nullptr, METADATA_PARAMS(0, nullptr) };
const UECodeGen_Private::FEnumPropertyParams Z_Construct_UClass_ANexusCharacterBase_Statics::NewProp_AscReplicationMode = { "AscReplicationMode", nullptr, (EPropertyFlags)0x0020080000000005, UECodeGen_Private::EPropertyGenFlags::Enum, RF_Public|RF_Transient|RF_MarkAsNative, nullptr, nullptr, 1, STRUCT_OFFSET(ANexusCharacterBase, AscReplicationMode), Z_Construct_UEnum_GameplayAbilities_EGameplayEffectReplicationMode, METADATA_PARAMS(UE_ARRAY_COUNT(NewProp_AscReplicationMode_MetaData), NewProp_AscReplicationMode_MetaData) }; // 3979288675
const UECodeGen_Private::FPropertyParamsBase* const Z_Construct_UClass_ANexusCharacterBase_Statics::PropPointers[] = {
	(const UECodeGen_Private::FPropertyParamsBase*)&Z_Construct_UClass_ANexusCharacterBase_Statics::NewProp_AbilitySystemComponent,
	(const UECodeGen_Private::FPropertyParamsBase*)&Z_Construct_UClass_ANexusCharacterBase_Statics::NewProp_BasicAttributeSet,
	(const UECodeGen_Private::FPropertyParamsBase*)&Z_Construct_UClass_ANexusCharacterBase_Statics::NewProp_BattleAttributeSet,
	(const UECodeGen_Private::FPropertyParamsBase*)&Z_Construct_UClass_ANexusCharacterBase_Statics::NewProp_AscReplicationMode_Underlying,
	(const UECodeGen_Private::FPropertyParamsBase*)&Z_Construct_UClass_ANexusCharacterBase_Statics::NewProp_AscReplicationMode,
};
static_assert(UE_ARRAY_COUNT(Z_Construct_UClass_ANexusCharacterBase_Statics::PropPointers) < 2048);
UObject* (*const Z_Construct_UClass_ANexusCharacterBase_Statics::DependentSingletons[])() = {
	(UObject* (*)())Z_Construct_UClass_ACharacter,
	(UObject* (*)())Z_Construct_UPackage__Script_Nexus,
};
static_assert(UE_ARRAY_COUNT(Z_Construct_UClass_ANexusCharacterBase_Statics::DependentSingletons) < 16);
const UECodeGen_Private::FImplementedInterfaceParams Z_Construct_UClass_ANexusCharacterBase_Statics::InterfaceParams[] = {
	{ Z_Construct_UClass_UAbilitySystemInterface_NoRegister, (int32)VTABLE_OFFSET(ANexusCharacterBase, IAbilitySystemInterface), false },  // 1199015870
};
const UECodeGen_Private::FClassParams Z_Construct_UClass_ANexusCharacterBase_Statics::ClassParams = {
	&ANexusCharacterBase::StaticClass,
	"Game",
	&StaticCppClassTypeInfo,
	DependentSingletons,
	nullptr,
	Z_Construct_UClass_ANexusCharacterBase_Statics::PropPointers,
	InterfaceParams,
	UE_ARRAY_COUNT(DependentSingletons),
	0,
	UE_ARRAY_COUNT(Z_Construct_UClass_ANexusCharacterBase_Statics::PropPointers),
	UE_ARRAY_COUNT(InterfaceParams),
	0x009001A4u,
	METADATA_PARAMS(UE_ARRAY_COUNT(Z_Construct_UClass_ANexusCharacterBase_Statics::Class_MetaDataParams), Z_Construct_UClass_ANexusCharacterBase_Statics::Class_MetaDataParams)
};
UClass* Z_Construct_UClass_ANexusCharacterBase()
{
	if (!Z_Registration_Info_UClass_ANexusCharacterBase.OuterSingleton)
	{
		UECodeGen_Private::ConstructUClass(Z_Registration_Info_UClass_ANexusCharacterBase.OuterSingleton, Z_Construct_UClass_ANexusCharacterBase_Statics::ClassParams);
	}
	return Z_Registration_Info_UClass_ANexusCharacterBase.OuterSingleton;
}
DEFINE_VTABLE_PTR_HELPER_CTOR(ANexusCharacterBase);
ANexusCharacterBase::~ANexusCharacterBase() {}
// ********** End Class ANexusCharacterBase ********************************************************

// ********** Begin Registration *******************************************************************
struct Z_CompiledInDeferFile_FID_UnrealProjects_team_dynamics_project_1_Nexus_Source_Nexus_GameplayAbilitySystem_Character_NexusCharacterBase_h__Script_Nexus_Statics
{
	static constexpr FClassRegisterCompiledInInfo ClassInfo[] = {
		{ Z_Construct_UClass_ANexusCharacterBase, ANexusCharacterBase::StaticClass, TEXT("ANexusCharacterBase"), &Z_Registration_Info_UClass_ANexusCharacterBase, CONSTRUCT_RELOAD_VERSION_INFO(FClassReloadVersionInfo, sizeof(ANexusCharacterBase), 2706460604U) },
	};
};
static FRegisterCompiledInInfo Z_CompiledInDeferFile_FID_UnrealProjects_team_dynamics_project_1_Nexus_Source_Nexus_GameplayAbilitySystem_Character_NexusCharacterBase_h__Script_Nexus_74504304(TEXT("/Script/Nexus"),
	Z_CompiledInDeferFile_FID_UnrealProjects_team_dynamics_project_1_Nexus_Source_Nexus_GameplayAbilitySystem_Character_NexusCharacterBase_h__Script_Nexus_Statics::ClassInfo, UE_ARRAY_COUNT(Z_CompiledInDeferFile_FID_UnrealProjects_team_dynamics_project_1_Nexus_Source_Nexus_GameplayAbilitySystem_Character_NexusCharacterBase_h__Script_Nexus_Statics::ClassInfo),
	nullptr, 0,
	nullptr, 0);
// ********** End Registration *********************************************************************

PRAGMA_ENABLE_DEPRECATION_WARNINGS

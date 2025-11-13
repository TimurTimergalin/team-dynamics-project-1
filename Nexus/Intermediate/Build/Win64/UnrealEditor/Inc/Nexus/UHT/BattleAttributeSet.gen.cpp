// Copyright Epic Games, Inc. All Rights Reserved.
/*===========================================================================
	Generated code exported from UnrealHeaderTool.
	DO NOT modify this manually! Edit the corresponding .h files instead!
===========================================================================*/

#include "UObject/GeneratedCppIncludes.h"
#include "Nexus/GameplayAbilitySystem/AttributeSets/BattleAttributeSet.h"
#include "AttributeSet.h"

PRAGMA_DISABLE_DEPRECATION_WARNINGS

void EmptyLinkFunctionForGeneratedCodeBattleAttributeSet() {}

// ********** Begin Cross Module References ********************************************************
GAMEPLAYABILITIES_API UClass* Z_Construct_UClass_UAttributeSet();
GAMEPLAYABILITIES_API UScriptStruct* Z_Construct_UScriptStruct_FGameplayAttributeData();
NEXUS_API UClass* Z_Construct_UClass_UBattleAttributeSet();
NEXUS_API UClass* Z_Construct_UClass_UBattleAttributeSet_NoRegister();
UPackage* Z_Construct_UPackage__Script_Nexus();
// ********** End Cross Module References **********************************************************

// ********** Begin Class UBattleAttributeSet Function OnRep_Damage ********************************
struct Z_Construct_UFunction_UBattleAttributeSet_OnRep_Damage_Statics
{
	struct BattleAttributeSet_eventOnRep_Damage_Parms
	{
		FGameplayAttributeData OldValue;
	};
#if WITH_METADATA
	static constexpr UECodeGen_Private::FMetaDataPairParam Function_MetaDataParams[] = {
		{ "ModuleRelativePath", "GameplayAbilitySystem/AttributeSets/BattleAttributeSet.h" },
	};
	static constexpr UECodeGen_Private::FMetaDataPairParam NewProp_OldValue_MetaData[] = {
		{ "NativeConst", "" },
	};
#endif // WITH_METADATA
	static const UECodeGen_Private::FStructPropertyParams NewProp_OldValue;
	static const UECodeGen_Private::FPropertyParamsBase* const PropPointers[];
	static const UECodeGen_Private::FFunctionParams FuncParams;
};
const UECodeGen_Private::FStructPropertyParams Z_Construct_UFunction_UBattleAttributeSet_OnRep_Damage_Statics::NewProp_OldValue = { "OldValue", nullptr, (EPropertyFlags)0x0010000008000182, UECodeGen_Private::EPropertyGenFlags::Struct, RF_Public|RF_Transient|RF_MarkAsNative, nullptr, nullptr, 1, STRUCT_OFFSET(BattleAttributeSet_eventOnRep_Damage_Parms, OldValue), Z_Construct_UScriptStruct_FGameplayAttributeData, METADATA_PARAMS(UE_ARRAY_COUNT(NewProp_OldValue_MetaData), NewProp_OldValue_MetaData) }; // 1532612004
const UECodeGen_Private::FPropertyParamsBase* const Z_Construct_UFunction_UBattleAttributeSet_OnRep_Damage_Statics::PropPointers[] = {
	(const UECodeGen_Private::FPropertyParamsBase*)&Z_Construct_UFunction_UBattleAttributeSet_OnRep_Damage_Statics::NewProp_OldValue,
};
static_assert(UE_ARRAY_COUNT(Z_Construct_UFunction_UBattleAttributeSet_OnRep_Damage_Statics::PropPointers) < 2048);
const UECodeGen_Private::FFunctionParams Z_Construct_UFunction_UBattleAttributeSet_OnRep_Damage_Statics::FuncParams = { { (UObject*(*)())Z_Construct_UClass_UBattleAttributeSet, nullptr, "OnRep_Damage", Z_Construct_UFunction_UBattleAttributeSet_OnRep_Damage_Statics::PropPointers, UE_ARRAY_COUNT(Z_Construct_UFunction_UBattleAttributeSet_OnRep_Damage_Statics::PropPointers), sizeof(Z_Construct_UFunction_UBattleAttributeSet_OnRep_Damage_Statics::BattleAttributeSet_eventOnRep_Damage_Parms), RF_Public|RF_Transient|RF_MarkAsNative, (EFunctionFlags)0x00420401, 0, 0, METADATA_PARAMS(UE_ARRAY_COUNT(Z_Construct_UFunction_UBattleAttributeSet_OnRep_Damage_Statics::Function_MetaDataParams), Z_Construct_UFunction_UBattleAttributeSet_OnRep_Damage_Statics::Function_MetaDataParams)},  };
static_assert(sizeof(Z_Construct_UFunction_UBattleAttributeSet_OnRep_Damage_Statics::BattleAttributeSet_eventOnRep_Damage_Parms) < MAX_uint16);
UFunction* Z_Construct_UFunction_UBattleAttributeSet_OnRep_Damage()
{
	static UFunction* ReturnFunction = nullptr;
	if (!ReturnFunction)
	{
		UECodeGen_Private::ConstructUFunction(&ReturnFunction, Z_Construct_UFunction_UBattleAttributeSet_OnRep_Damage_Statics::FuncParams);
	}
	return ReturnFunction;
}
DEFINE_FUNCTION(UBattleAttributeSet::execOnRep_Damage)
{
	P_GET_STRUCT_REF(FGameplayAttributeData,Z_Param_Out_OldValue);
	P_FINISH;
	P_NATIVE_BEGIN;
	P_THIS->OnRep_Damage(Z_Param_Out_OldValue);
	P_NATIVE_END;
}
// ********** End Class UBattleAttributeSet Function OnRep_Damage **********************************

// ********** Begin Class UBattleAttributeSet Function OnRep_HealingAmount *************************
struct Z_Construct_UFunction_UBattleAttributeSet_OnRep_HealingAmount_Statics
{
	struct BattleAttributeSet_eventOnRep_HealingAmount_Parms
	{
		FGameplayAttributeData OldValue;
	};
#if WITH_METADATA
	static constexpr UECodeGen_Private::FMetaDataPairParam Function_MetaDataParams[] = {
		{ "ModuleRelativePath", "GameplayAbilitySystem/AttributeSets/BattleAttributeSet.h" },
	};
	static constexpr UECodeGen_Private::FMetaDataPairParam NewProp_OldValue_MetaData[] = {
		{ "NativeConst", "" },
	};
#endif // WITH_METADATA
	static const UECodeGen_Private::FStructPropertyParams NewProp_OldValue;
	static const UECodeGen_Private::FPropertyParamsBase* const PropPointers[];
	static const UECodeGen_Private::FFunctionParams FuncParams;
};
const UECodeGen_Private::FStructPropertyParams Z_Construct_UFunction_UBattleAttributeSet_OnRep_HealingAmount_Statics::NewProp_OldValue = { "OldValue", nullptr, (EPropertyFlags)0x0010000008000182, UECodeGen_Private::EPropertyGenFlags::Struct, RF_Public|RF_Transient|RF_MarkAsNative, nullptr, nullptr, 1, STRUCT_OFFSET(BattleAttributeSet_eventOnRep_HealingAmount_Parms, OldValue), Z_Construct_UScriptStruct_FGameplayAttributeData, METADATA_PARAMS(UE_ARRAY_COUNT(NewProp_OldValue_MetaData), NewProp_OldValue_MetaData) }; // 1532612004
const UECodeGen_Private::FPropertyParamsBase* const Z_Construct_UFunction_UBattleAttributeSet_OnRep_HealingAmount_Statics::PropPointers[] = {
	(const UECodeGen_Private::FPropertyParamsBase*)&Z_Construct_UFunction_UBattleAttributeSet_OnRep_HealingAmount_Statics::NewProp_OldValue,
};
static_assert(UE_ARRAY_COUNT(Z_Construct_UFunction_UBattleAttributeSet_OnRep_HealingAmount_Statics::PropPointers) < 2048);
const UECodeGen_Private::FFunctionParams Z_Construct_UFunction_UBattleAttributeSet_OnRep_HealingAmount_Statics::FuncParams = { { (UObject*(*)())Z_Construct_UClass_UBattleAttributeSet, nullptr, "OnRep_HealingAmount", Z_Construct_UFunction_UBattleAttributeSet_OnRep_HealingAmount_Statics::PropPointers, UE_ARRAY_COUNT(Z_Construct_UFunction_UBattleAttributeSet_OnRep_HealingAmount_Statics::PropPointers), sizeof(Z_Construct_UFunction_UBattleAttributeSet_OnRep_HealingAmount_Statics::BattleAttributeSet_eventOnRep_HealingAmount_Parms), RF_Public|RF_Transient|RF_MarkAsNative, (EFunctionFlags)0x00420401, 0, 0, METADATA_PARAMS(UE_ARRAY_COUNT(Z_Construct_UFunction_UBattleAttributeSet_OnRep_HealingAmount_Statics::Function_MetaDataParams), Z_Construct_UFunction_UBattleAttributeSet_OnRep_HealingAmount_Statics::Function_MetaDataParams)},  };
static_assert(sizeof(Z_Construct_UFunction_UBattleAttributeSet_OnRep_HealingAmount_Statics::BattleAttributeSet_eventOnRep_HealingAmount_Parms) < MAX_uint16);
UFunction* Z_Construct_UFunction_UBattleAttributeSet_OnRep_HealingAmount()
{
	static UFunction* ReturnFunction = nullptr;
	if (!ReturnFunction)
	{
		UECodeGen_Private::ConstructUFunction(&ReturnFunction, Z_Construct_UFunction_UBattleAttributeSet_OnRep_HealingAmount_Statics::FuncParams);
	}
	return ReturnFunction;
}
DEFINE_FUNCTION(UBattleAttributeSet::execOnRep_HealingAmount)
{
	P_GET_STRUCT_REF(FGameplayAttributeData,Z_Param_Out_OldValue);
	P_FINISH;
	P_NATIVE_BEGIN;
	P_THIS->OnRep_HealingAmount(Z_Param_Out_OldValue);
	P_NATIVE_END;
}
// ********** End Class UBattleAttributeSet Function OnRep_HealingAmount ***************************

// ********** Begin Class UBattleAttributeSet Function OnRep_HealingTime ***************************
struct Z_Construct_UFunction_UBattleAttributeSet_OnRep_HealingTime_Statics
{
	struct BattleAttributeSet_eventOnRep_HealingTime_Parms
	{
		FGameplayAttributeData OldValue;
	};
#if WITH_METADATA
	static constexpr UECodeGen_Private::FMetaDataPairParam Function_MetaDataParams[] = {
		{ "ModuleRelativePath", "GameplayAbilitySystem/AttributeSets/BattleAttributeSet.h" },
	};
	static constexpr UECodeGen_Private::FMetaDataPairParam NewProp_OldValue_MetaData[] = {
		{ "NativeConst", "" },
	};
#endif // WITH_METADATA
	static const UECodeGen_Private::FStructPropertyParams NewProp_OldValue;
	static const UECodeGen_Private::FPropertyParamsBase* const PropPointers[];
	static const UECodeGen_Private::FFunctionParams FuncParams;
};
const UECodeGen_Private::FStructPropertyParams Z_Construct_UFunction_UBattleAttributeSet_OnRep_HealingTime_Statics::NewProp_OldValue = { "OldValue", nullptr, (EPropertyFlags)0x0010000008000182, UECodeGen_Private::EPropertyGenFlags::Struct, RF_Public|RF_Transient|RF_MarkAsNative, nullptr, nullptr, 1, STRUCT_OFFSET(BattleAttributeSet_eventOnRep_HealingTime_Parms, OldValue), Z_Construct_UScriptStruct_FGameplayAttributeData, METADATA_PARAMS(UE_ARRAY_COUNT(NewProp_OldValue_MetaData), NewProp_OldValue_MetaData) }; // 1532612004
const UECodeGen_Private::FPropertyParamsBase* const Z_Construct_UFunction_UBattleAttributeSet_OnRep_HealingTime_Statics::PropPointers[] = {
	(const UECodeGen_Private::FPropertyParamsBase*)&Z_Construct_UFunction_UBattleAttributeSet_OnRep_HealingTime_Statics::NewProp_OldValue,
};
static_assert(UE_ARRAY_COUNT(Z_Construct_UFunction_UBattleAttributeSet_OnRep_HealingTime_Statics::PropPointers) < 2048);
const UECodeGen_Private::FFunctionParams Z_Construct_UFunction_UBattleAttributeSet_OnRep_HealingTime_Statics::FuncParams = { { (UObject*(*)())Z_Construct_UClass_UBattleAttributeSet, nullptr, "OnRep_HealingTime", Z_Construct_UFunction_UBattleAttributeSet_OnRep_HealingTime_Statics::PropPointers, UE_ARRAY_COUNT(Z_Construct_UFunction_UBattleAttributeSet_OnRep_HealingTime_Statics::PropPointers), sizeof(Z_Construct_UFunction_UBattleAttributeSet_OnRep_HealingTime_Statics::BattleAttributeSet_eventOnRep_HealingTime_Parms), RF_Public|RF_Transient|RF_MarkAsNative, (EFunctionFlags)0x00420401, 0, 0, METADATA_PARAMS(UE_ARRAY_COUNT(Z_Construct_UFunction_UBattleAttributeSet_OnRep_HealingTime_Statics::Function_MetaDataParams), Z_Construct_UFunction_UBattleAttributeSet_OnRep_HealingTime_Statics::Function_MetaDataParams)},  };
static_assert(sizeof(Z_Construct_UFunction_UBattleAttributeSet_OnRep_HealingTime_Statics::BattleAttributeSet_eventOnRep_HealingTime_Parms) < MAX_uint16);
UFunction* Z_Construct_UFunction_UBattleAttributeSet_OnRep_HealingTime()
{
	static UFunction* ReturnFunction = nullptr;
	if (!ReturnFunction)
	{
		UECodeGen_Private::ConstructUFunction(&ReturnFunction, Z_Construct_UFunction_UBattleAttributeSet_OnRep_HealingTime_Statics::FuncParams);
	}
	return ReturnFunction;
}
DEFINE_FUNCTION(UBattleAttributeSet::execOnRep_HealingTime)
{
	P_GET_STRUCT_REF(FGameplayAttributeData,Z_Param_Out_OldValue);
	P_FINISH;
	P_NATIVE_BEGIN;
	P_THIS->OnRep_HealingTime(Z_Param_Out_OldValue);
	P_NATIVE_END;
}
// ********** End Class UBattleAttributeSet Function OnRep_HealingTime *****************************

// ********** Begin Class UBattleAttributeSet ******************************************************
void UBattleAttributeSet::StaticRegisterNativesUBattleAttributeSet()
{
	UClass* Class = UBattleAttributeSet::StaticClass();
	static const FNameNativePtrPair Funcs[] = {
		{ "OnRep_Damage", &UBattleAttributeSet::execOnRep_Damage },
		{ "OnRep_HealingAmount", &UBattleAttributeSet::execOnRep_HealingAmount },
		{ "OnRep_HealingTime", &UBattleAttributeSet::execOnRep_HealingTime },
	};
	FNativeFunctionRegistrar::RegisterFunctions(Class, Funcs, UE_ARRAY_COUNT(Funcs));
}
FClassRegistrationInfo Z_Registration_Info_UClass_UBattleAttributeSet;
UClass* UBattleAttributeSet::GetPrivateStaticClass()
{
	using TClass = UBattleAttributeSet;
	if (!Z_Registration_Info_UClass_UBattleAttributeSet.InnerSingleton)
	{
		GetPrivateStaticClassBody(
			StaticPackage(),
			TEXT("BattleAttributeSet"),
			Z_Registration_Info_UClass_UBattleAttributeSet.InnerSingleton,
			StaticRegisterNativesUBattleAttributeSet,
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
	return Z_Registration_Info_UClass_UBattleAttributeSet.InnerSingleton;
}
UClass* Z_Construct_UClass_UBattleAttributeSet_NoRegister()
{
	return UBattleAttributeSet::GetPrivateStaticClass();
}
struct Z_Construct_UClass_UBattleAttributeSet_Statics
{
#if WITH_METADATA
	static constexpr UECodeGen_Private::FMetaDataPairParam Class_MetaDataParams[] = {
#if !UE_BUILD_SHIPPING
		{ "Comment", "/**\n * \n */" },
#endif
		{ "IncludePath", "GameplayAbilitySystem/AttributeSets/BattleAttributeSet.h" },
		{ "ModuleRelativePath", "GameplayAbilitySystem/AttributeSets/BattleAttributeSet.h" },
	};
	static constexpr UECodeGen_Private::FMetaDataPairParam NewProp_Damage_MetaData[] = {
		{ "Category", "Attributes" },
#if !UE_BUILD_SHIPPING
		{ "Comment", "// Battle Attributes\n" },
#endif
		{ "ModuleRelativePath", "GameplayAbilitySystem/AttributeSets/BattleAttributeSet.h" },
#if !UE_BUILD_SHIPPING
		{ "ToolTip", "Battle Attributes" },
#endif
	};
	static constexpr UECodeGen_Private::FMetaDataPairParam NewProp_HealingAmount_MetaData[] = {
		{ "Category", "Attributes" },
		{ "ModuleRelativePath", "GameplayAbilitySystem/AttributeSets/BattleAttributeSet.h" },
	};
	static constexpr UECodeGen_Private::FMetaDataPairParam NewProp_HealingTime_MetaData[] = {
		{ "Category", "Attributes" },
		{ "ModuleRelativePath", "GameplayAbilitySystem/AttributeSets/BattleAttributeSet.h" },
	};
#endif // WITH_METADATA
	static const UECodeGen_Private::FStructPropertyParams NewProp_Damage;
	static const UECodeGen_Private::FStructPropertyParams NewProp_HealingAmount;
	static const UECodeGen_Private::FStructPropertyParams NewProp_HealingTime;
	static const UECodeGen_Private::FPropertyParamsBase* const PropPointers[];
	static UObject* (*const DependentSingletons[])();
	static constexpr FClassFunctionLinkInfo FuncInfo[] = {
		{ &Z_Construct_UFunction_UBattleAttributeSet_OnRep_Damage, "OnRep_Damage" }, // 2190556535
		{ &Z_Construct_UFunction_UBattleAttributeSet_OnRep_HealingAmount, "OnRep_HealingAmount" }, // 3838573014
		{ &Z_Construct_UFunction_UBattleAttributeSet_OnRep_HealingTime, "OnRep_HealingTime" }, // 2032394735
	};
	static_assert(UE_ARRAY_COUNT(FuncInfo) < 2048);
	static constexpr FCppClassTypeInfoStatic StaticCppClassTypeInfo = {
		TCppClassTypeTraits<UBattleAttributeSet>::IsAbstract,
	};
	static const UECodeGen_Private::FClassParams ClassParams;
};
const UECodeGen_Private::FStructPropertyParams Z_Construct_UClass_UBattleAttributeSet_Statics::NewProp_Damage = { "Damage", "OnRep_Damage", (EPropertyFlags)0x0010000100000034, UECodeGen_Private::EPropertyGenFlags::Struct, RF_Public|RF_Transient|RF_MarkAsNative, nullptr, nullptr, 1, STRUCT_OFFSET(UBattleAttributeSet, Damage), Z_Construct_UScriptStruct_FGameplayAttributeData, METADATA_PARAMS(UE_ARRAY_COUNT(NewProp_Damage_MetaData), NewProp_Damage_MetaData) }; // 1532612004
const UECodeGen_Private::FStructPropertyParams Z_Construct_UClass_UBattleAttributeSet_Statics::NewProp_HealingAmount = { "HealingAmount", "OnRep_HealingAmount", (EPropertyFlags)0x0010000100000034, UECodeGen_Private::EPropertyGenFlags::Struct, RF_Public|RF_Transient|RF_MarkAsNative, nullptr, nullptr, 1, STRUCT_OFFSET(UBattleAttributeSet, HealingAmount), Z_Construct_UScriptStruct_FGameplayAttributeData, METADATA_PARAMS(UE_ARRAY_COUNT(NewProp_HealingAmount_MetaData), NewProp_HealingAmount_MetaData) }; // 1532612004
const UECodeGen_Private::FStructPropertyParams Z_Construct_UClass_UBattleAttributeSet_Statics::NewProp_HealingTime = { "HealingTime", "OnRep_HealingTime", (EPropertyFlags)0x0010000100000034, UECodeGen_Private::EPropertyGenFlags::Struct, RF_Public|RF_Transient|RF_MarkAsNative, nullptr, nullptr, 1, STRUCT_OFFSET(UBattleAttributeSet, HealingTime), Z_Construct_UScriptStruct_FGameplayAttributeData, METADATA_PARAMS(UE_ARRAY_COUNT(NewProp_HealingTime_MetaData), NewProp_HealingTime_MetaData) }; // 1532612004
const UECodeGen_Private::FPropertyParamsBase* const Z_Construct_UClass_UBattleAttributeSet_Statics::PropPointers[] = {
	(const UECodeGen_Private::FPropertyParamsBase*)&Z_Construct_UClass_UBattleAttributeSet_Statics::NewProp_Damage,
	(const UECodeGen_Private::FPropertyParamsBase*)&Z_Construct_UClass_UBattleAttributeSet_Statics::NewProp_HealingAmount,
	(const UECodeGen_Private::FPropertyParamsBase*)&Z_Construct_UClass_UBattleAttributeSet_Statics::NewProp_HealingTime,
};
static_assert(UE_ARRAY_COUNT(Z_Construct_UClass_UBattleAttributeSet_Statics::PropPointers) < 2048);
UObject* (*const Z_Construct_UClass_UBattleAttributeSet_Statics::DependentSingletons[])() = {
	(UObject* (*)())Z_Construct_UClass_UAttributeSet,
	(UObject* (*)())Z_Construct_UPackage__Script_Nexus,
};
static_assert(UE_ARRAY_COUNT(Z_Construct_UClass_UBattleAttributeSet_Statics::DependentSingletons) < 16);
const UECodeGen_Private::FClassParams Z_Construct_UClass_UBattleAttributeSet_Statics::ClassParams = {
	&UBattleAttributeSet::StaticClass,
	nullptr,
	&StaticCppClassTypeInfo,
	DependentSingletons,
	FuncInfo,
	Z_Construct_UClass_UBattleAttributeSet_Statics::PropPointers,
	nullptr,
	UE_ARRAY_COUNT(DependentSingletons),
	UE_ARRAY_COUNT(FuncInfo),
	UE_ARRAY_COUNT(Z_Construct_UClass_UBattleAttributeSet_Statics::PropPointers),
	0,
	0x003000A0u,
	METADATA_PARAMS(UE_ARRAY_COUNT(Z_Construct_UClass_UBattleAttributeSet_Statics::Class_MetaDataParams), Z_Construct_UClass_UBattleAttributeSet_Statics::Class_MetaDataParams)
};
UClass* Z_Construct_UClass_UBattleAttributeSet()
{
	if (!Z_Registration_Info_UClass_UBattleAttributeSet.OuterSingleton)
	{
		UECodeGen_Private::ConstructUClass(Z_Registration_Info_UClass_UBattleAttributeSet.OuterSingleton, Z_Construct_UClass_UBattleAttributeSet_Statics::ClassParams);
	}
	return Z_Registration_Info_UClass_UBattleAttributeSet.OuterSingleton;
}
#if VALIDATE_CLASS_REPS
void UBattleAttributeSet::ValidateGeneratedRepEnums(const TArray<struct FRepRecord>& ClassReps) const
{
	static FName Name_Damage(TEXT("Damage"));
	static FName Name_HealingAmount(TEXT("HealingAmount"));
	static FName Name_HealingTime(TEXT("HealingTime"));
	const bool bIsValid = true
		&& Name_Damage == ClassReps[(int32)ENetFields_Private::Damage].Property->GetFName()
		&& Name_HealingAmount == ClassReps[(int32)ENetFields_Private::HealingAmount].Property->GetFName()
		&& Name_HealingTime == ClassReps[(int32)ENetFields_Private::HealingTime].Property->GetFName();
	checkf(bIsValid, TEXT("UHT Generated Rep Indices do not match runtime populated Rep Indices for properties in UBattleAttributeSet"));
}
#endif
DEFINE_VTABLE_PTR_HELPER_CTOR(UBattleAttributeSet);
UBattleAttributeSet::~UBattleAttributeSet() {}
// ********** End Class UBattleAttributeSet ********************************************************

// ********** Begin Registration *******************************************************************
struct Z_CompiledInDeferFile_FID_UnrealProjects_team_dynamics_project_1_Nexus_Source_Nexus_GameplayAbilitySystem_AttributeSets_BattleAttributeSet_h__Script_Nexus_Statics
{
	static constexpr FClassRegisterCompiledInInfo ClassInfo[] = {
		{ Z_Construct_UClass_UBattleAttributeSet, UBattleAttributeSet::StaticClass, TEXT("UBattleAttributeSet"), &Z_Registration_Info_UClass_UBattleAttributeSet, CONSTRUCT_RELOAD_VERSION_INFO(FClassReloadVersionInfo, sizeof(UBattleAttributeSet), 1267711093U) },
	};
};
static FRegisterCompiledInInfo Z_CompiledInDeferFile_FID_UnrealProjects_team_dynamics_project_1_Nexus_Source_Nexus_GameplayAbilitySystem_AttributeSets_BattleAttributeSet_h__Script_Nexus_3754278241(TEXT("/Script/Nexus"),
	Z_CompiledInDeferFile_FID_UnrealProjects_team_dynamics_project_1_Nexus_Source_Nexus_GameplayAbilitySystem_AttributeSets_BattleAttributeSet_h__Script_Nexus_Statics::ClassInfo, UE_ARRAY_COUNT(Z_CompiledInDeferFile_FID_UnrealProjects_team_dynamics_project_1_Nexus_Source_Nexus_GameplayAbilitySystem_AttributeSets_BattleAttributeSet_h__Script_Nexus_Statics::ClassInfo),
	nullptr, 0,
	nullptr, 0);
// ********** End Registration *********************************************************************

PRAGMA_ENABLE_DEPRECATION_WARNINGS

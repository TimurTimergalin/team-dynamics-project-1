// Copyright Epic Games, Inc. All Rights Reserved.
/*===========================================================================
	Generated code exported from UnrealHeaderTool.
	DO NOT modify this manually! Edit the corresponding .h files instead!
===========================================================================*/

#include "UObject/GeneratedCppIncludes.h"
#include "Nexus/GameplayAbilitySystem/AttributeSets/BasicAttributeSet.h"
#include "AttributeSet.h"

PRAGMA_DISABLE_DEPRECATION_WARNINGS

void EmptyLinkFunctionForGeneratedCodeBasicAttributeSet() {}

// ********** Begin Cross Module References ********************************************************
GAMEPLAYABILITIES_API UClass* Z_Construct_UClass_UAttributeSet();
GAMEPLAYABILITIES_API UScriptStruct* Z_Construct_UScriptStruct_FGameplayAttributeData();
NEXUS_API UClass* Z_Construct_UClass_UBasicAttributeSet();
NEXUS_API UClass* Z_Construct_UClass_UBasicAttributeSet_NoRegister();
UPackage* Z_Construct_UPackage__Script_Nexus();
// ********** End Cross Module References **********************************************************

// ********** Begin Class UBasicAttributeSet Function OnRep_Health *********************************
struct Z_Construct_UFunction_UBasicAttributeSet_OnRep_Health_Statics
{
	struct BasicAttributeSet_eventOnRep_Health_Parms
	{
		FGameplayAttributeData OldValue;
	};
#if WITH_METADATA
	static constexpr UECodeGen_Private::FMetaDataPairParam Function_MetaDataParams[] = {
		{ "ModuleRelativePath", "GameplayAbilitySystem/AttributeSets/BasicAttributeSet.h" },
	};
	static constexpr UECodeGen_Private::FMetaDataPairParam NewProp_OldValue_MetaData[] = {
		{ "NativeConst", "" },
	};
#endif // WITH_METADATA
	static const UECodeGen_Private::FStructPropertyParams NewProp_OldValue;
	static const UECodeGen_Private::FPropertyParamsBase* const PropPointers[];
	static const UECodeGen_Private::FFunctionParams FuncParams;
};
const UECodeGen_Private::FStructPropertyParams Z_Construct_UFunction_UBasicAttributeSet_OnRep_Health_Statics::NewProp_OldValue = { "OldValue", nullptr, (EPropertyFlags)0x0010000008000182, UECodeGen_Private::EPropertyGenFlags::Struct, RF_Public|RF_Transient|RF_MarkAsNative, nullptr, nullptr, 1, STRUCT_OFFSET(BasicAttributeSet_eventOnRep_Health_Parms, OldValue), Z_Construct_UScriptStruct_FGameplayAttributeData, METADATA_PARAMS(UE_ARRAY_COUNT(NewProp_OldValue_MetaData), NewProp_OldValue_MetaData) }; // 1532612004
const UECodeGen_Private::FPropertyParamsBase* const Z_Construct_UFunction_UBasicAttributeSet_OnRep_Health_Statics::PropPointers[] = {
	(const UECodeGen_Private::FPropertyParamsBase*)&Z_Construct_UFunction_UBasicAttributeSet_OnRep_Health_Statics::NewProp_OldValue,
};
static_assert(UE_ARRAY_COUNT(Z_Construct_UFunction_UBasicAttributeSet_OnRep_Health_Statics::PropPointers) < 2048);
const UECodeGen_Private::FFunctionParams Z_Construct_UFunction_UBasicAttributeSet_OnRep_Health_Statics::FuncParams = { { (UObject*(*)())Z_Construct_UClass_UBasicAttributeSet, nullptr, "OnRep_Health", Z_Construct_UFunction_UBasicAttributeSet_OnRep_Health_Statics::PropPointers, UE_ARRAY_COUNT(Z_Construct_UFunction_UBasicAttributeSet_OnRep_Health_Statics::PropPointers), sizeof(Z_Construct_UFunction_UBasicAttributeSet_OnRep_Health_Statics::BasicAttributeSet_eventOnRep_Health_Parms), RF_Public|RF_Transient|RF_MarkAsNative, (EFunctionFlags)0x00420401, 0, 0, METADATA_PARAMS(UE_ARRAY_COUNT(Z_Construct_UFunction_UBasicAttributeSet_OnRep_Health_Statics::Function_MetaDataParams), Z_Construct_UFunction_UBasicAttributeSet_OnRep_Health_Statics::Function_MetaDataParams)},  };
static_assert(sizeof(Z_Construct_UFunction_UBasicAttributeSet_OnRep_Health_Statics::BasicAttributeSet_eventOnRep_Health_Parms) < MAX_uint16);
UFunction* Z_Construct_UFunction_UBasicAttributeSet_OnRep_Health()
{
	static UFunction* ReturnFunction = nullptr;
	if (!ReturnFunction)
	{
		UECodeGen_Private::ConstructUFunction(&ReturnFunction, Z_Construct_UFunction_UBasicAttributeSet_OnRep_Health_Statics::FuncParams);
	}
	return ReturnFunction;
}
DEFINE_FUNCTION(UBasicAttributeSet::execOnRep_Health)
{
	P_GET_STRUCT_REF(FGameplayAttributeData,Z_Param_Out_OldValue);
	P_FINISH;
	P_NATIVE_BEGIN;
	P_THIS->OnRep_Health(Z_Param_Out_OldValue);
	P_NATIVE_END;
}
// ********** End Class UBasicAttributeSet Function OnRep_Health ***********************************

// ********** Begin Class UBasicAttributeSet Function OnRep_MaxHealth ******************************
struct Z_Construct_UFunction_UBasicAttributeSet_OnRep_MaxHealth_Statics
{
	struct BasicAttributeSet_eventOnRep_MaxHealth_Parms
	{
		FGameplayAttributeData OldValue;
	};
#if WITH_METADATA
	static constexpr UECodeGen_Private::FMetaDataPairParam Function_MetaDataParams[] = {
		{ "ModuleRelativePath", "GameplayAbilitySystem/AttributeSets/BasicAttributeSet.h" },
	};
	static constexpr UECodeGen_Private::FMetaDataPairParam NewProp_OldValue_MetaData[] = {
		{ "NativeConst", "" },
	};
#endif // WITH_METADATA
	static const UECodeGen_Private::FStructPropertyParams NewProp_OldValue;
	static const UECodeGen_Private::FPropertyParamsBase* const PropPointers[];
	static const UECodeGen_Private::FFunctionParams FuncParams;
};
const UECodeGen_Private::FStructPropertyParams Z_Construct_UFunction_UBasicAttributeSet_OnRep_MaxHealth_Statics::NewProp_OldValue = { "OldValue", nullptr, (EPropertyFlags)0x0010000008000182, UECodeGen_Private::EPropertyGenFlags::Struct, RF_Public|RF_Transient|RF_MarkAsNative, nullptr, nullptr, 1, STRUCT_OFFSET(BasicAttributeSet_eventOnRep_MaxHealth_Parms, OldValue), Z_Construct_UScriptStruct_FGameplayAttributeData, METADATA_PARAMS(UE_ARRAY_COUNT(NewProp_OldValue_MetaData), NewProp_OldValue_MetaData) }; // 1532612004
const UECodeGen_Private::FPropertyParamsBase* const Z_Construct_UFunction_UBasicAttributeSet_OnRep_MaxHealth_Statics::PropPointers[] = {
	(const UECodeGen_Private::FPropertyParamsBase*)&Z_Construct_UFunction_UBasicAttributeSet_OnRep_MaxHealth_Statics::NewProp_OldValue,
};
static_assert(UE_ARRAY_COUNT(Z_Construct_UFunction_UBasicAttributeSet_OnRep_MaxHealth_Statics::PropPointers) < 2048);
const UECodeGen_Private::FFunctionParams Z_Construct_UFunction_UBasicAttributeSet_OnRep_MaxHealth_Statics::FuncParams = { { (UObject*(*)())Z_Construct_UClass_UBasicAttributeSet, nullptr, "OnRep_MaxHealth", Z_Construct_UFunction_UBasicAttributeSet_OnRep_MaxHealth_Statics::PropPointers, UE_ARRAY_COUNT(Z_Construct_UFunction_UBasicAttributeSet_OnRep_MaxHealth_Statics::PropPointers), sizeof(Z_Construct_UFunction_UBasicAttributeSet_OnRep_MaxHealth_Statics::BasicAttributeSet_eventOnRep_MaxHealth_Parms), RF_Public|RF_Transient|RF_MarkAsNative, (EFunctionFlags)0x00420401, 0, 0, METADATA_PARAMS(UE_ARRAY_COUNT(Z_Construct_UFunction_UBasicAttributeSet_OnRep_MaxHealth_Statics::Function_MetaDataParams), Z_Construct_UFunction_UBasicAttributeSet_OnRep_MaxHealth_Statics::Function_MetaDataParams)},  };
static_assert(sizeof(Z_Construct_UFunction_UBasicAttributeSet_OnRep_MaxHealth_Statics::BasicAttributeSet_eventOnRep_MaxHealth_Parms) < MAX_uint16);
UFunction* Z_Construct_UFunction_UBasicAttributeSet_OnRep_MaxHealth()
{
	static UFunction* ReturnFunction = nullptr;
	if (!ReturnFunction)
	{
		UECodeGen_Private::ConstructUFunction(&ReturnFunction, Z_Construct_UFunction_UBasicAttributeSet_OnRep_MaxHealth_Statics::FuncParams);
	}
	return ReturnFunction;
}
DEFINE_FUNCTION(UBasicAttributeSet::execOnRep_MaxHealth)
{
	P_GET_STRUCT_REF(FGameplayAttributeData,Z_Param_Out_OldValue);
	P_FINISH;
	P_NATIVE_BEGIN;
	P_THIS->OnRep_MaxHealth(Z_Param_Out_OldValue);
	P_NATIVE_END;
}
// ********** End Class UBasicAttributeSet Function OnRep_MaxHealth ********************************

// ********** Begin Class UBasicAttributeSet Function OnRep_MaxStamina *****************************
struct Z_Construct_UFunction_UBasicAttributeSet_OnRep_MaxStamina_Statics
{
	struct BasicAttributeSet_eventOnRep_MaxStamina_Parms
	{
		FGameplayAttributeData OldValue;
	};
#if WITH_METADATA
	static constexpr UECodeGen_Private::FMetaDataPairParam Function_MetaDataParams[] = {
		{ "ModuleRelativePath", "GameplayAbilitySystem/AttributeSets/BasicAttributeSet.h" },
	};
	static constexpr UECodeGen_Private::FMetaDataPairParam NewProp_OldValue_MetaData[] = {
		{ "NativeConst", "" },
	};
#endif // WITH_METADATA
	static const UECodeGen_Private::FStructPropertyParams NewProp_OldValue;
	static const UECodeGen_Private::FPropertyParamsBase* const PropPointers[];
	static const UECodeGen_Private::FFunctionParams FuncParams;
};
const UECodeGen_Private::FStructPropertyParams Z_Construct_UFunction_UBasicAttributeSet_OnRep_MaxStamina_Statics::NewProp_OldValue = { "OldValue", nullptr, (EPropertyFlags)0x0010000008000182, UECodeGen_Private::EPropertyGenFlags::Struct, RF_Public|RF_Transient|RF_MarkAsNative, nullptr, nullptr, 1, STRUCT_OFFSET(BasicAttributeSet_eventOnRep_MaxStamina_Parms, OldValue), Z_Construct_UScriptStruct_FGameplayAttributeData, METADATA_PARAMS(UE_ARRAY_COUNT(NewProp_OldValue_MetaData), NewProp_OldValue_MetaData) }; // 1532612004
const UECodeGen_Private::FPropertyParamsBase* const Z_Construct_UFunction_UBasicAttributeSet_OnRep_MaxStamina_Statics::PropPointers[] = {
	(const UECodeGen_Private::FPropertyParamsBase*)&Z_Construct_UFunction_UBasicAttributeSet_OnRep_MaxStamina_Statics::NewProp_OldValue,
};
static_assert(UE_ARRAY_COUNT(Z_Construct_UFunction_UBasicAttributeSet_OnRep_MaxStamina_Statics::PropPointers) < 2048);
const UECodeGen_Private::FFunctionParams Z_Construct_UFunction_UBasicAttributeSet_OnRep_MaxStamina_Statics::FuncParams = { { (UObject*(*)())Z_Construct_UClass_UBasicAttributeSet, nullptr, "OnRep_MaxStamina", Z_Construct_UFunction_UBasicAttributeSet_OnRep_MaxStamina_Statics::PropPointers, UE_ARRAY_COUNT(Z_Construct_UFunction_UBasicAttributeSet_OnRep_MaxStamina_Statics::PropPointers), sizeof(Z_Construct_UFunction_UBasicAttributeSet_OnRep_MaxStamina_Statics::BasicAttributeSet_eventOnRep_MaxStamina_Parms), RF_Public|RF_Transient|RF_MarkAsNative, (EFunctionFlags)0x00420401, 0, 0, METADATA_PARAMS(UE_ARRAY_COUNT(Z_Construct_UFunction_UBasicAttributeSet_OnRep_MaxStamina_Statics::Function_MetaDataParams), Z_Construct_UFunction_UBasicAttributeSet_OnRep_MaxStamina_Statics::Function_MetaDataParams)},  };
static_assert(sizeof(Z_Construct_UFunction_UBasicAttributeSet_OnRep_MaxStamina_Statics::BasicAttributeSet_eventOnRep_MaxStamina_Parms) < MAX_uint16);
UFunction* Z_Construct_UFunction_UBasicAttributeSet_OnRep_MaxStamina()
{
	static UFunction* ReturnFunction = nullptr;
	if (!ReturnFunction)
	{
		UECodeGen_Private::ConstructUFunction(&ReturnFunction, Z_Construct_UFunction_UBasicAttributeSet_OnRep_MaxStamina_Statics::FuncParams);
	}
	return ReturnFunction;
}
DEFINE_FUNCTION(UBasicAttributeSet::execOnRep_MaxStamina)
{
	P_GET_STRUCT_REF(FGameplayAttributeData,Z_Param_Out_OldValue);
	P_FINISH;
	P_NATIVE_BEGIN;
	P_THIS->OnRep_MaxStamina(Z_Param_Out_OldValue);
	P_NATIVE_END;
}
// ********** End Class UBasicAttributeSet Function OnRep_MaxStamina *******************************

// ********** Begin Class UBasicAttributeSet Function OnRep_Stamina ********************************
struct Z_Construct_UFunction_UBasicAttributeSet_OnRep_Stamina_Statics
{
	struct BasicAttributeSet_eventOnRep_Stamina_Parms
	{
		FGameplayAttributeData OldValue;
	};
#if WITH_METADATA
	static constexpr UECodeGen_Private::FMetaDataPairParam Function_MetaDataParams[] = {
		{ "ModuleRelativePath", "GameplayAbilitySystem/AttributeSets/BasicAttributeSet.h" },
	};
	static constexpr UECodeGen_Private::FMetaDataPairParam NewProp_OldValue_MetaData[] = {
		{ "NativeConst", "" },
	};
#endif // WITH_METADATA
	static const UECodeGen_Private::FStructPropertyParams NewProp_OldValue;
	static const UECodeGen_Private::FPropertyParamsBase* const PropPointers[];
	static const UECodeGen_Private::FFunctionParams FuncParams;
};
const UECodeGen_Private::FStructPropertyParams Z_Construct_UFunction_UBasicAttributeSet_OnRep_Stamina_Statics::NewProp_OldValue = { "OldValue", nullptr, (EPropertyFlags)0x0010000008000182, UECodeGen_Private::EPropertyGenFlags::Struct, RF_Public|RF_Transient|RF_MarkAsNative, nullptr, nullptr, 1, STRUCT_OFFSET(BasicAttributeSet_eventOnRep_Stamina_Parms, OldValue), Z_Construct_UScriptStruct_FGameplayAttributeData, METADATA_PARAMS(UE_ARRAY_COUNT(NewProp_OldValue_MetaData), NewProp_OldValue_MetaData) }; // 1532612004
const UECodeGen_Private::FPropertyParamsBase* const Z_Construct_UFunction_UBasicAttributeSet_OnRep_Stamina_Statics::PropPointers[] = {
	(const UECodeGen_Private::FPropertyParamsBase*)&Z_Construct_UFunction_UBasicAttributeSet_OnRep_Stamina_Statics::NewProp_OldValue,
};
static_assert(UE_ARRAY_COUNT(Z_Construct_UFunction_UBasicAttributeSet_OnRep_Stamina_Statics::PropPointers) < 2048);
const UECodeGen_Private::FFunctionParams Z_Construct_UFunction_UBasicAttributeSet_OnRep_Stamina_Statics::FuncParams = { { (UObject*(*)())Z_Construct_UClass_UBasicAttributeSet, nullptr, "OnRep_Stamina", Z_Construct_UFunction_UBasicAttributeSet_OnRep_Stamina_Statics::PropPointers, UE_ARRAY_COUNT(Z_Construct_UFunction_UBasicAttributeSet_OnRep_Stamina_Statics::PropPointers), sizeof(Z_Construct_UFunction_UBasicAttributeSet_OnRep_Stamina_Statics::BasicAttributeSet_eventOnRep_Stamina_Parms), RF_Public|RF_Transient|RF_MarkAsNative, (EFunctionFlags)0x00420401, 0, 0, METADATA_PARAMS(UE_ARRAY_COUNT(Z_Construct_UFunction_UBasicAttributeSet_OnRep_Stamina_Statics::Function_MetaDataParams), Z_Construct_UFunction_UBasicAttributeSet_OnRep_Stamina_Statics::Function_MetaDataParams)},  };
static_assert(sizeof(Z_Construct_UFunction_UBasicAttributeSet_OnRep_Stamina_Statics::BasicAttributeSet_eventOnRep_Stamina_Parms) < MAX_uint16);
UFunction* Z_Construct_UFunction_UBasicAttributeSet_OnRep_Stamina()
{
	static UFunction* ReturnFunction = nullptr;
	if (!ReturnFunction)
	{
		UECodeGen_Private::ConstructUFunction(&ReturnFunction, Z_Construct_UFunction_UBasicAttributeSet_OnRep_Stamina_Statics::FuncParams);
	}
	return ReturnFunction;
}
DEFINE_FUNCTION(UBasicAttributeSet::execOnRep_Stamina)
{
	P_GET_STRUCT_REF(FGameplayAttributeData,Z_Param_Out_OldValue);
	P_FINISH;
	P_NATIVE_BEGIN;
	P_THIS->OnRep_Stamina(Z_Param_Out_OldValue);
	P_NATIVE_END;
}
// ********** End Class UBasicAttributeSet Function OnRep_Stamina **********************************

// ********** Begin Class UBasicAttributeSet *******************************************************
void UBasicAttributeSet::StaticRegisterNativesUBasicAttributeSet()
{
	UClass* Class = UBasicAttributeSet::StaticClass();
	static const FNameNativePtrPair Funcs[] = {
		{ "OnRep_Health", &UBasicAttributeSet::execOnRep_Health },
		{ "OnRep_MaxHealth", &UBasicAttributeSet::execOnRep_MaxHealth },
		{ "OnRep_MaxStamina", &UBasicAttributeSet::execOnRep_MaxStamina },
		{ "OnRep_Stamina", &UBasicAttributeSet::execOnRep_Stamina },
	};
	FNativeFunctionRegistrar::RegisterFunctions(Class, Funcs, UE_ARRAY_COUNT(Funcs));
}
FClassRegistrationInfo Z_Registration_Info_UClass_UBasicAttributeSet;
UClass* UBasicAttributeSet::GetPrivateStaticClass()
{
	using TClass = UBasicAttributeSet;
	if (!Z_Registration_Info_UClass_UBasicAttributeSet.InnerSingleton)
	{
		GetPrivateStaticClassBody(
			StaticPackage(),
			TEXT("BasicAttributeSet"),
			Z_Registration_Info_UClass_UBasicAttributeSet.InnerSingleton,
			StaticRegisterNativesUBasicAttributeSet,
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
	return Z_Registration_Info_UClass_UBasicAttributeSet.InnerSingleton;
}
UClass* Z_Construct_UClass_UBasicAttributeSet_NoRegister()
{
	return UBasicAttributeSet::GetPrivateStaticClass();
}
struct Z_Construct_UClass_UBasicAttributeSet_Statics
{
#if WITH_METADATA
	static constexpr UECodeGen_Private::FMetaDataPairParam Class_MetaDataParams[] = {
#if !UE_BUILD_SHIPPING
		{ "Comment", "/**\n * \n */" },
#endif
		{ "IncludePath", "GameplayAbilitySystem/AttributeSets/BasicAttributeSet.h" },
		{ "ModuleRelativePath", "GameplayAbilitySystem/AttributeSets/BasicAttributeSet.h" },
	};
	static constexpr UECodeGen_Private::FMetaDataPairParam NewProp_Health_MetaData[] = {
		{ "Category", "Attributes" },
#if !UE_BUILD_SHIPPING
		{ "Comment", "// Health Attributes\n" },
#endif
		{ "ModuleRelativePath", "GameplayAbilitySystem/AttributeSets/BasicAttributeSet.h" },
#if !UE_BUILD_SHIPPING
		{ "ToolTip", "Health Attributes" },
#endif
	};
	static constexpr UECodeGen_Private::FMetaDataPairParam NewProp_MaxHealth_MetaData[] = {
		{ "Category", "Attributes" },
		{ "ModuleRelativePath", "GameplayAbilitySystem/AttributeSets/BasicAttributeSet.h" },
	};
	static constexpr UECodeGen_Private::FMetaDataPairParam NewProp_Stamina_MetaData[] = {
		{ "Category", "Attributes" },
#if !UE_BUILD_SHIPPING
		{ "Comment", "// Stamina Attributes\n" },
#endif
		{ "ModuleRelativePath", "GameplayAbilitySystem/AttributeSets/BasicAttributeSet.h" },
#if !UE_BUILD_SHIPPING
		{ "ToolTip", "Stamina Attributes" },
#endif
	};
	static constexpr UECodeGen_Private::FMetaDataPairParam NewProp_MaxStamina_MetaData[] = {
		{ "Category", "Attributes" },
		{ "ModuleRelativePath", "GameplayAbilitySystem/AttributeSets/BasicAttributeSet.h" },
	};
#endif // WITH_METADATA
	static const UECodeGen_Private::FStructPropertyParams NewProp_Health;
	static const UECodeGen_Private::FStructPropertyParams NewProp_MaxHealth;
	static const UECodeGen_Private::FStructPropertyParams NewProp_Stamina;
	static const UECodeGen_Private::FStructPropertyParams NewProp_MaxStamina;
	static const UECodeGen_Private::FPropertyParamsBase* const PropPointers[];
	static UObject* (*const DependentSingletons[])();
	static constexpr FClassFunctionLinkInfo FuncInfo[] = {
		{ &Z_Construct_UFunction_UBasicAttributeSet_OnRep_Health, "OnRep_Health" }, // 4105426027
		{ &Z_Construct_UFunction_UBasicAttributeSet_OnRep_MaxHealth, "OnRep_MaxHealth" }, // 3552342407
		{ &Z_Construct_UFunction_UBasicAttributeSet_OnRep_MaxStamina, "OnRep_MaxStamina" }, // 3387719776
		{ &Z_Construct_UFunction_UBasicAttributeSet_OnRep_Stamina, "OnRep_Stamina" }, // 3385329268
	};
	static_assert(UE_ARRAY_COUNT(FuncInfo) < 2048);
	static constexpr FCppClassTypeInfoStatic StaticCppClassTypeInfo = {
		TCppClassTypeTraits<UBasicAttributeSet>::IsAbstract,
	};
	static const UECodeGen_Private::FClassParams ClassParams;
};
const UECodeGen_Private::FStructPropertyParams Z_Construct_UClass_UBasicAttributeSet_Statics::NewProp_Health = { "Health", "OnRep_Health", (EPropertyFlags)0x0010000100000034, UECodeGen_Private::EPropertyGenFlags::Struct, RF_Public|RF_Transient|RF_MarkAsNative, nullptr, nullptr, 1, STRUCT_OFFSET(UBasicAttributeSet, Health), Z_Construct_UScriptStruct_FGameplayAttributeData, METADATA_PARAMS(UE_ARRAY_COUNT(NewProp_Health_MetaData), NewProp_Health_MetaData) }; // 1532612004
const UECodeGen_Private::FStructPropertyParams Z_Construct_UClass_UBasicAttributeSet_Statics::NewProp_MaxHealth = { "MaxHealth", "OnRep_MaxHealth", (EPropertyFlags)0x0010000100000034, UECodeGen_Private::EPropertyGenFlags::Struct, RF_Public|RF_Transient|RF_MarkAsNative, nullptr, nullptr, 1, STRUCT_OFFSET(UBasicAttributeSet, MaxHealth), Z_Construct_UScriptStruct_FGameplayAttributeData, METADATA_PARAMS(UE_ARRAY_COUNT(NewProp_MaxHealth_MetaData), NewProp_MaxHealth_MetaData) }; // 1532612004
const UECodeGen_Private::FStructPropertyParams Z_Construct_UClass_UBasicAttributeSet_Statics::NewProp_Stamina = { "Stamina", "OnRep_Stamina", (EPropertyFlags)0x0010000100000034, UECodeGen_Private::EPropertyGenFlags::Struct, RF_Public|RF_Transient|RF_MarkAsNative, nullptr, nullptr, 1, STRUCT_OFFSET(UBasicAttributeSet, Stamina), Z_Construct_UScriptStruct_FGameplayAttributeData, METADATA_PARAMS(UE_ARRAY_COUNT(NewProp_Stamina_MetaData), NewProp_Stamina_MetaData) }; // 1532612004
const UECodeGen_Private::FStructPropertyParams Z_Construct_UClass_UBasicAttributeSet_Statics::NewProp_MaxStamina = { "MaxStamina", "OnRep_MaxStamina", (EPropertyFlags)0x0010000100000034, UECodeGen_Private::EPropertyGenFlags::Struct, RF_Public|RF_Transient|RF_MarkAsNative, nullptr, nullptr, 1, STRUCT_OFFSET(UBasicAttributeSet, MaxStamina), Z_Construct_UScriptStruct_FGameplayAttributeData, METADATA_PARAMS(UE_ARRAY_COUNT(NewProp_MaxStamina_MetaData), NewProp_MaxStamina_MetaData) }; // 1532612004
const UECodeGen_Private::FPropertyParamsBase* const Z_Construct_UClass_UBasicAttributeSet_Statics::PropPointers[] = {
	(const UECodeGen_Private::FPropertyParamsBase*)&Z_Construct_UClass_UBasicAttributeSet_Statics::NewProp_Health,
	(const UECodeGen_Private::FPropertyParamsBase*)&Z_Construct_UClass_UBasicAttributeSet_Statics::NewProp_MaxHealth,
	(const UECodeGen_Private::FPropertyParamsBase*)&Z_Construct_UClass_UBasicAttributeSet_Statics::NewProp_Stamina,
	(const UECodeGen_Private::FPropertyParamsBase*)&Z_Construct_UClass_UBasicAttributeSet_Statics::NewProp_MaxStamina,
};
static_assert(UE_ARRAY_COUNT(Z_Construct_UClass_UBasicAttributeSet_Statics::PropPointers) < 2048);
UObject* (*const Z_Construct_UClass_UBasicAttributeSet_Statics::DependentSingletons[])() = {
	(UObject* (*)())Z_Construct_UClass_UAttributeSet,
	(UObject* (*)())Z_Construct_UPackage__Script_Nexus,
};
static_assert(UE_ARRAY_COUNT(Z_Construct_UClass_UBasicAttributeSet_Statics::DependentSingletons) < 16);
const UECodeGen_Private::FClassParams Z_Construct_UClass_UBasicAttributeSet_Statics::ClassParams = {
	&UBasicAttributeSet::StaticClass,
	nullptr,
	&StaticCppClassTypeInfo,
	DependentSingletons,
	FuncInfo,
	Z_Construct_UClass_UBasicAttributeSet_Statics::PropPointers,
	nullptr,
	UE_ARRAY_COUNT(DependentSingletons),
	UE_ARRAY_COUNT(FuncInfo),
	UE_ARRAY_COUNT(Z_Construct_UClass_UBasicAttributeSet_Statics::PropPointers),
	0,
	0x003000A0u,
	METADATA_PARAMS(UE_ARRAY_COUNT(Z_Construct_UClass_UBasicAttributeSet_Statics::Class_MetaDataParams), Z_Construct_UClass_UBasicAttributeSet_Statics::Class_MetaDataParams)
};
UClass* Z_Construct_UClass_UBasicAttributeSet()
{
	if (!Z_Registration_Info_UClass_UBasicAttributeSet.OuterSingleton)
	{
		UECodeGen_Private::ConstructUClass(Z_Registration_Info_UClass_UBasicAttributeSet.OuterSingleton, Z_Construct_UClass_UBasicAttributeSet_Statics::ClassParams);
	}
	return Z_Registration_Info_UClass_UBasicAttributeSet.OuterSingleton;
}
#if VALIDATE_CLASS_REPS
void UBasicAttributeSet::ValidateGeneratedRepEnums(const TArray<struct FRepRecord>& ClassReps) const
{
	static FName Name_Health(TEXT("Health"));
	static FName Name_MaxHealth(TEXT("MaxHealth"));
	static FName Name_Stamina(TEXT("Stamina"));
	static FName Name_MaxStamina(TEXT("MaxStamina"));
	const bool bIsValid = true
		&& Name_Health == ClassReps[(int32)ENetFields_Private::Health].Property->GetFName()
		&& Name_MaxHealth == ClassReps[(int32)ENetFields_Private::MaxHealth].Property->GetFName()
		&& Name_Stamina == ClassReps[(int32)ENetFields_Private::Stamina].Property->GetFName()
		&& Name_MaxStamina == ClassReps[(int32)ENetFields_Private::MaxStamina].Property->GetFName();
	checkf(bIsValid, TEXT("UHT Generated Rep Indices do not match runtime populated Rep Indices for properties in UBasicAttributeSet"));
}
#endif
DEFINE_VTABLE_PTR_HELPER_CTOR(UBasicAttributeSet);
UBasicAttributeSet::~UBasicAttributeSet() {}
// ********** End Class UBasicAttributeSet *********************************************************

// ********** Begin Registration *******************************************************************
struct Z_CompiledInDeferFile_FID_UnrealProjects_team_dynamics_project_1_Nexus_Source_Nexus_GameplayAbilitySystem_AttributeSets_BasicAttributeSet_h__Script_Nexus_Statics
{
	static constexpr FClassRegisterCompiledInInfo ClassInfo[] = {
		{ Z_Construct_UClass_UBasicAttributeSet, UBasicAttributeSet::StaticClass, TEXT("UBasicAttributeSet"), &Z_Registration_Info_UClass_UBasicAttributeSet, CONSTRUCT_RELOAD_VERSION_INFO(FClassReloadVersionInfo, sizeof(UBasicAttributeSet), 1269293096U) },
	};
};
static FRegisterCompiledInInfo Z_CompiledInDeferFile_FID_UnrealProjects_team_dynamics_project_1_Nexus_Source_Nexus_GameplayAbilitySystem_AttributeSets_BasicAttributeSet_h__Script_Nexus_641013908(TEXT("/Script/Nexus"),
	Z_CompiledInDeferFile_FID_UnrealProjects_team_dynamics_project_1_Nexus_Source_Nexus_GameplayAbilitySystem_AttributeSets_BasicAttributeSet_h__Script_Nexus_Statics::ClassInfo, UE_ARRAY_COUNT(Z_CompiledInDeferFile_FID_UnrealProjects_team_dynamics_project_1_Nexus_Source_Nexus_GameplayAbilitySystem_AttributeSets_BasicAttributeSet_h__Script_Nexus_Statics::ClassInfo),
	nullptr, 0,
	nullptr, 0);
// ********** End Registration *********************************************************************

PRAGMA_ENABLE_DEPRECATION_WARNINGS

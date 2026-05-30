// Fill out your copyright notice in the Description page of Project Settings.

#pragma once

#include "CoreMinimal.h"
#include "GameFramework/CharacterMovementComponent.h"
#include "TagDuelsCharacterMovementComponent.generated.h"


UCLASS(ClassGroup=(Custom), meta=(BlueprintSpawnableComponent))
class TAGDUELS_API UTagDuelsCharacterMovementComponent : public UCharacterMovementComponent
{
	GENERATED_BODY()

	class FTagDuelsSavedMove : public FSavedMove_Character
	{
	public:

		typedef FSavedMove_Character Super;

		///@brief Resets all saved variables.
		virtual void Clear() override;

		///@brief Store input commands in the compressed flags.
		virtual uint8 GetCompressedFlags() const override;

		///@brief This is used to check whether or not two moves can be combined into one.
		///Basically you just check to make sure that the saved variables are the same.
		virtual bool CanCombineWith(const FSavedMovePtr& NewMove, ACharacter* Character, float MaxDelta) const override;

		///@brief Sets up the move before sending it to the server. 
		virtual void SetMoveFor(ACharacter* Character, float InDeltaTime, FVector const& NewAccel, class FNetworkPredictionData_Client_Character & ClientData) override;
		///@brief Sets variables on character movement component before making a predictive correction.
		virtual void PrepMoveFor(class ACharacter* Character) override;

		// Sprint
		uint8 SavedRequestToStartSprinting : 1;
	};

	class FTagDuelsNetworkPredictionData_Client : public FNetworkPredictionData_Client_Character
	{
	public:
		FTagDuelsNetworkPredictionData_Client(const UCharacterMovementComponent& ClientMovement);

		typedef FNetworkPredictionData_Client_Character Super;

		///@brief Allocates a new copy of our custom saved move
		virtual FSavedMovePtr AllocateNewMove() override;
	};

public:
	// Sets default values for this component's properties
	UTagDuelsCharacterMovementComponent();

	UPROPERTY(EditAnywhere, BlueprintReadWrite, Category = "Sprint")
	float MaxSprintSpeed;

	UPROPERTY(EditAnywhere, BlueprintReadWrite, Category = "Net")
    float MaxSmoothNetUpdateDist;

	UPROPERTY(EditAnywhere, BlueprintReadWrite, Category = "Net")
	float NoSmoothNetUpdateDist;

	UPROPERTY(Transient, BlueprintReadOnly)
	uint8 RequestToStartSprinting : 1;

public:
	// Called every frame
	virtual void TickComponent(float DeltaTime, ELevelTick TickType,
							   FActorComponentTickFunction* ThisTickFunction) override;

	//Overrides
	virtual float GetMaxSpeed() const override;
	
	virtual void UpdateFromCompressedFlags(uint8 Flags) override;
	
	virtual class FNetworkPredictionData_Client* GetPredictionData_Client() const override;

	// Sprint
	UFUNCTION(BlueprintCallable, Category = "Sprint")
	void StartSprinting();
	
	UFUNCTION(BlueprintCallable, Category = "Sprint")
	void StopSprinting();

protected:
	// Called when the game starts
	virtual void BeginPlay() override;
};

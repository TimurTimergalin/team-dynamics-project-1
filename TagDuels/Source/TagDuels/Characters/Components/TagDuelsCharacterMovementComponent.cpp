// Fill out your copyright notice in the Description page of Project Settings.


#include "TagDuelsCharacterMovementComponent.h"
#include "TagDuels/Characters/TagDuelsCharacterBase.h"


// Sets default values for this component's properties
UTagDuelsCharacterMovementComponent::UTagDuelsCharacterMovementComponent()
{
	// Set this component to be initialized when the game starts, and to be ticked every frame.  You can turn these features
	// off to improve performance if you don't need them.
	PrimaryComponentTick.bCanEverTick = true;

	MaxSmoothNetUpdateDist = 150.f;
	NoSmoothNetUpdateDist = 200.f;
}


// Called when the game starts
void UTagDuelsCharacterMovementComponent::BeginPlay()
{
	Super::BeginPlay();
}


// Called every frame
void UTagDuelsCharacterMovementComponent::TickComponent(float DeltaTime, ELevelTick TickType,
                                                        FActorComponentTickFunction* ThisTickFunction)
{
	Super::TickComponent(DeltaTime, TickType, ThisTickFunction);
}

float UTagDuelsCharacterMovementComponent::GetMaxSpeed() const
{
	ATagDuelsCharacterBase* Owner = Cast<ATagDuelsCharacterBase>(GetOwner());
	if (!Owner)
	{
		UE_LOG(LogTemp, Error, TEXT("%s() No Owner"), *FString(__FUNCTION__));
		return Super::GetMaxSpeed();
	}

	if (Owner->GetAbilitySystemComponent()->HasMatchingGameplayTag(FGameplayTag::RequestGameplayTag(FName("State.Movement.Stunned"))))
	{
		return 0.0f;
	}

	if (IsCrouching())
	{
		return MaxWalkSpeedCrouched;
	}

	if (RequestToStartSprinting)
	{
		return MaxSprintSpeed;
	}

	return MaxWalkSpeed;
}

void UTagDuelsCharacterMovementComponent::UpdateFromCompressedFlags(uint8 Flags)
{
	Super::UpdateFromCompressedFlags(Flags);

	//The Flags parameter contains the compressed input flags that are stored in the saved move.
	//UpdateFromCompressed flags simply copies the flags from the saved move into the movement component.
	//It basically just resets the movement component to the state when the move was made so it can simulate from there.
	RequestToStartSprinting = (Flags & FSavedMove_Character::FLAG_Custom_0) != 0;
}

FNetworkPredictionData_Client * UTagDuelsCharacterMovementComponent::GetPredictionData_Client() const
{
	check(PawnOwner != NULL);

	if (!ClientPredictionData)
	{
		UTagDuelsCharacterMovementComponent* MutableThis = const_cast<UTagDuelsCharacterMovementComponent*>(this);

		MutableThis->ClientPredictionData = new FTagDuelsNetworkPredictionData_Client(*this);
		MutableThis->ClientPredictionData->MaxSmoothNetUpdateDist = MaxSmoothNetUpdateDist;
		MutableThis->ClientPredictionData->NoSmoothNetUpdateDist = NoSmoothNetUpdateDist;
	}

	return ClientPredictionData;
}

void UTagDuelsCharacterMovementComponent::StartSprinting()
{
	// Debug messages
	FString CurrentTime = FDateTime::Now().ToString(); // [citation:2]
	if (GetOwner()->GetLocalRole() == ROLE_Authority)
	{
		FString DebugText = FString::Printf(TEXT("Server(%s): StartSprint"), *CurrentTime);
		GEngine->AddOnScreenDebugMessage(-1, 5.f, FColor::Green, DebugText);
	}
	else
	{
		FString DebugText = FString::Printf(TEXT("Client(%s): StartSprint"), *CurrentTime);
		GEngine->AddOnScreenDebugMessage(-1, 5.f, FColor::Green, DebugText);
	}
	
	RequestToStartSprinting = true;
}

void UTagDuelsCharacterMovementComponent::StopSprinting()
{
	// Debug messages
	FString CurrentTime = FDateTime::Now().ToString(); // [citation:2]
	if (GetOwner()->GetLocalRole() == ROLE_Authority)
	{
		FString DebugText = FString::Printf(TEXT("Server(%s): StopSprint"), *CurrentTime);
		GEngine->AddOnScreenDebugMessage(-1, 5.f, FColor::Red, DebugText);
	}
	else
	{
		FString DebugText = FString::Printf(TEXT("Client(%s): StopSprint"), *CurrentTime);
		GEngine->AddOnScreenDebugMessage(-1, 5.f, FColor::Red, DebugText);
	}
	
	RequestToStartSprinting = false;
}

void UTagDuelsCharacterMovementComponent::FTagDuelsSavedMove::Clear()
{
	Super::Clear();

	SavedRequestToStartSprinting = false;
}

uint8 UTagDuelsCharacterMovementComponent::FTagDuelsSavedMove::GetCompressedFlags() const
{
	uint8 Result = Super::GetCompressedFlags();

	if (SavedRequestToStartSprinting)
	{
		Result |= FLAG_Custom_0;
	}

	return Result;
}

bool UTagDuelsCharacterMovementComponent::FTagDuelsSavedMove::CanCombineWith(const FSavedMovePtr & NewMove, ACharacter * Character, float MaxDelta) const
{
	//Set which moves can be combined together. This will depend on the bit flags that are used.
	if (SavedRequestToStartSprinting != ((FTagDuelsSavedMove*)&NewMove)->SavedRequestToStartSprinting)
	{
		return false;
	}

	return Super::CanCombineWith(NewMove, Character, MaxDelta);
}

void UTagDuelsCharacterMovementComponent::FTagDuelsSavedMove::SetMoveFor(ACharacter * Character, float InDeltaTime, FVector const & NewAccel, FNetworkPredictionData_Client_Character & ClientData)
{
	Super::SetMoveFor(Character, InDeltaTime, NewAccel, ClientData);

	UTagDuelsCharacterMovementComponent* CharacterMovement = Cast<UTagDuelsCharacterMovementComponent>(Character->GetCharacterMovement());
	if (CharacterMovement)
	{
		SavedRequestToStartSprinting = CharacterMovement->RequestToStartSprinting;
	}

	/*// Round acceleration, so sent version and locally used version always match
	Acceleration.X = FMath::RoundToFloat(Acceleration.X);
	Acceleration.Y = FMath::RoundToFloat(Acceleration.Y);
	Acceleration.Z = FMath::RoundToFloat(Acceleration.Z);*/
}

void UTagDuelsCharacterMovementComponent::FTagDuelsSavedMove::PrepMoveFor(ACharacter * Character)
{
	Super::PrepMoveFor(Character);

	UTagDuelsCharacterMovementComponent* CharacterMovement = Cast<UTagDuelsCharacterMovementComponent>(Character->GetCharacterMovement());
	if (CharacterMovement)
	{
		CharacterMovement->RequestToStartSprinting = SavedRequestToStartSprinting;
	}
}

UTagDuelsCharacterMovementComponent::FTagDuelsNetworkPredictionData_Client::FTagDuelsNetworkPredictionData_Client(
	const UCharacterMovementComponent& ClientMovement) : Super(ClientMovement)
{
}


FSavedMovePtr UTagDuelsCharacterMovementComponent::FTagDuelsNetworkPredictionData_Client::AllocateNewMove()
{
	return FSavedMovePtr(new FTagDuelsSavedMove());
}

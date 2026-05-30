#include "TagDuels/UI/MainMenuWidget.h"
#include "Components//Button.h"
#include "Components/EditableText.h"

UMainMenuWidget::UMainMenuWidget(const FObjectInitializer& ObjectInitializer)
	: Super(ObjectInitializer)
{
}

void UMainMenuWidget::NativeConstruct()
{
	if (ConnectButton)
	{
		ConnectButton->OnClicked.AddDynamic(this, &UMainMenuWidget::OnConnectClicked);
	}
	Super::NativeConstruct();
}

void UMainMenuWidget::NativeDestruct()
{
	Super::NativeDestruct();
}

void UMainMenuWidget::OnConnectClicked()
{
	if (!ConnectButton)
	{
		return;
	}
	if (!AddressLine)
	{
		return;
	}

	ConnectButton->SetIsEnabled(false);
	AddressLine->SetIsEnabled(false);

	FText address = AddressLine->GetText();  

	if (address.IsEmpty())
	{
		return;
	}

	APlayerController* PC = GetOwningPlayer();
	if (!PC)
	{
		return;
	}

	PC->ClientTravel(address.ToString(), TRAVEL_Absolute);
}

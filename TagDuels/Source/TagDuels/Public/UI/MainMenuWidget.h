#pragma once

#include "CoreMinimal.h"
#include "Blueprint/UserWidget.h"
#include "MainMenuWidget.generated.h"

UCLASS(Blueprintable, BlueprintType)
class TAGDUELS_API UMainMenuWidget : public UUserWidget
{
	GENERATED_BODY()
	
public:
	UMainMenuWidget(const FObjectInitializer& ObjectInitializer);
	virtual void NativeConstruct() override;
	virtual void NativeDestruct() override;
    
protected:
	UPROPERTY(BlueprintReadWrite, meta = (BindWidget))
	class UButton* ConnectButton;
	
	UPROPERTY(BlueprintReadWrite, meta = (BindWidget))
	class UEditableText* AddressLine;

	UFUNCTION()
	void OnConnectClicked();
};

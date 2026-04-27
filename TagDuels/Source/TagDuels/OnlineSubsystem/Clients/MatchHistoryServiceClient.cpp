#include "MatchHistoryServiceClient.h"

#include "HttpModule.h"
#include "Interfaces/IHttpResponse.h"
#include "TagDuels/OnlineSubsystem/utils.h"
#include "Dom/JsonObject.h"
#include "Serialization/JsonSerializer.h"

namespace
{
    // Convert a protobuf JSON match result enum to our MatchResolution.
    MatchResolution ParseMatchResult(int32 ProtoResult)
    {
        switch (ProtoResult)
        {
        case 1: return MatchResolution::FirstWins;   // MATCH_RESULT_PLAYER1_WIN
        case 2: return MatchResolution::Draw;        // MATCH_RESULT_DRAW
        case 3: return MatchResolution::SecondWins;  // MATCH_RESULT_PLAYER2_WIN
        default: return MatchResolution::Draw;       // Unspecified or Cancelled -> Draw
        }
    }

    TOptional<FMatchHistory> ConvertMatch(const TSharedPtr<FJsonObject>& MatchObj)
    {
        if (!MatchObj.IsValid())
            return {};

        FMatchHistory History;

        // --- Player1 ---
        const TSharedPtr<FJsonObject>* Player1Obj = nullptr;
        if (!MatchObj->TryGetObjectField(TEXT("player1"), Player1Obj) || !Player1Obj)
            return {};

        if (!(*Player1Obj)->TryGetNumberField(TEXT("id"), History.Player1.Id)) return {};
        if (!(*Player1Obj)->TryGetStringField(TEXT("name"), History.Player1.Name)) return {};
        (*Player1Obj)->TryGetNumberField(TEXT("rating"), History.Player1.Rating); // optional, defaults to 0

        // --- Player2 ---
        const TSharedPtr<FJsonObject>* Player2Obj = nullptr;
        if (!MatchObj->TryGetObjectField(TEXT("player2"), Player2Obj) || !Player2Obj)
            return {};

        if (!(*Player2Obj)->TryGetNumberField(TEXT("id"), History.Player2.Id)) return {};
        if (!(*Player2Obj)->TryGetStringField(TEXT("name"), History.Player2.Name)) return {};
        (*Player2Obj)->TryGetNumberField(TEXT("rating"), History.Player2.Rating);

        // --- Rounds ---
        const TArray<TSharedPtr<FJsonValue>>* RoundsArray = nullptr;
        if (MatchObj->TryGetArrayField(TEXT("rounds"), RoundsArray))
        {
            for (const auto& RoundVal : *RoundsArray)
            {
                const TSharedPtr<FJsonObject>* RoundObj = nullptr;
                if (!RoundVal->TryGetObject(RoundObj) || !RoundObj)
                    continue;

                FRoundData Round;
                bool bIsPlayer1Killer = false;
                if ((*RoundObj)->TryGetBoolField(TEXT("isPlayer1Killer"), bIsPlayer1Killer))
                {
                    Round.RoundKiller_ = bIsPlayer1Killer ? RoundKiller::First : RoundKiller::Second;
                }

                int64 TimeMillis = 0;
                if ((*RoundObj)->TryGetNumberField(TEXT("timeMillis"), TimeMillis))
                {
                    Round.Duration = FTimespan::FromMilliseconds(TimeMillis);
                }

                History.Rounds.Add(Round);
            }
        }

        // --- End Timestamp ---
        int64 EndTimestamp = 0;
        if (!MatchObj->TryGetNumberField(TEXT("endTimestamp"), EndTimestamp))
            return {};

        // Convert milliseconds since epoch to FDateTime.
        History.EndTime = FDateTime::FromUnixTimestamp(EndTimestamp / 1000) +
                          FTimespan::FromMilliseconds(EndTimestamp % 1000);

        // --- Match Result ---
        int32 ResultVal = 0;
        if (MatchObj->TryGetNumberField(TEXT("matchResult"), ResultVal))
        {
            History.Resolution = ParseMatchResult(ResultVal);
        }

        // --- Match ID ---
        if (!MatchObj->TryGetStringField(TEXT("matchId"), History.MatchId))
            return {};

        return History;
    }

    TOptional<FMatchHistoryPage> ConvertMatchHistoryResponse(FHttpResponsePtr Response)
    {
        if (!Response.IsValid())
        {
            UE_LOG(LogTemp, Error, TEXT("GetMatchHistory: invalid response"));
            return {};
        }

        if (Response->GetResponseCode() != 200)
        {
            UE_LOG(LogTemp, Error, TEXT("GetMatchHistory: HTTP error %d, body: %s"),
                   Response->GetResponseCode(), *Response->GetContentAsString());
            return {};
        }

        FString JsonString = Response->GetContentAsString();
        TSharedPtr<FJsonObject> JsonObject;
        TSharedRef<TJsonReader<>> Reader = TJsonReaderFactory<>::Create(JsonString);

        if (!FJsonSerializer::Deserialize(Reader, JsonObject) || !JsonObject.IsValid())
        {
            UE_LOG(LogTemp, Error, TEXT("GetMatchHistory: failed to parse JSON: %s"), *JsonString);
            return {};
        }

        FMatchHistoryPage Page;

        JsonObject->TryGetStringField(TEXT("pagekey"), Page.NextPageKey);

        const TArray<TSharedPtr<FJsonValue>>* MatchesArray = nullptr;
        if (!JsonObject->TryGetArrayField(TEXT("matches"), MatchesArray))
        {
            return Page;
        }

        for (const auto& MatchValue : *MatchesArray)
        {
            const TSharedPtr<FJsonObject>* MatchObj = nullptr;
            if (!MatchValue->TryGetObject(MatchObj) || !MatchObj)
                continue;

            auto Converted = ConvertMatch(*MatchObj);
            if (Converted.IsSet())
            {
                Page.Matches.Add(Converted.GetValue());
            }
            else
            {
                UE_LOG(LogTemp, Warning, TEXT("GetMatchHistory: failed to convert a match, skipping"));
            }
        }

        return Page;
    }
}

MatchHistoryServiceClient::MatchHistoryServiceClient(const FString& Address)
    : Address(Address)
{
}

TSharedPtr<IHttpRequest> MatchHistoryServiceClient::CreateMatchHistoryRequest(int64 UserId, const FString& PageKey) const
{
    FString Url = FString::Printf(TEXT("http://%s/v1/match_history?user_id=%lld"), *Address, UserId);
    if (!PageKey.IsEmpty())
    {
        Url += FString::Printf(TEXT("&pagekey=%s"), *PageKey);
    }

    TSharedRef<IHttpRequest> Request = FHttpModule::Get().CreateRequest();
    Request->SetURL(Url);
    Request->SetVerb(TEXT("GET"));
    Request->SetHeader(TEXT("Accept"), TEXT("application/json"));

    return Request;
}

TFuture<TOptional<FMatchHistoryPage>> MatchHistoryServiceClient::GetMatchHistory(int64 UserId, const FString& PageKey) const
{
    TSharedPtr<IHttpRequest> Request = CreateMatchHistoryRequest(UserId, PageKey);
    return MakeHttpRequest(Request).Next(ConvertMatchHistoryResponse);
}

TOptional<MatchHistoryServiceClient> CreateMatchHistoryServiceClient()
{
    FString Address;
    if (!GConfig)
    {
        UE_LOG(LogTemp, Error, TEXT("CreateMatchHistoryServiceClient: GConfig absent"));
        return {};
    }
    if (!GConfig->GetString(TEXT("OnlineSubsystemAddresses"), TEXT("MatchHistoryService"), Address, GGameIni))
    {
        UE_LOG(LogTemp, Error, TEXT("CreateMatchHistoryServiceClient: OnlineSubsystemAddresses:MatchHistoryService not set"));
        return {};
    }

    return MatchHistoryServiceClient(Address);
}

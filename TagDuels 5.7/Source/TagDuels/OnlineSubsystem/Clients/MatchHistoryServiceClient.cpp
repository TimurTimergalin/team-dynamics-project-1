#include "MatchHistoryServiceClient.h"

#include "HttpModule.h"
#include "Interfaces/IHttpResponse.h"
#include "TagDuels/OnlineSubsystem/utils.h"
#include "Dom/JsonObject.h"
#include "Serialization/JsonSerializer.h"

namespace
{
    // Convert a protobuf JSON match result enum to our MatchResolution.
    MatchResolution ParseMatchResult(const FString& ProtoResult)
    {
        if (ProtoResult == TEXT("MATCH_RESULT_PLAYER1_WIN")) return MatchResolution::FirstWins;
        if (ProtoResult == TEXT("MATCH_RESULT_PLAYER2_WIN")) return MatchResolution::SecondWins;
        if (ProtoResult == TEXT("MATCH_RESULT_DRAW"))        return MatchResolution::Draw;
        return MatchResolution::Draw;
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

        FString IdStr;
        if (!(*Player1Obj)->TryGetStringField(TEXT("id"), IdStr)) return {};
        TOptional<int64> P1Id = StrToInt64(IdStr);
        if (!P1Id.IsSet()) return {};
        History.Player1.Id = P1Id.GetValue();
        if (!(*Player1Obj)->TryGetStringField(TEXT("name"), History.Player1.Name)) return {};
        FString RatingStr;
        if ((*Player1Obj)->TryGetStringField(TEXT("rating"), RatingStr))
        {
            TOptional<int64> R = StrToInt64(RatingStr);
            if (R.IsSet()) History.Player1.Rating = R.GetValue();
        }

        // --- Player2 ---
        const TSharedPtr<FJsonObject>* Player2Obj = nullptr;
        if (!MatchObj->TryGetObjectField(TEXT("player2"), Player2Obj) || !Player2Obj)
            return {};

        if (!(*Player2Obj)->TryGetStringField(TEXT("id"), IdStr)) return {};
        TOptional<int64> P2Id = StrToInt64(IdStr);
        if (!P2Id.IsSet()) return {};
        History.Player2.Id = P2Id.GetValue();
        if (!(*Player2Obj)->TryGetStringField(TEXT("name"), History.Player2.Name)) return {};
        if ((*Player2Obj)->TryGetStringField(TEXT("rating"), RatingStr))
        {
            TOptional<int64> R = StrToInt64(RatingStr);
            if (R.IsSet()) History.Player2.Rating = R.GetValue();
        }

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
        FString EndTimestampStr;
        if (!MatchObj->TryGetStringField(TEXT("endTimestamp"), EndTimestampStr))
            return {};
        TOptional<int64> EndTimestamp = StrToInt64(EndTimestampStr);
        if (!EndTimestamp.IsSet()) return {};
        History.EndTime = FDateTime::FromUnixTimestamp(EndTimestamp.GetValue() / 1000) +
                          FTimespan::FromMilliseconds(EndTimestamp.GetValue() % 1000);

        // --- Match Result ---
        FString ResultStr;
        if (MatchObj->TryGetStringField(TEXT("matchResult"), ResultStr))
            History.Resolution = ParseMatchResult(ResultStr);

        // --- Match ID (optional) ---
        MatchObj->TryGetStringField(TEXT("matchId"), History.MatchId);

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

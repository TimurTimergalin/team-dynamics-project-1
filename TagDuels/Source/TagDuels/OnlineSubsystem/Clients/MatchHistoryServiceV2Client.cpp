#include "MatchHistoryServiceV2Client.h"

#include "HttpModule.h"
#include "Interfaces/IHttpResponse.h"
#include "TagDuels/OnlineSubsystem/utils.h"
#include "Dom/JsonObject.h"
#include "Serialization/JsonSerializer.h"

namespace
{
    MatchResolution ParseMatchResult(const FString& ProtoResult)
    {
        if (ProtoResult == TEXT("MATCH_RESULT_PLAYER1_WIN")) return MatchResolution::FirstWins;
        if (ProtoResult == TEXT("MATCH_RESULT_PLAYER2_WIN")) return MatchResolution::SecondWins;
        return MatchResolution::Draw;
    }

    TOptional<FMatchHistory> ConvertMatch(const TSharedPtr<FJsonObject>& MatchObj)
    {
        if (!MatchObj.IsValid()) return {};

        FMatchHistory History;

        auto ParsePlayer = [&](const FString& Field, FMatchHistoryPlayerData& Player) -> bool
        {
            const TSharedPtr<FJsonObject>* PlayerObj = nullptr;
            if (!MatchObj->TryGetObjectField(Field, PlayerObj) || !PlayerObj) return false;

            FString IdStr;
            if (!(*PlayerObj)->TryGetStringField(TEXT("id"), IdStr)) return false;
            TOptional<int64> Id = StrToInt64(IdStr);
            if (!Id.IsSet()) return false;
            Player.Id = Id.GetValue();
            (*PlayerObj)->TryGetStringField(TEXT("name"), Player.Name);
            FString RatingStr;
            if ((*PlayerObj)->TryGetStringField(TEXT("rating"), RatingStr))
            {
                TOptional<int64> R = StrToInt64(RatingStr);
                if (R.IsSet()) Player.Rating = R.GetValue();
            }
            return true;
        };

        if (!ParsePlayer(TEXT("player1"), History.Player1)) return {};
        if (!ParsePlayer(TEXT("player2"), History.Player2)) return {};

        const TArray<TSharedPtr<FJsonValue>>* RoundsArray = nullptr;
        if (MatchObj->TryGetArrayField(TEXT("rounds"), RoundsArray))
        {
            for (const auto& RoundVal : *RoundsArray)
            {
                const TSharedPtr<FJsonObject>* RoundObj = nullptr;
                if (!RoundVal->TryGetObject(RoundObj) || !RoundObj) continue;

                FRoundData Round;
                bool bIsPlayer1Killer = false;
                if ((*RoundObj)->TryGetBoolField(TEXT("isPlayer1Killer"), bIsPlayer1Killer))
                    Round.RoundKiller_ = bIsPlayer1Killer ? RoundKiller::First : RoundKiller::Second;

                int64 TimeMillis = 0;
                if ((*RoundObj)->TryGetNumberField(TEXT("timeMillis"), TimeMillis))
                    Round.Duration = FTimespan::FromMilliseconds(TimeMillis);

                History.Rounds.Add(Round);
            }
        }

        FString EndTimestampStr;
        if (!MatchObj->TryGetStringField(TEXT("endTimestamp"), EndTimestampStr)) return {};
        TOptional<int64> EndTimestamp = StrToInt64(EndTimestampStr);
        if (!EndTimestamp.IsSet()) return {};
        History.EndTime = FDateTime::FromUnixTimestamp(EndTimestamp.GetValue() / 1000) +
                          FTimespan::FromMilliseconds(EndTimestamp.GetValue() % 1000);

        FString ResultStr;
        if (MatchObj->TryGetStringField(TEXT("matchResult"), ResultStr))
            History.Resolution = ParseMatchResult(ResultStr);

        MatchObj->TryGetStringField(TEXT("matchId"), History.MatchId);

        return History;
    }

    TOptional<FMatchHistoryPage> ConvertGetMatchHistoryResponse(FHttpResponsePtr Response)
    {
        if (!Response.IsValid())
        {
            UE_LOG(LogTemp, Error, TEXT("V2 GetMatchHistory: invalid response"));
            return {};
        }
        if (Response->GetResponseCode() != 200)
        {
            UE_LOG(LogTemp, Error, TEXT("V2 GetMatchHistory: HTTP error %d, body: %s"),
                   Response->GetResponseCode(), *Response->GetContentAsString());
            return {};
        }

        TSharedPtr<FJsonObject> JsonObject;
        TSharedRef<TJsonReader<>> Reader = TJsonReaderFactory<>::Create(Response->GetContentAsString());
        if (!FJsonSerializer::Deserialize(Reader, JsonObject) || !JsonObject.IsValid())
        {
            UE_LOG(LogTemp, Error, TEXT("V2 GetMatchHistory: failed to parse JSON"));
            return {};
        }

        FMatchHistoryPage Page;
        JsonObject->TryGetStringField(TEXT("pagekey"), Page.NextPageKey);

        const TArray<TSharedPtr<FJsonValue>>* MatchesArray = nullptr;
        if (JsonObject->TryGetArrayField(TEXT("matches"), MatchesArray))
        {
            for (const auto& MatchValue : *MatchesArray)
            {
                const TSharedPtr<FJsonObject>* MatchObj = nullptr;
                if (!MatchValue->TryGetObject(MatchObj) || !MatchObj) continue;

                auto Converted = ConvertMatch(*MatchObj);
                if (Converted.IsSet())
                    Page.Matches.Add(Converted.GetValue());
                else
                    UE_LOG(LogTemp, Warning, TEXT("V2 GetMatchHistory: failed to convert a match, skipping"));
            }
        }

        return Page;
    }

    TOptional<int64> ConvertGetRatingResponse(FHttpResponsePtr Response)
    {
        if (!Response.IsValid())
        {
            UE_LOG(LogTemp, Error, TEXT("V2 GetRating: invalid response"));
            return {};
        }
        if (Response->GetResponseCode() != 200)
        {
            UE_LOG(LogTemp, Error, TEXT("V2 GetRating: HTTP error %d, body: %s"),
                   Response->GetResponseCode(), *Response->GetContentAsString());
            return {};
        }

        TSharedPtr<FJsonObject> JsonObject;
        TSharedRef<TJsonReader<>> Reader = TJsonReaderFactory<>::Create(Response->GetContentAsString());
        if (!FJsonSerializer::Deserialize(Reader, JsonObject) || !JsonObject.IsValid())
        {
            UE_LOG(LogTemp, Error, TEXT("V2 GetRating: failed to parse JSON"));
            return {};
        }

        const TSharedPtr<FJsonObject>* RatingObj = nullptr;
        if (!JsonObject->TryGetObjectField(TEXT("rating"), RatingObj) || !RatingObj)
        {
            UE_LOG(LogTemp, Error, TEXT("V2 GetRating: missing 'rating' field"));
            return {};
        }

        FString DisplayStr;
        if (!(*RatingObj)->TryGetStringField(TEXT("displayValue"), DisplayStr))
        {
            UE_LOG(LogTemp, Error, TEXT("V2 GetRating: missing 'displayValue' field"));
            return {};
        }

        return StrToInt64(DisplayStr);
    }
}

MatchHistoryServiceV2Client::MatchHistoryServiceV2Client(const FString& Address)
    : Address(Address)
{
}

TSharedPtr<IHttpRequest> MatchHistoryServiceV2Client::CreateGetMatchHistoryRequest(int64 UserId, const FString& PageKey) const
{
    FString Url = FString::Printf(TEXT("http://%s/v1/match_history?user_id=%lld"), *Address, UserId);
    if (!PageKey.IsEmpty())
        Url += FString::Printf(TEXT("&pagekey=%s"), *PageKey);

    TSharedRef<IHttpRequest> Request = FHttpModule::Get().CreateRequest();
    Request->SetURL(Url);
    Request->SetVerb(TEXT("GET"));
    Request->SetHeader(TEXT("Accept"), TEXT("application/json"));
    return Request;
}

TSharedPtr<IHttpRequest> MatchHistoryServiceV2Client::CreateGetRatingRequest(int64 UserId) const
{
    TSharedRef<IHttpRequest> Request = FHttpModule::Get().CreateRequest();
    Request->SetURL(FString::Printf(TEXT("http://%s/v1/ratings/%lld"), *Address, UserId));
    Request->SetVerb(TEXT("GET"));
    Request->SetHeader(TEXT("Accept"), TEXT("application/json"));
    return Request;
}

TFuture<TOptional<FMatchHistoryPage>> MatchHistoryServiceV2Client::GetMatchHistory(int64 UserId, const FString& PageKey) const
{
    return MakeHttpRequest(CreateGetMatchHistoryRequest(UserId, PageKey)).Next(ConvertGetMatchHistoryResponse);
}

TFuture<TOptional<int64>> MatchHistoryServiceV2Client::GetRating(int64 UserId) const
{
    return MakeHttpRequest(CreateGetRatingRequest(UserId)).Next(ConvertGetRatingResponse);
}

TOptional<MatchHistoryServiceV2Client> CreateMatchHistoryServiceV2Client()
{
    if (!GConfig)
    {
        UE_LOG(LogTemp, Error, TEXT("CreateMatchHistoryServiceV2Client: GConfig absent"));
        return {};
    }
    FString Address;
    if (!GConfig->GetString(TEXT("OnlineSubsystemAddresses"), TEXT("MatchHistoryServiceV2"), Address, GGameIni))
    {
        UE_LOG(LogTemp, Error, TEXT("CreateMatchHistoryServiceV2Client: OnlineSubsystemAddresses:MatchHistoryServiceV2 not set"));
        return {};
    }
    return MatchHistoryServiceV2Client(Address);
}

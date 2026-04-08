package redis

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"logging"
	"strconv"
	"team_dynamics/match_service/models"
)

type CreateMatchResult struct {
	MatchId *string
	Fail1   *models.PlayerFailResponse
	Fail2   *models.PlayerFailResponse
}

type MatchKvRepo interface {
	SaveCreateMatch(ctx context.Context, match models.Match, player1, player2 models.Player) (CreateMatchResult, error)
	GetMatchByPlayerId(ctx context.Context, playerId int64) (*models.Match, *models.Player, *models.Player, error)
	GetMatchById(ctx context.Context, matchId string) (*models.Match, *models.Player, *models.Player, error)
	SaveMatchStart(ctx context.Context, matchId string) error
	SaveMatchFinish(ctx context.Context, matchId string) error
	RemoveMatch(ctx context.Context, matchId string) error
	RestartMatch(ctx context.Context, matchId string) error
}

type matchKvRepoImpl struct {
	rdb               *redis.Client
	createMatchScript *redis.Script
}

func MakeMatchKvRepo(rdb *redis.Client) MatchKvRepo {
	return &matchKvRepoImpl{
		rdb:               rdb,
		createMatchScript: redis.NewScript(createMatchScriptCode),
	}
}

var createMatchScriptCode = `
local MatchStatus = {
    REQUESTED = '1',
    ONGOING   = '2',
    FINISHED  = '3',
}

local PlayerFailResponse = {
    REENTER = "1",
    REMOVE  = "2",
}

local p1_keys = {  -- player1
    match_id = KEYS[1],
    reg_id   = KEYS[2],
    rating   = KEYS[3],
    name     = KEYS[4],
}

local p2_keys = {  -- player2
    match_id = KEYS[5],
    reg_id   = KEYS[6],
    rating   = KEYS[7],
    name     = KEYS[8],
}

local m_keys = {  -- match
    status  = KEYS[9],
    fleet   = KEYS[10],
    player1 = KEYS[11],
    player2 = KEYS[12],
}

local p1cm_keys = {  -- player1's current match
    status  = KEYS[13],
    fleet   = KEYS[14],
    player1 = KEYS[15],
    player2 = KEYS[16],
}

local p2cm_keys = {  -- player2's current match
    status  = KEYS[17],
    fleet   = KEYS[18],
    player1 = KEYS[19],
    player2 = KEYS[20],
}

local args = {
    p1_reg_id = ARGV[1],
    p1_rating = ARGV[2],
    p1_name   = ARGV[3],
    p2_reg_id = ARGV[4],
    p2_rating = ARGV[5],
    p2_name   = ARGV[6],
    m_id      = ARGV[7],
    m_fleet   = ARGV[8],
    p1_id     = ARGV[9],
    p2_id     = ARGV[10],
}

local p1_new_vals = {
    match_id = args.m_id,
    reg_id   = args.p1_reg_id,
    rating   = args.p1_rating,
    name     = args.p1_name,
}

local p2_new_vals = {
    match_id = args.m_id,
    reg_id   = args.p2_reg_id,
    rating   = args.p2_rating,
    name     = args.p2_name,
}

local m_new_vals = {
    status  = MatchStatus.REQUESTED,
    fleet   = args.m_fleet,
    player1 = args.p1_id,
    player2 = args.p2_id,
}

local function clear_match(m_keys_)
    redis.call('DEL', m_keys_.status, m_keys_.fleet, m_keys_.player1, m_keys_.player2)
end

local function write_match(m_keys_, m_vals)
    for _, field in ipairs({'status', 'fleet', 'player1', 'player2'}) do
        redis.call('SET', m_keys_[field], m_vals[field])
    end
end

local function write_player(p_keys, p_vals)
    for _, field in ipairs({'match_id', 'reg_id', 'rating', 'name'}) do
        redis.call('SET', p_keys[field], p_vals[field])
    end
end

local player1_reg_id_saved = redis.call('GET', p1_keys.reg_id)
local player2_reg_id_saved = redis.call('GET', p2_keys.reg_id)

local player1_repeated = player1_reg_id_saved and (player1_reg_id_saved == p1_new_vals.reg_id)
local player2_repeated = player2_reg_id_saved and (player2_reg_id_saved == p2_new_vals.reg_id)

if player1_repeated or player2_repeated then
    local player1_response = player1_repeated and PlayerFailResponse.REMOVE or PlayerFailResponse.REENTER
    local player2_response = player2_repeated and PlayerFailResponse.REMOVE or PlayerFailResponse.REENTER
    return {1, player1_response, player2_response}
end

local player1_has_match = (p1cm_keys.status ~= "") and (redis.call('GET', p1cm_keys.status) ~= MatchStatus.FINISHED)
local player2_has_match = (p2cm_keys.status ~= "") and (redis.call('GET', p2cm_keys.status) ~= MatchStatus.FINISHED)

if player1_has_match or player2_has_match then
    local player1_response = player1_has_match and PlayerFailResponse.REMOVE or PlayerFailResponse.REENTER
    local player2_response = player2_has_match and PlayerFailResponse.REMOVE or PlayerFailResponse.REENTER
    return {2, player1_response, player2_response}
end

if player1_has_match and p1cm_keys.status ~= "" then
    clear_match(p1cm_keys)
end
if player2_has_match and p2cm_keys.status ~= "" then
    clear_match(p2cm_keys)
end

write_match(m_keys, m_new_vals)
write_player(p1_keys, p1_new_vals)
write_player(p2_keys, p2_new_vals)

return {0}
`

func makeMatchId() string {
	return uuid.New().String()
}

func ParseMatchStatus(raw string) (models.MatchStatus, error) {
	res64, err := strconv.ParseInt(raw, 10, 32)
	if err != nil {
		return 0, err
	}
	res32 := int32(res64)
	return models.MatchStatus(res32), nil
}

func ParseFailResponse(raw string) (models.PlayerFailResponse, error) {
	res64, err := strconv.ParseInt(raw, 10, 32)
	if err != nil {
		return 0, err
	}
	res32 := int32(res64)
	return models.PlayerFailResponse(res32), nil
}

func (r *matchKvRepoImpl) readPlayer(ctx context.Context, keys playerKeys) (*models.Player, error) {
	resultSl, err := r.rdb.MGet(ctx, keys.keys()...).Result()
	if err != nil {
		return nil, err
	}
	if resultSl[0] == nil || resultSl[1] == nil || resultSl[2] == nil || resultSl[3] == nil {
		return nil, redis.Nil
	}
	rating, err := strconv.ParseInt(resultSl[2].(string), 10, 64)
	if err != nil {
		panic("player rating is not a number")
	}
	return &models.Player{
		Id:     keys.id,
		Name:   resultSl[3].(string),
		Rating: rating,
		RegId:  resultSl[1].(string),
	}, nil
}

func (r *matchKvRepoImpl) readMatch(ctx context.Context, keys realMatchKeys) (*models.Match, error) {
	resultSl, err := r.rdb.MGet(ctx, keys.keys()...).Result()
	if err != nil {
		return nil, err
	}
	if resultSl[0] == nil || resultSl[1] == nil || resultSl[2] == nil || resultSl[3] == nil {
		return nil, redis.Nil
	}
	status, err := ParseMatchStatus(resultSl[0].(string))
	if err != nil {
		panic("match status is not a number")
	}
	player1Id, err := strconv.ParseInt(resultSl[2].(string), 10, 64)
	if err != nil {
		panic("player1 id is not a number")
	}
	player2Id, err := strconv.ParseInt(resultSl[4].(string), 10, 64)
	if err != nil {
		panic("player2 id is not a number")
	}
	return &models.Match{
		MatchId:   keys.id,
		Player1Id: player1Id,
		Player2Id: player2Id,
		Status:    status,
		Fleet:     resultSl[1].(string),
	}, nil
}

func (r *matchKvRepoImpl) SaveCreateMatch(ctx context.Context, match models.Match, player1, player2 models.Player) (CreateMatchResult, error) {
	match.MatchId = makeMatchId()
	p1Keys := playerKeys{player1.Id}
	p2Keys := playerKeys{player2.Id}
	mKeys := realMatchKeys{match.MatchId}
	var result CreateMatchResult
	err := r.rdb.Watch(ctx, func(tx *redis.Tx) error {
		p1cmId, err := tx.Get(ctx, p1Keys.matchId()).Result()
		if err != nil && !errors.Is(err, redis.Nil) {
			return err
		}
		if errors.Is(err, redis.Nil) {
			p1cmId = ""
		}
		p2cmId, err := tx.Get(ctx, p2Keys.matchId()).Result()
		if err != nil && !errors.Is(err, redis.Nil) {
			return err
		}
		if errors.Is(err, redis.Nil) {
			p2cmId = ""
		}

		var p1cmKeys matchKeys
		if p1cmId == "" {
			p1cmKeys = absentMatchKeys{}
		} else {
			p1cmKeys = realMatchKeys{p1cmId}
		}
		var p2cmKeys matchKeys
		if p2cmId == "" {
			p2cmKeys = absentMatchKeys{}
		} else {
			p2cmKeys = realMatchKeys{p2cmId}
		}

		keys := []string{
			p1Keys.matchId(),
			p1Keys.regId(),
			p1Keys.rating(),
			p1Keys.name(),
			p2Keys.matchId(),
			p2Keys.regId(),
			p2Keys.rating(),
			p2Keys.name(),
			mKeys.status(),
			mKeys.fleet(),
			mKeys.player1(),
			mKeys.player2(),
			p1cmKeys.status(),
			p1cmKeys.fleet(),
			p1cmKeys.player1(),
			p1cmKeys.player2(),
			p2cmKeys.status(),
			p2cmKeys.fleet(),
			p2cmKeys.player1(),
			p2cmKeys.player2(),
		}

		args := []interface{}{
			player1.RegId,
			player1.Rating,
			player1.Name,
			player2.RegId,
			player2.Rating,
			player2.Name,
			match.MatchId,
			match.Fleet,
			match.Player1Id,
			match.Player2Id,
		}
		var scriptRes *redis.Cmd
		_, err = tx.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
			scriptRes = r.createMatchScript.Run(ctx, pipe, keys, args...)
			return nil
		})
		if err != nil {
			return err
		}
		res, err := scriptRes.Result()
		if err != nil {
			return err
		}
		resSl := res.([]interface{})
		if len(resSl) == 0 {
			panic("Empty response from script")
		}
		scriptStatus := resSl[0].(string)
		if scriptStatus == "0" {
			result.MatchId = &match.MatchId
		} else {
			p1Resp, err := ParseFailResponse(resSl[1].(string))
			if err != nil {
				panic("Player1 fail resp is not a number")
			}
			p2Resp, err := ParseFailResponse(resSl[2].(string))
			if err != nil {
				panic("Player2 fail resp is not a number")
			}
			result.Fail1 = &p1Resp
			result.Fail2 = &p2Resp
		}
		return nil
	}, p1Keys.matchId(), p2Keys.matchId())
	return result, err
}

func (r *matchKvRepoImpl) GetMatchById(ctx context.Context, matchId string) (*models.Match, *models.Player, *models.Player, error) {
	mKeys := realMatchKeys{matchId}
	player1Id, err := r.rdb.Get(ctx, mKeys.player1()).Int64()
	if err != nil {
		return nil, nil, nil, err
	}
	player2Id, err := r.rdb.Get(ctx, mKeys.player1()).Int64()
	if err != nil {
		return nil, nil, nil, err
	}

	player1, err := r.readPlayer(ctx, playerKeys{player1Id})
	if err != nil {
		return nil, nil, nil, err
	}
	player2, err := r.readPlayer(ctx, playerKeys{player2Id})
	if err != nil {
		return nil, nil, nil, err
	}
	match, err := r.readMatch(ctx, mKeys)
	if err != nil {
		return nil, nil, nil, err
	}

	return match, player1, player2, nil
}

func (r *matchKvRepoImpl) GetMatchByPlayerId(ctx context.Context, playerId int64) (*models.Match, *models.Player, *models.Player, error) {
	pKeys := playerKeys{playerId}
	matchId, err := r.rdb.Get(ctx, pKeys.matchId()).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil, nil, nil
		}
		return nil, nil, nil, err
	}
	return r.GetMatchById(ctx, matchId)
}

func (r *matchKvRepoImpl) changeMatchStatus(ctx context.Context, matchId string, from models.MatchStatus, to models.MatchStatus) error {
	logger := logging.GetLogger(ctx)
	mKeys := realMatchKeys{matchId}
	err := r.rdb.Watch(ctx, func(tx *redis.Tx) error {
		statusRaw, err := tx.Get(ctx, mKeys.status()).Int()
		if errors.Is(err, redis.Nil) {
			logger.Warn("Trying to change cancelled match's status", "matchId", matchId)
			return nil
		}
		if err != nil {
			return err
		}
		status := models.MatchStatus(int32(statusRaw))
		if status != from {
			logger.Warn("Wrong status to change", "matchId", matchId)
			return nil
		}
		_, err = tx.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
			pipe.Set(ctx, mKeys.status(), strconv.FormatInt(int64(to), 10), 0)
			return nil
		})
		return err
	}, mKeys.status())
	return err
}

func (r *matchKvRepoImpl) SaveMatchStart(ctx context.Context, matchId string) error {
	return r.changeMatchStatus(ctx, matchId, models.Requested, models.Ongoing)
}

func (r *matchKvRepoImpl) SaveMatchFinish(ctx context.Context, matchId string) error {
	return r.changeMatchStatus(ctx, matchId, models.Ongoing, models.Finished)
}

func (r *matchKvRepoImpl) RemoveMatch(ctx context.Context, matchId string) error {
	mKeys := realMatchKeys{matchId}
	err := r.rdb.Watch(ctx, func(tx *redis.Tx) error {
		player1, err := tx.Get(ctx, mKeys.player1()).Int64()
		if err != nil {
			if errors.Is(err, redis.Nil) {
				return nil
			}
			return fmt.Errorf("error while getting player1 id: %w", err)
		}
		player2, err := tx.Get(ctx, mKeys.player2()).Int64()
		if err != nil {
			if errors.Is(err, redis.Nil) {
				return nil
			}
			return fmt.Errorf("error while getting player2 id: %w", err)
		}
		p1Keys := playerKeys{player1}
		p2Keys := playerKeys{player2}
		keysToDelete := append(append(mKeys.keys(), p1Keys.keys()...), p2Keys.keys()...)
		_, err = tx.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
			pipe.Del(ctx, keysToDelete...)
			return nil
		})
		return err
	}, mKeys.player1(), mKeys.player2())
	return err
}

func (r *matchKvRepoImpl) RestartMatch(ctx context.Context, matchId string) error {
	logger := logging.GetLogger(ctx)
	mKeys := realMatchKeys{matchId}
	err := r.rdb.Watch(ctx, func(tx *redis.Tx) error {
		player1, err := tx.Get(ctx, mKeys.player1()).Int64()
		if err != nil {
			return fmt.Errorf("error while getting player1 id: %w", err)
		}
		player2, err := tx.Get(ctx, mKeys.player2()).Int64()
		if err != nil {
			return fmt.Errorf("error while getting player2 id: %w", err)
		}
		status, err := tx.Get(ctx, mKeys.status()).Int()
		if err != nil {
			return fmt.Errorf("error while getting status id: %w", err)
		}
		if status != int(models.Finished) {
			logger.Info("trying to renew non-finished match", "matchId", matchId)
			return nil
		}

		fleet, err := tx.Get(ctx, mKeys.fleet()).Result()
		if err != nil {
			return fmt.Errorf("error while getting fleet: %w", err)
		}
		p1Keys := playerKeys{player1}
		p2Keys := playerKeys{player2}
		newMatchId := makeMatchId()
		newMKeys := realMatchKeys{newMatchId}
		_, err = tx.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
			pipe.Del(ctx, mKeys.keys()...)
			pipe.MSet(ctx, []interface{}{
				newMKeys.status(), models.Ongoing,
				newMKeys.fleet(), fleet,
				newMKeys.player1(), player1,
				newMKeys.player2(), player2,
				p1Keys.matchId(), newMatchId,
				p2Keys.matchId(), newMatchId,
			})
			return nil
		})
		return err
	}, mKeys.player1(), mKeys.player2(), mKeys.status(), mKeys.fleet())
	return err
}

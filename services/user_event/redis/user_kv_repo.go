package redis

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"log/slog"
	"team_dynamics/user_event/models"
)

type UserKvRepo interface {
	Register(ctx context.Context, player *models.PlayerUserData) (bool, error)
	Unregister(ctx context.Context, playerId int64) error
	Subscribe(ctx context.Context, playerId int64, otherPlayersId []int64) error
	NotifyBusy(ctx context.Context, playerId int64) error
	NotifyFree(ctx context.Context, playerId int64) error
	GetPlayerStatus(ctx context.Context, playerId int64) (models.PlayerStatus, error)
	GetSubscriptionsStatus(ctx context.Context, playerId int64) (map[int64]models.PlayerStatus, error)
	CreateChallenge(ctx context.Context, from *models.PlayerUserData, to int64) (string, error)
	AcceptChallenge(ctx context.Context, messageId string, from, to int64) (*models.PlayerUserData, error)
	SendAcceptMessage(ctx context.Context, playerId int64, address string) error
	CleanupAfterFailedAccept(ctx context.Context, senderId, receiverId int64) error
	DeclineChallenge(ctx context.Context, messageId string, from, to int64) error
	CancelChallenge(ctx context.Context, messageId string, from, to int64) error
	ReadMessages(ctx context.Context, playerId int64) ([]*models.Message, error)
	CommitMessages(ctx context.Context, playerId int64, messages []*models.Message) error
}

type userKvRepoImpl struct {
	rdb    *redis.Client
	logger *slog.Logger
}

func MakeUserKvRepo(rdb *redis.Client, logger *slog.Logger) UserKvRepo {
	return &userKvRepoImpl{rdb: rdb, logger: logger}
}

func (r *userKvRepoImpl) Register(ctx context.Context, player *models.PlayerUserData) (bool, error) {
	keys := PlayerKeySet{player.Id}
	var registered bool
	err := r.rdb.Watch(ctx, func(tx *redis.Tx) error {
		n, err := tx.Exists(ctx, keys.name()).Result()
		if err != nil {
			return err
		}
		if n == 1 {
			registered = false
			return nil
		}
		_, err = tx.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
			pipe.MSet(ctx, []interface{}{
				keys.name(), player.Name,
				keys.rating(), player.Rating,
				keys.status(), int64(models.Online),
			})
			return nil
		})
		if err != nil {
			return err
		}
		registered = true
		return nil
	}, keys.name())
	if err != nil {
		r.logger.Error("Register: failed", "player_id", player.Id, "error", err)
	} else {
		r.logger.Debug("Register: ok", "player_id", player.Id, "registered", registered)
	}
	return registered, err
}

func (r *userKvRepoImpl) Subscribe(ctx context.Context, playerId int64, otherPlayersId []int64) error {
	keys := PlayerKeySet{playerId}
	members := make([]interface{}, len(otherPlayersId))
	for i, id := range otherPlayersId {
		members[i] = id
	}
	_, err := r.rdb.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.Del(ctx, keys.subscriptions())
		if len(members) > 0 {
			pipe.SAdd(ctx, keys.subscriptions(), members...)
		}
		return nil
	})
	if err != nil {
		r.logger.Error("Subscribe: failed", "player_id", playerId, "error", err)
	} else {
		r.logger.Debug("Subscribe: ok", "player_id", playerId, "count", len(otherPlayersId))
	}
	return err
}

func (r *userKvRepoImpl) GetPlayerStatus(ctx context.Context, playerId int64) (models.PlayerStatus, error) {
	keys := PlayerKeySet{playerId}
	val, err := r.rdb.Get(ctx, keys.status()).Int64()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return models.Offline, nil
		}
		return 0, err
	}
	return models.PlayerStatus(val), nil
}

func (r *userKvRepoImpl) GetSubscriptionsStatus(ctx context.Context, playerId int64) (map[int64]models.PlayerStatus, error) {
	keys := PlayerKeySet{playerId}

	ids, err := r.rdb.SMembers(ctx, keys.subscriptions()).Result()
	if err != nil {
		return nil, err
	}
	if len(ids) == 0 {
		return map[int64]models.PlayerStatus{}, nil
	}

	statusKeys := make([]string, len(ids))
	parsedIds := make([]int64, len(ids))
	for i, idStr := range ids {
		var id int64
		fmt.Sscanf(idStr, "%d", &id)
		parsedIds[i] = id
		statusKeys[i] = PlayerKeySet{id}.status()
	}

	vals, err := r.rdb.MGet(ctx, statusKeys...).Result()
	if err != nil {
		return nil, err
	}

	result := make(map[int64]models.PlayerStatus, len(parsedIds))
	for i, val := range vals {
		if val == nil {
			result[parsedIds[i]] = models.Offline
			continue
		}
		var status int64
		fmt.Sscanf(val.(string), "%d", &status)
		result[parsedIds[i]] = models.PlayerStatus(status)
	}
	return result, nil
}

func (r *userKvRepoImpl) NotifyBusy(ctx context.Context, playerId int64) error {
	keys := PlayerKeySet{playerId}
	_, err := r.rdb.Set(ctx, keys.status(), int64(models.Busy), 0).Result()
	if err != nil {
		r.logger.Error("NotifyBusy: failed", "player_id", playerId, "error", err)
	} else {
		r.logger.Debug("NotifyBusy: ok", "player_id", playerId)
	}
	return err
}

func (r *userKvRepoImpl) NotifyFree(ctx context.Context, playerId int64) error {
	keys := PlayerKeySet{playerId}
	_, err := r.rdb.Set(ctx, keys.status(), int64(models.Online), 0).Result()
	if err != nil {
		r.logger.Error("NotifyFree: failed", "player_id", playerId, "error", err)
	} else {
		r.logger.Debug("NotifyFree: ok", "player_id", playerId)
	}
	return err
}

func (r *userKvRepoImpl) CreateChallenge(ctx context.Context, from *models.PlayerUserData, to int64) (string, error) {
	toKeys := PlayerKeySet{to}
	fromKeys := PlayerKeySet{from.Id}
	var msgId string
	err := r.rdb.Watch(ctx, func(tx *redis.Tx) error {
		statusVal, err := tx.Get(ctx, fromKeys.status()).Int64()
		if err != nil {
			if errors.Is(err, redis.Nil) {
				return errors.New("sender is not registered")
			}
			return err
		}
		if models.PlayerStatus(statusVal) != models.Online {
			return errors.New("sender is not online")
		}
		existing, err := tx.Exists(ctx, fromKeys.currentChallenge()).Result()
		if err != nil {
			return err
		}
		if existing == 1 {
			return errors.New("sender already has an outgoing challenge")
		}
		statusVal, err = tx.Get(ctx, toKeys.status()).Int64()
		if err != nil {
			if errors.Is(err, redis.Nil) {
				return errors.New("target player is not registered")
			}
			return err
		}
		if models.PlayerStatus(statusVal) != models.Online {
			return errors.New("target player is not online")
		}
		msgId = uuid.New().String()
		msgKeys := messageKeySet{msgId}
		_, err = tx.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
			pipe.MSet(ctx, []interface{}{
				msgKeys.type_(), int64(models.Challenge),
				msgKeys.senderName(), from.Name,
				msgKeys.senderId(), from.Id,
				msgKeys.senderRating(), from.Rating,
				msgKeys.receiverId(), to,
			})
			pipe.RPush(ctx, toKeys.mailbox(), msgId)
			pipe.Set(ctx, fromKeys.currentChallenge(), msgId, 0)
			return nil
		})
		return err
	}, toKeys.status(), fromKeys.status(), fromKeys.currentChallenge())
	if err != nil {
		r.logger.Error("CreateChallenge: failed", "from", from.Id, "to", to, "error", err)
	} else {
		r.logger.Debug("CreateChallenge: ok", "from", from.Id, "to", to, "msg_id", msgId)
	}
	return msgId, err
}

func (r *userKvRepoImpl) AcceptChallenge(ctx context.Context, messageId string, from, to int64) (*models.PlayerUserData, error) {
	fromKeys := PlayerKeySet{from}
	toKeys := PlayerKeySet{to}
	var result *models.PlayerUserData
	err := r.rdb.Watch(ctx, func(tx *redis.Tx) error {
		current, err := tx.Get(ctx, fromKeys.currentChallenge()).Result()
		if err != nil {
			return err
		}
		if current != messageId {
			return errors.New("challenge mismatch")
		}
		statusVal, err := tx.Get(ctx, toKeys.status()).Int64()
		if err != nil {
			if errors.Is(err, redis.Nil) {
				return errors.New("accepting player is not registered")
			}
			return err
		}
		if models.PlayerStatus(statusVal) != models.Online {
			return errors.New("accepting player is not online")
		}
		existing, err := tx.Exists(ctx, toKeys.currentChallenge()).Result()
		if err != nil {
			return err
		}
		if existing == 1 {
			return errors.New("accepting player already has an outgoing challenge")
		}
		msgKeys := messageKeySet{messageId}
		vals, err := tx.MGet(ctx,
			msgKeys.senderName(),
			msgKeys.senderId(),
			msgKeys.senderRating(),
		).Result()
		if err != nil {
			return err
		}
		var name string
		var senderId, rating int64
		if vals[0] != nil {
			name = vals[0].(string)
		}
		if vals[1] != nil {
			fmt.Sscanf(vals[1].(string), "%d", &senderId)
		}
		if vals[2] != nil {
			fmt.Sscanf(vals[2].(string), "%d", &rating)
		}
		_, err = tx.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
			pipe.Del(ctx, fromKeys.currentChallenge())
			pipe.Del(ctx,
				msgKeys.type_(),
				msgKeys.senderName(),
				msgKeys.senderId(),
				msgKeys.senderRating(),
				msgKeys.receiverId(),
			)
			pipe.Set(ctx, fromKeys.status(), int64(models.Busy), 0)
			pipe.Set(ctx, toKeys.status(), int64(models.Busy), 0)
			return nil
		})
		if err != nil {
			return err
		}
		result = &models.PlayerUserData{Id: senderId, Name: name, Rating: rating}
		return nil
	}, fromKeys.currentChallenge(), toKeys.status(), toKeys.currentChallenge())
	if err != nil {
		r.logger.Error("AcceptChallenge: failed", "from", from, "to", to, "msg_id", messageId, "error", err)
	} else {
		r.logger.Debug("AcceptChallenge: ok", "from", from, "to", to, "msg_id", messageId)
	}
	return result, err
}

func (r *userKvRepoImpl) DeclineChallenge(ctx context.Context, messageId string, from, to int64) error {
	fromKeys := PlayerKeySet{from}
	err := r.rdb.Watch(ctx, func(tx *redis.Tx) error {
		current, err := tx.Get(ctx, fromKeys.currentChallenge()).Result()
		if err != nil {
			return err
		}
		if current != messageId {
			return errors.New("challenge mismatch")
		}
		newMsgId := uuid.New().String()
		oldMsgKeys := messageKeySet{messageId}
		newMsgKeys := messageKeySet{newMsgId}
		_, err = tx.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
			pipe.Del(ctx, fromKeys.currentChallenge())
			pipe.Del(ctx,
				oldMsgKeys.type_(),
				oldMsgKeys.senderName(),
				oldMsgKeys.senderId(),
				oldMsgKeys.senderRating(),
				oldMsgKeys.receiverId(),
			)
			pipe.MSet(ctx, []interface{}{
				newMsgKeys.type_(), int64(models.ChallengeDeclined),
				newMsgKeys.originalMessageId(), messageId,
			})
			pipe.RPush(ctx, fromKeys.mailbox(), newMsgId)
			return nil
		})
		return err
	}, fromKeys.currentChallenge())
	if err != nil {
		r.logger.Error("DeclineChallenge: failed", "from", from, "to", to, "msg_id", messageId, "error", err)
	} else {
		r.logger.Debug("DeclineChallenge: ok", "from", from, "to", to, "msg_id", messageId)
	}
	return err
}

func (r *userKvRepoImpl) CancelChallenge(ctx context.Context, messageId string, from, to int64) error {
	fromKeys := PlayerKeySet{from}
	toKeys := PlayerKeySet{to}
	err := r.rdb.Watch(ctx, func(tx *redis.Tx) error {
		current, err := tx.Get(ctx, fromKeys.currentChallenge()).Result()
		if err != nil {
			return err
		}
		if current != messageId {
			return errors.New("challenge mismatch")
		}
		msgKeys := messageKeySet{messageId}
		receiverIdStr, err := tx.Get(ctx, msgKeys.receiverId()).Result()
		if err != nil {
			return err
		}
		var receiverId int64
		fmt.Sscanf(receiverIdStr, "%d", &receiverId)
		if receiverId != to {
			return errors.New("receiver mismatch")
		}
		newMsgId := uuid.New().String()
		newMsgKeys := messageKeySet{newMsgId}
		_, err = tx.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
			pipe.Del(ctx, fromKeys.currentChallenge())
			pipe.Del(ctx,
				msgKeys.type_(),
				msgKeys.senderName(),
				msgKeys.senderId(),
				msgKeys.senderRating(),
				msgKeys.receiverId(),
			)
			pipe.MSet(ctx, []interface{}{
				newMsgKeys.type_(), int64(models.ChallengeCancelled),
				newMsgKeys.senderId(), from,
			})
			pipe.RPush(ctx, toKeys.mailbox(), newMsgId)
			return nil
		})
		return err
	}, fromKeys.currentChallenge())
	if err != nil {
		r.logger.Error("CancelChallenge: failed", "from", from, "to", to, "msg_id", messageId, "error", err)
	} else {
		r.logger.Debug("CancelChallenge: ok", "from", from, "to", to, "msg_id", messageId)
	}
	return err
}

func (r *userKvRepoImpl) Unregister(ctx context.Context, playerId int64) error {
	keys := PlayerKeySet{playerId}
	err := r.rdb.Del(ctx, keys.keys()...).Err()
	if err != nil {
		r.logger.Error("Unregister: failed", "player_id", playerId, "error", err)
	} else {
		r.logger.Debug("Unregister: ok", "player_id", playerId)
	}
	return err
}

func (r *userKvRepoImpl) SendAcceptMessage(ctx context.Context, playerId int64, address string) error {
	playerKeys := PlayerKeySet{playerId}
	err := r.rdb.Watch(ctx, func(tx *redis.Tx) error {
		exists, err := tx.Exists(ctx, playerKeys.name()).Result()
		if err != nil {
			return err
		}
		if exists == 0 {
			return errors.New("player mailbox does not exist")
		}
		msgId := uuid.New().String()
		msgKeys := messageKeySet{msgId}
		_, err = tx.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
			pipe.MSet(ctx, []interface{}{
				msgKeys.type_(), int64(models.ChallengeAccepted),
				msgKeys.address(), address,
			})
			pipe.RPush(ctx, playerKeys.mailbox(), msgId)
			return nil
		})
		return err
	}, playerKeys.name())
	if err != nil {
		r.logger.Error("SendAcceptMessage: failed", "player_id", playerId, "error", err)
	} else {
		r.logger.Debug("SendAcceptMessage: ok", "player_id", playerId)
	}
	return err
}

func (r *userKvRepoImpl) ReadMessages(ctx context.Context, playerId int64) ([]*models.Message, error) {
	playerKeys := PlayerKeySet{playerId}

	ids, err := r.rdb.LRange(ctx, playerKeys.mailbox(), 0, -1).Result()
	if err != nil {
		return nil, err
	}
	if len(ids) == 0 {
		return []*models.Message{}, nil
	}

	allKeys := make([]string, 0, len(ids)*7)
	for _, id := range ids {
		mk := messageKeySet{id}
		allKeys = append(allKeys,
			mk.type_(),
			mk.senderName(),
			mk.senderId(),
			mk.senderRating(),
			mk.receiverId(),
			mk.address(),
			mk.originalMessageId(),
		)
	}

	vals, err := r.rdb.MGet(ctx, allKeys...).Result()
	if err != nil {
		return nil, err
	}

	messages := make([]*models.Message, 0, len(ids))
	for i, id := range ids {
		base := i * 7
		if vals[base] == nil {
			continue
		}
		msg := &models.Message{Id: id}
		fmt.Sscanf(vals[base].(string), "%d", (*int64)(&msg.Type))
		if v := vals[base+1]; v != nil {
			msg.SenderName = v.(string)
		}
		if v := vals[base+2]; v != nil {
			fmt.Sscanf(v.(string), "%d", &msg.SenderId)
		}
		if v := vals[base+3]; v != nil {
			fmt.Sscanf(v.(string), "%d", &msg.SenderRating)
		}
		if v := vals[base+4]; v != nil {
			fmt.Sscanf(v.(string), "%d", &msg.ReceiverId)
		}
		if v := vals[base+5]; v != nil {
			msg.Address = v.(string)
		}
		if v := vals[base+6]; v != nil {
			msg.OriginalMessageId = v.(string)
		}
		messages = append(messages, msg)
	}
	return messages, nil
}

func (r *userKvRepoImpl) CommitMessages(ctx context.Context, playerId int64, messages []*models.Message) error {
	if len(messages) == 0 {
		return nil
	}
	playerKeys := PlayerKeySet{playerId}
	_, err := r.rdb.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.LPopCount(ctx, playerKeys.mailbox(), len(messages))
		for _, msg := range messages {
			if msg.Type == models.Challenge {
				continue
			}
			mk := messageKeySet{msg.Id}
			pipe.Del(ctx,
				mk.type_(),
				mk.senderName(),
				mk.senderId(),
				mk.senderRating(),
				mk.receiverId(),
				mk.address(),
				mk.originalMessageId(),
			)
		}
		return nil
	})
	if err != nil {
		r.logger.Error("CommitMessages: failed", "player_id", playerId, "count", len(messages), "error", err)
	} else {
		r.logger.Debug("CommitMessages: ok", "player_id", playerId, "count", len(messages))
	}
	return err
}

func (r *userKvRepoImpl) CleanupAfterFailedAccept(ctx context.Context, senderId, receiverId int64) error {
	senderKeys := PlayerKeySet{senderId}
	newMsgId := uuid.New().String()
	newMsgKeys := messageKeySet{newMsgId}
	_, err := r.rdb.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.Del(ctx, senderKeys.currentChallenge())
		pipe.MSet(ctx, []interface{}{
			newMsgKeys.type_(), int64(models.ChallengeDeclined),
		})
		pipe.RPush(ctx, senderKeys.mailbox(), newMsgId)
		pipe.Set(ctx, senderKeys.status(), int64(models.Online), 0)
		pipe.Set(ctx, PlayerKeySet{receiverId}.status(), int64(models.Online), 0)
		return nil
	})
	if err != nil {
		r.logger.Error("CleanupAfterFailedAccept: failed", "sender_id", senderId, "receiver_id", receiverId, "error", err)
	} else {
		r.logger.Debug("CleanupAfterFailedAccept: ok", "sender_id", senderId, "receiver_id", receiverId)
	}
	return err
}

package client

import (
	"context"
	"fmt"
	"log/slog"
	msPb "team_dynamics/api/proto/match_service"
	"team_dynamics/user_event/connection"
	"team_dynamics/user_event/downstream"
	ueJson "team_dynamics/user_event/json"
	"team_dynamics/user_event/models"
	"team_dynamics/user_event/redis"
	"time"
)

type Client interface {
	Close()
	Run(ctx context.Context)
}

type clientImpl struct {
	conn                      *connection.Connection
	disconnectCh              chan<- Client
	closeCh                   chan struct{}
	checkStatusTicker         *time.Ticker
	userKvRepo                redis.UserKvRepo
	logger                    *slog.Logger
	player                    *models.PlayerUserData
	msFactory                 downstream.MatchServiceClientFactory
	config                    *UserEventConfig
	currentChallengeMessageId string
	currentChallengeTargetId  int64
	lastKnownStatuses         map[int64]models.PlayerStatus
}

func (c *clientImpl) Close() {
	select {
	case c.closeCh <- struct{}{}:
	default:
	}
}

func (c *clientImpl) terminate(ctx context.Context) {
	if err := c.conn.Close(); err != nil {
		c.logger.Warn("failed to close connection", "error", err)
	} else {
		c.logger.Debug("connection closed")
	}
	if err := c.userKvRepo.Unregister(context.Background(), c.player.Id); err != nil {
		c.logger.Error("failed to remove connection", "player_id", c.player.Id, "error", err)
	} else {
		c.logger.Debug("connection removed", "player_id", c.player.Id)
	}
	c.disconnectCh <- c
}

func (c *clientImpl) onPanic(rec any, out *[]*ueJson.Event, cleanup *[]func() error) bool {
	c.logger.Error("Panic recovered", "panic", rec)
	event, err := ueJson.MakeEvent(ueJson.ErrorPayload{Message: fmt.Sprintf("Internal error: %v", rec)})
	if err != nil {
		*out = append(*out, event)
	}
	return true
}

func (c *clientImpl) onSubscribe(ctx context.Context, payload ueJson.SubscribePayload, out *[]*ueJson.Event, cleanup *[]func() error) bool {
	c.logger.Debug("Subscribe called")
	if err := c.userKvRepo.Subscribe(ctx, c.player.Id, payload.Users); err != nil {
		c.logger.Error("failed to subscribe", "error", err)
		if event, err := ueJson.MakeEvent(ueJson.ErrorPayload{Message: err.Error()}); err == nil {
			*out = append(*out, event)
		}
		return false
	}
	c.lastKnownStatuses = make(map[int64]models.PlayerStatus)
	return c.checkSubscriptionsStatus(ctx, out, cleanup)
}

func (c *clientImpl) onChallenge(ctx context.Context, payload ueJson.ChallengePayload, out *[]*ueJson.Event, cleanup *[]func() error) bool {
	msgId, err := c.userKvRepo.CreateChallenge(ctx, c.player, payload.UserId)
	if err != nil {
		c.logger.Error("failed to create challenge", "error", err)
		if event, err := ueJson.MakeEvent(ueJson.ErrorPayload{Message: err.Error()}); err == nil {
			*out = append(*out, event)
		}
		return false
	}
	c.currentChallengeMessageId = msgId
	c.currentChallengeTargetId = payload.UserId
	return false
}

func (c *clientImpl) onCancelChallenge(ctx context.Context, payload ueJson.CancelChallengePayload, out *[]*ueJson.Event, cleanup *[]func() error) bool {
	if err := c.userKvRepo.CancelChallenge(ctx, c.currentChallengeMessageId, c.player.Id, c.currentChallengeTargetId); err != nil {
		c.logger.Error("failed to cancel challenge", "error", err)
		if event, err := ueJson.MakeEvent(ueJson.ErrorPayload{Message: err.Error()}); err == nil {
			*out = append(*out, event)
		}
		return false
	}
	c.currentChallengeMessageId = ""
	c.currentChallengeTargetId = 0
	return false
}

func (c *clientImpl) AcceptChallenge(ctx context.Context, payload ueJson.AcceptChallengePayload, out *[]*ueJson.Event, cleanup *[]func() error) bool {
	sender, err := c.userKvRepo.AcceptChallenge(ctx, payload.MessageId, payload.UserId, c.player.Id)
	if err != nil {
		c.logger.Error("AcceptChallenge: repo failed", "error", err)
		if event, err := ueJson.MakeEvent(ueJson.ErrorPayload{Message: err.Error()}); err == nil {
			*out = append(*out, event)
		}
		return false
	}

	regId := payload.MessageId
	startResp, err := c.msFactory.StartMatch(ctx, &msPb.StartMatchRequest{
		Matches: []*msPb.InputMatch{
			{
				Player1: &msPb.PlayerData{
					PlayerId:     &sender.Id,
					PlayerName:   &sender.Name,
					PlayerRating: &sender.Rating,
					RegId:        &regId,
				},
				Player2: &msPb.PlayerData{
					PlayerId:     &c.player.Id,
					PlayerName:   &c.player.Name,
					PlayerRating: &c.player.Rating,
					RegId:        &regId,
				},
				Fleet: nil,
			},
		},
	})

	cleanupRepo := func() {
		if err := c.userKvRepo.CleanupAfterFailedAccept(ctx, sender.Id, c.player.Id); err != nil {
			c.logger.Error("AcceptChallenge: cleanup failed", "error", err)
		}
	}

	if err != nil || len(startResp.Results) == 0 {
		c.logger.Error("AcceptChallenge: StartMatch failed", "error", err)
		cleanupRepo()
		if event, err := ueJson.MakeEvent(ueJson.ErrorPayload{Message: "failed to start match"}); err == nil {
			*out = append(*out, event)
		}
		return false
	}

	result := startResp.Results[0]
	p1Reenter := result.Player1FailResponse == msPb.PlayerFailResponse_PLAYER_FAIL_RESPONSE_REENTER
	p2Reenter := result.Player2FailResponse == msPb.PlayerFailResponse_PLAYER_FAIL_RESPONSE_REENTER

	if p1Reenter || p2Reenter {
		c.logger.Error("AcceptChallenge: StartMatch returned reenter")
		cleanupRepo()
		if event, err := ueJson.MakeEvent(ueJson.ErrorPayload{Message: "match start failed, please retry"}); err == nil {
			*out = append(*out, event)
		}
		return false
	}

	getResp, err := c.msFactory.GetMatch(ctx, &msPb.GetMatchRequest{PlayerId: &sender.Id})
	if err != nil || getResp == nil || getResp.ConnectionInfo == nil {
		c.logger.Error("AcceptChallenge: GetMatch failed", "error", err)
		cleanupRepo()
		if event, err := ueJson.MakeEvent(ueJson.ErrorPayload{Message: "failed to get match address"}); err == nil {
			*out = append(*out, event)
		}
		return false
	}

	p1Id := int64(0)
	p2Id := int64(0)
	if getResp.Player1Id != nil {
		p1Id = *getResp.Player1Id
	}
	if getResp.Player2Id != nil {
		p2Id = *getResp.Player2Id
	}
	validPair := (p1Id == sender.Id && p2Id == c.player.Id) || (p1Id == c.player.Id && p2Id == sender.Id)
	if !validPair {
		c.logger.Error("AcceptChallenge: player ids mismatch", "p1", p1Id, "p2", p2Id)
		cleanupRepo()
		if event, err := ueJson.MakeEvent(ueJson.ErrorPayload{Message: "match player mismatch"}); err == nil {
			*out = append(*out, event)
		}
		return false
	}

	address := *getResp.ConnectionInfo.Address
	if err := c.userKvRepo.SendAcceptMessage(ctx, sender.Id, address); err != nil {
		c.logger.Error("Cannot send accept: %v", "error", err)
	}
	if event, err := ueJson.MakeEvent(ueJson.MatchStartedPayload{Address: address}); err == nil {
		*out = append(*out, event)
	}
	return false
}

func (c *clientImpl) DeclineChallenge(ctx context.Context, payload ueJson.DeclineChallengePayload, out *[]*ueJson.Event, cleanup *[]func() error) bool {
	if err := c.userKvRepo.DeclineChallenge(ctx, payload.MessageId, payload.UserId, c.player.Id); err != nil {
		c.logger.Error("failed to decline challenge", "error", err)
		if event, err := ueJson.MakeEvent(ueJson.ErrorPayload{Message: err.Error()}); err == nil {
			*out = append(*out, event)
		}
		return false
	}
	return false
}

func (c *clientImpl) onNotifyBusy(ctx context.Context, out *[]*ueJson.Event, cleanup *[]func() error) bool {
	if err := c.userKvRepo.NotifyBusy(ctx, c.player.Id); err != nil {
		c.logger.Error("failed to notify busy", "error", err)
		if event, err := ueJson.MakeEvent(ueJson.ErrorPayload{Message: err.Error()}); err == nil {
			*out = append(*out, event)
		}
	}
	return false
}

func (c *clientImpl) onNotifyFree(ctx context.Context, out *[]*ueJson.Event, cleanup *[]func() error) bool {
	if err := c.userKvRepo.NotifyFree(ctx, c.player.Id); err != nil {
		c.logger.Error("failed to notify free", "error", err)
		if event, err := ueJson.MakeEvent(ueJson.ErrorPayload{Message: err.Error()}); err == nil {
			*out = append(*out, event)
		}
	}
	return false
}

func (c *clientImpl) onMessage(ctx context.Context, req *connection.ReceivedMessage, out *[]*ueJson.Event, cleanup *[]func() error) bool {
	if req.Err != nil {
		c.logger.Info("Error while reading", "error", req.Err)
		return true
	}
	switch req.Request.Type {
	case ueJson.Subscribe:
		return c.onSubscribe(ctx, req.Request.Payload.(ueJson.SubscribePayload), out, cleanup)
	case ueJson.Challenge:
		return c.onChallenge(ctx, req.Request.Payload.(ueJson.ChallengePayload), out, cleanup)
	case ueJson.CancelChallenge:
		return c.onCancelChallenge(ctx, req.Request.Payload.(ueJson.CancelChallengePayload), out, cleanup)
	case ueJson.NotifyBusy:
		return c.onNotifyBusy(ctx, out, cleanup)
	case ueJson.NotifyFree:
		return c.onNotifyFree(ctx, out, cleanup)
	case ueJson.AcceptChallenge:
		return c.AcceptChallenge(ctx, req.Request.Payload.(ueJson.AcceptChallengePayload), out, cleanup)
	case ueJson.DeclineChallenge:
		return c.DeclineChallenge(ctx, req.Request.Payload.(ueJson.DeclineChallengePayload), out, cleanup)
	}
	c.logger.Error("Unknown request type", "request type", req.Request.Type)
	return false
}

func (c *clientImpl) checkSubscriptionsStatus(ctx context.Context, out *[]*ueJson.Event, cleanup *[]func() error) bool {
	statuses, err := c.userKvRepo.GetSubscriptionsStatus(ctx, c.player.Id)
	if err != nil {
		c.logger.Error("checkSubscriptionsStatus: failed", "error", err)
		return false
	}
	for userId, status := range statuses {
		if prev, ok := c.lastKnownStatuses[userId]; ok && prev == status {
			continue
		}
		c.lastKnownStatuses[userId] = status
		var jsonStatus ueJson.UserStatus
		c.logger.Debug("Status of user changed", "userId", userId, "receiverId", c.player.Id)
		switch status {
		case models.Online:
			jsonStatus = ueJson.Online
		case models.Busy:
			jsonStatus = ueJson.InGame
		default:
			jsonStatus = ueJson.Offline
		}
		if event, err := ueJson.MakeEvent(ueJson.StatusChangedPayload{UserId: userId, NewStatus: jsonStatus}); err == nil {
			*out = append(*out, event)
		}
	}
	return false
}

func (c *clientImpl) checkMailbox(ctx context.Context, out *[]*ueJson.Event, cleanup *[]func() error) bool {
	messages, err := c.userKvRepo.ReadMessages(ctx, c.player.Id)
	if err != nil {
		c.logger.Error("checkMailbox: ReadMessages failed", "error", err)
		return false
	}
	if len(messages) == 0 {
		return false
	}
	for _, msg := range messages {
		var event *ueJson.Event
		switch msg.Type {
		case models.Challenge:
			event, err = ueJson.MakeEvent(ueJson.ChallengeReceivedPayload{
				MessageId: msg.Id,
				UserId:    msg.SenderId,
				UserName:  msg.SenderName,
			})
		case models.ChallengeAccepted:
			event, err = ueJson.MakeEvent(ueJson.ChallengeAcceptedPayload{Address: msg.Address})
		case models.ChallengeDeclined:
			event, err = ueJson.MakeEvent(ueJson.ChallengeDeclinedPayload{})
		case models.ChallengeCancelled:
			event, err = ueJson.MakeEvent(ueJson.ChallengeCancelledPayload{UserId: msg.SenderId})
		default:
			c.logger.Warn("checkMailbox: unknown message type", "type", msg.Type)
			continue
		}
		if err != nil {
			c.logger.Error("checkMailbox: MakeEvent failed", "error", err)
			continue
		}
		*out = append(*out, event)
	}
	*cleanup = append(*cleanup, func() error {
		return c.userKvRepo.CommitMessages(ctx, c.player.Id, messages)
	})
	return false
}

func (c *clientImpl) onCheckStatusTick(ctx context.Context, out *[]*ueJson.Event, cleanup *[]func() error) bool {
	terminate1 := c.checkSubscriptionsStatus(ctx, out, cleanup)
	terminate2 := c.checkMailbox(ctx, out, cleanup)
	return terminate1 || terminate2
}

func (c *clientImpl) onEnd(ctx context.Context, out *[]*ueJson.Event, cleanup *[]func() error) bool {
	return true
}

func (c *clientImpl) handle(ctx context.Context, out *[]*ueJson.Event, cleanup *[]func() error) (res bool) {
	defer func(resOut *bool) {
		if r := recover(); r != nil {
			*resOut = c.onPanic(r, out, cleanup)
		}
	}(&res)
	select {
	case msg := <-c.conn.Messages():
		return c.onMessage(ctx, &msg, out, cleanup)
	case <-c.checkStatusTicker.C:
		return c.onCheckStatusTick(ctx, out, cleanup)
	case <-c.closeCh:
		return c.onEnd(ctx, out, cleanup)
	case <-ctx.Done():
		return c.onEnd(ctx, out, cleanup)
	}
}

func (c *clientImpl) Run(ctx context.Context) {
	defer c.checkStatusTicker.Stop()
	for {
		events := make([]*ueJson.Event, 0)
		cleanup := make([]func() error, 0)
		stop := c.handle(ctx, &events, &cleanup)
		var err error = nil
		if len(events) > 0 {
			err = c.conn.Write(&ueJson.Response{Events: events})
		}
		if err != nil {
			c.logger.Warn("Error while sending events", "error", err)
		} else {
			for _, fn := range cleanup {
				if err := func() (retErr error) {
					defer func() {
						if r := recover(); r != nil {
							retErr = fmt.Errorf("panic in cleanup: %v", r)
						}
					}()
					return fn()
				}(); err != nil {
					c.logger.Error("cleanup error", "error", err)
				}
			}
		}
		if stop {
			c.terminate(ctx)
			break
		}
	}
}

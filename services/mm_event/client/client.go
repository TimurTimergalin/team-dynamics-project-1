package client

import (
	"context"
	"log/slog"
	msPb "team_dynamics/api/proto/match_service"
	"team_dynamics/mm_event/connection"
	mmeJson "team_dynamics/mm_event/json"
	"team_dynamics/mm_event/models"
	"team_dynamics/mm_event/redis"
	"time"
)

type Client interface {
	Close()
	Run()
}

type clientState int32

const (
	Stale clientState = iota + 1
	InQueue
	WaitingForMatch
	Done
	Erroneous
)

type clientImpl struct {
	conn                    *connection.Connection
	disconnectCh            chan<- Client
	ctx                     context.Context
	cancel                  context.CancelFunc
	checkPoolPresenceTicker *time.Ticker
	checkMatchTicker        *time.Ticker
	state                   clientState
	mmPoolRepo              redis.MMPoolRepo
	logger                  *slog.Logger
	player                  *models.Player
	msClient                msPb.MatchServiceClient
	config                  *Config
}

func (c *clientImpl) Close() {
	c.cancel()
}

func (c *clientImpl) terminate() {
	if err := c.conn.Close(); err != nil {
		c.logger.Warn("failed to close connection", "error", err)
	} else {
		c.logger.Debug("connection closed")
	}
	if err := c.mmPoolRepo.RemoveConnection(c.ctx, c.player.Id); err != nil {
		c.logger.Error("failed to remove connection", "player_id", c.player.Id, "error", err)
	} else {
		c.logger.Debug("connection removed", "player_id", c.player.Id)
	}
	c.disconnectCh <- c
}

func (c *clientImpl) removeFromPool() error {
	if c.state != InQueue {
		return nil
	}
	err := c.mmPoolRepo.RemoveFromPool(c.ctx, c.player.Id)
	if err != nil {
		c.logger.Error("failed to remove player from pool", "player_id", c.player.Id, "error", err)
		c.state = Erroneous
		return err
	}
	c.state = Stale
	c.logger.Debug("removed player from pool, state changed to Stale", "player_id", c.player.Id)
	return nil
}

func (c *clientImpl) addToPool() error {
	if c.state != Stale {
		return nil
	}
	err := c.mmPoolRepo.AddToPool(c.ctx, c.player)
	if err != nil {
		c.logger.Error("failed to add player to pool", "player_id", c.player.Id, "error", err)
		c.state = Erroneous
		return err
	}
	c.state = InQueue
	c.logger.Debug("added player to pool, state changed to InQueue", "player_id", c.player.Id)
	return nil
}

func (c *clientImpl) checkInPool() error {
	if c.state != InQueue {
		return nil
	}
	inPool, err := c.mmPoolRepo.CheckInPool(c.ctx, c.player.Id)
	if err != nil {
		c.logger.Error("failed to check player in pool", "player_id", c.player.Id, "error", err)
		c.state = Erroneous
		return err
	}
	if !inPool {
		c.state = WaitingForMatch
		c.logger.Debug("player no longer in pool, state changed to WaitingForMatch", "player_id", c.player.Id)
	}
	return nil
}

func (c *clientImpl) checkMatch() (*string, error) {
	if c.state != WaitingForMatch {
		return nil, nil
	}
	req := &msPb.GetMatchRequest{
		PlayerId: &c.player.Id,
	}
	resp, err := c.msClient.GetMatch(c.ctx, req)
	if err != nil {
		c.logger.Error("failed to get match", "player_id", c.player.Id, "error", err)
		c.state = Erroneous
		return nil, err
	}
	if resp.ConnectionInfo != nil && resp.ConnectionInfo.Address != nil {
		c.state = Done
		c.logger.Debug("match found, state changed to Done", "player_id", c.player.Id, "address", *resp.ConnectionInfo.Address)
		return resp.ConnectionInfo.Address, nil
	}
	return nil, nil
}

func (c *clientImpl) onPanic() bool {
	defer func() {
		if r := recover(); r != nil {
			c.logger.Error("panic in removeFromPool", "recover", r)
		}
	}()
	if err := c.removeFromPool(); err != nil {
		c.logger.Error("failed to remove from pool on finish", "error", err)
	}
	return true
}

func (c *clientImpl) onMessage(msg connection.ReceivedMessage) bool {
	if msg.Err != nil {
		c.logger.Error("read error", "error", msg.Err)
		if err := c.removeFromPool(); err != nil {
			c.logger.Error("failed to remove from pool on read error", "error", err)
		}
		return true
	}
	req := msg.Request
	if req == nil {
		return false
	}
	switch req.Type {
	case mmeJson.Unregister:
		if err := c.removeFromPool(); err != nil {
			errMsg := err.Error()
			resp := &mmeJson.Response{
				Type:         mmeJson.Error,
				ErrorMessage: &errMsg,
			}
			if sendErr := c.conn.Write(resp); sendErr != nil {
				c.logger.Error("failed to send error response", "error", sendErr)
			}
			c.logger.Error("failed to unregister", "error", err)
			return true
		}
		// Success: send Unregistered response
		resp := &mmeJson.Response{
			Type: mmeJson.Unregistered,
		}
		if err := c.conn.Write(resp); err != nil {
			c.logger.Error("failed to send unregistered response", "error", err)
		} else {
			c.logger.Debug("unregistered response sent")
		}
		return false
	case mmeJson.Register:
		if err := c.addToPool(); err != nil {
			errMsg := err.Error()
			resp := &mmeJson.Response{
				Type:         mmeJson.Error,
				ErrorMessage: &errMsg,
			}
			if sendErr := c.conn.Write(resp); sendErr != nil {
				c.logger.Error("failed to send error response", "error", sendErr)
			}
			c.logger.Error("failed to register", "error", err)
			return true
		}
		// Success: send Registered response
		resp := &mmeJson.Response{
			Type: mmeJson.Registered,
		}
		if err := c.conn.Write(resp); err != nil {
			c.logger.Error("failed to send registered response", "error", err)
		} else {
			c.logger.Debug("registered response sent")
		}
		return false
	default:
		c.logger.Warn("unknown request type", "type", req.Type)
		return false
	}
}

func (c *clientImpl) onCancel() bool {
	if err := c.removeFromPool(); err != nil {
		c.logger.Error("failed to remove from pool on cancel", "error", err)
	}
	return true
}

func (c *clientImpl) onCheckInPoolTick() bool {
	if err := c.checkInPool(); err != nil {
		errMsg := err.Error()
		resp := &mmeJson.Response{
			Type:         mmeJson.Error,
			ErrorMessage: &errMsg,
		}
		if sendErr := c.conn.Write(resp); sendErr != nil {
			c.logger.Error("failed to send error response on checkInPool tick", "error", sendErr)
		}
		c.logger.Error("checkInPool tick error", "error", err)
		return true
	}
	return false
}

func (c *clientImpl) onCheckMatchTick() bool {
	address, err := c.checkMatch()
	if err != nil {
		errMsg := err.Error()
		resp := &mmeJson.Response{
			Type:         mmeJson.Error,
			ErrorMessage: &errMsg,
		}
		if sendErr := c.conn.Write(resp); sendErr != nil {
			c.logger.Error("failed to send error response on checkMatch tick", "error", sendErr)
		}
		c.logger.Error("checkMatch tick error", "error", err)
		return true
	}
	if address != nil {
		resp := &mmeJson.Response{
			Type:    mmeJson.Match,
			Address: address,
		}
		if err := c.conn.Write(resp); err != nil {
			c.logger.Error("failed to send match response", "error", err)
		} else {
			c.logger.Debug("match response sent", "address", *address)
		}
		return true
	}
	return false
}

func (c *clientImpl) handle() (shouldTerminate bool) {
	defer func() {
		if r := recover(); r != nil {
			c.logger.Error("panic in handle", "recover", r)
			shouldTerminate = c.onPanic()
		}
	}()
	select {
	case msg := <-c.conn.Messages():
		return c.onMessage(msg)
	case <-c.ctx.Done():
		return c.onCancel()
	case <-c.checkPoolPresenceTicker.C:
		return c.onCheckInPoolTick()
	case <-c.checkMatchTicker.C:
		return c.onCheckMatchTick()
	case <-time.After(c.config.MessageReceivedTimeout):
		return c.onCancel()
	}
}

func (c *clientImpl) Run() {
	defer c.cancel()
	defer c.checkPoolPresenceTicker.Stop()
	defer c.checkMatchTicker.Stop()
	for {
		if c.handle() {
			c.terminate()
			break
		}
	}
}

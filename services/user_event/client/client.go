package client

import (
	"context"
	"log/slog"
	"team_dynamics/user_event/connection"
	"team_dynamics/user_event/downstream"
	"team_dynamics/user_event/models"
	"team_dynamics/user_event/redis"
	"time"
)

type Client interface {
	Close()
	Run()
}

type State int32

const (
	Free State = 1 + iota
	Busy
	Erroneous
)

type clientImpl struct {
	conn              *connection.Connection
	disconnectCh      chan<- Client
	ctx               context.Context
	cancel            context.CancelFunc
	checkStatusTicker *time.Ticker
	state             State
	userKvRepo        redis.UserKvRepo
	logger            *slog.Logger
	player            *models.PlayerUserData
	msFactory         downstream.MatchServiceClientFactory
	config            *UserEventConfig
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
	if err := c.userKvRepo.Unregister(c.ctx, c.player.Id); err != nil {
		c.logger.Error("failed to remove connection", "player_id", c.player.Id, "error", err)
	} else {
		c.logger.Debug("connection removed", "player_id", c.player.Id)
	}
	c.disconnectCh <- c
}

func (c *clientImpl) Run() {
	panic("Not implemented")
}

package client

import (
	"context"
	"log/slog"
	"time"
)

type Hub interface {
	Register(client Client)
	Run(ctx context.Context)
}

type hubImpl struct {
	registerCh   chan Client
	disconnectCh <-chan Client
	clients      map[Client]struct{}
	ctx          context.Context
	logger       *slog.Logger
	config       *UserEventConfig
}

func NewHub(disconnectCh <-chan Client, registerCh chan Client, logger *slog.Logger, config *UserEventConfig) Hub {
	return &hubImpl{
		registerCh:   registerCh,
		disconnectCh: disconnectCh,
		clients:      make(map[Client]struct{}),
		logger:       logger,
		config:       config,
	}
}

func (h *hubImpl) Register(client Client) {
	select {
	case h.registerCh <- client:
		go client.Run(h.ctx)
	case <-time.After(5 * time.Second):
		h.logger.Warn("register channel full, dropping client")
		go client.Run(h.ctx)
		client.Close()
	}
}

func (h *hubImpl) Run(ctx context.Context) {
	h.ctx = ctx
	for {
		select {
		case <-h.ctx.Done():
			for client := range h.clients {
				client.Close()
			}
			h.logger.Info("hub stopped")
			return
		case client := <-h.registerCh:
			h.clients[client] = struct{}{}
			h.logger.Debug("client registered", "client", client)
		case client := <-h.disconnectCh:
			delete(h.clients, client)
			h.logger.Debug("client unregistered", "client", client)
		}
	}
}

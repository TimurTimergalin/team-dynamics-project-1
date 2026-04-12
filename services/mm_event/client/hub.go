package client

import (
	"context"
	"log/slog"
	"time"
)

type Hub interface {
	Register(client Client)
	Close()
	Run()
}

type hubImpl struct {
	registerCh   chan Client
	disconnectCh <-chan Client
	clients      map[Client]struct{}
	ctx          context.Context
	cancel       context.CancelFunc
	logger       *slog.Logger
	config       *Config
}

func NewHub(disconnectCh <-chan Client, registerCh chan Client, logger *slog.Logger, config *Config) Hub {
	ctx, cancel := context.WithCancel(context.Background())
	return &hubImpl{
		registerCh:   registerCh,
		disconnectCh: disconnectCh,
		clients:      make(map[Client]struct{}),
		ctx:          ctx,
		cancel:       cancel,
		logger:       logger,
		config:       config,
	}
}

func (h *hubImpl) Register(client Client) {
	select {
	case h.registerCh <- client:
		go client.Run()
	case <-time.After(h.config.HubRegisterTimeout):
		h.logger.Warn("register channel full, dropping client")
		go client.Run()
		client.Close()
	}
}

func (h *hubImpl) Close() {
	h.cancel()
}

func (h *hubImpl) Run() {
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

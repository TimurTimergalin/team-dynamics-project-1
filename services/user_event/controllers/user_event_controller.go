package controllers

import (
	"context"
	"errors"
	"log"
	"net/http"
	"team_dynamics/user_event/client"
)

type UserEventController struct {
	Hub     client.Hub
	Factory client.Factory
	Server  *http.Server
}

func NewUserEventController(ctx context.Context, hub client.Hub, factory client.Factory, address string) *UserEventController {
	mux := http.NewServeMux()
	mux.HandleFunc("/events", func(w http.ResponseWriter, r *http.Request) {
		cl, err := factory.MakeClient(ctx, w, r)
		if err != nil {
			return
		}
		hub.Register(cl)
	})

	return &UserEventController{
		Hub:     hub,
		Factory: factory,
		Server: &http.Server{
			Addr:    address,
			Handler: mux,
		},
	}
}

func (c *UserEventController) Run(ctx context.Context) {
	go c.Hub.Run(ctx)

	if err := c.Server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("errors while listening: %v", err)
	}
}

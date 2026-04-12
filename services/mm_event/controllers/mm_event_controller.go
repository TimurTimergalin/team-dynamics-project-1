package controllers

import (
	"errors"
	"log"
	"net/http"
	"team_dynamics/mm_event/client"
)

type MMEventController struct {
	Hub     client.Hub
	Factory client.Factory
	Server  *http.Server
}

func NewMMEventController(Hub client.Hub, Factory client.Factory, address string) *MMEventController {
	http.HandleFunc("/events", func(w http.ResponseWriter, r *http.Request) {
		cl, _ := Factory.MakeClient(w, r)
		if cl != nil {
			Hub.Register(cl)
		}
	})

	return &MMEventController{
		Hub,
		Factory,
		&http.Server{
			Addr:    address,
			Handler: nil,
		},
	}
}

func (c *MMEventController) Run() {
	defer c.Hub.Close()
	go c.Hub.Run()

	if err := c.Server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("errors while listening: %v", err)
	}
}

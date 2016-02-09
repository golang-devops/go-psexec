package main

import (
	"github.com/ayufan/golang-kardianos-service"
	"net/http"
)

type app struct {
	logger service.Logger
}

func (a *app) Run(logger service.Logger) {
	a.logger = logger

	h := &handler{logger}
	http.HandleFunc("/", h.handler)

	a.logger.Infof("Now serving on '%s'", *address)
	http.ListenAndServe(*address, nil)
}

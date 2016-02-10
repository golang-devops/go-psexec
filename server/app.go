package main

import (
	"github.com/ayufan/golang-kardianos-service"
	"github.com/labstack/echo"
	"gopkg.in/tylerb/graceful.v1"
	"time"
)

type app struct {
	logger          service.Logger
	gracefulTimeout time.Duration
	srv             *graceful.Server
}

func (a *app) Run(logger service.Logger) {
	a.logger = logger
	a.gracefulTimeout = 30 * time.Second //Because a command could be busy executing

	h := &handler{logger}

	e := echo.New()
	t := &htmlTemplateRenderer{}
	e.SetRenderer(t)

	// Unrestricted group
	e.Get("/webui", h.handleWebUIFunc)

	// Restricted group
	r := e.Group("/auth")
	r.Use(JWTAuth(SigningKey))
	r.Post("/exec", h.handleExecFunc)

	a.logger.Infof("Now serving on '%s'", *address)

	a.srv = &graceful.Server{
		Timeout: a.gracefulTimeout,
		Server:  e.Server(*address),
	}

	a.srv.ListenAndServe()
}

func (a *app) OnStop() {
	a.srv.Stop(a.gracefulTimeout)
	<-a.srv.StopChan()
}

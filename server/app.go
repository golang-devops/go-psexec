package main

import (
	"crypto/rsa"
	"github.com/ayufan/golang-kardianos-service"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"gopkg.in/tylerb/graceful.v1"
	"time"
)

type app struct {
	logger          service.Logger
	gracefulTimeout time.Duration
	privateKey      *rsa.PrivateKey
	srv             *graceful.Server
}

func (a *app) Run(logger service.Logger) {
	a.logger = logger
	a.gracefulTimeout = 30 * time.Second //Because a command could be busy executing
	a.privateKey = readPemKey()

	h := &handler{logger, a.privateKey}

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	t := &htmlTemplateRenderer{}
	e.SetRenderer(t)

	// Unrestricted group
	e.Post("/token", h.handleGenerateTokenFunc)
	e.Get("/webui", h.handleWebUIFunc)

	// Restricted group
	r := e.Group("/auth")
	r.Use(GetClientPubkey())
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

package main

import (
	"crypto/rsa"
	"github.com/ayufan/golang-kardianos-service"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"gopkg.in/tylerb/graceful.v1"
	"time"

	"github.com/golang-devops/go-psexec/shared"
)

type app struct {
	logger          service.Logger
	gracefulTimeout time.Duration
	privateKey      *rsa.PrivateKey
	srv             *graceful.Server

	AllowedPublicKeys []*shared.AllowedPublicKey
}

func (a *app) Run(logger service.Logger) {
	a.logger = logger
	a.gracefulTimeout = 30 * time.Second //Because a command could be busy executing
	a.privateKey = readPemKey()

	if len(*allowedPublicKeysFileFlag) > 0 {
		a.AllowedPublicKeys = shared.LoadAllowedPublicKeysFile(*allowedPublicKeysFileFlag)
		if len(a.AllowedPublicKeys) == 0 {
			logger.Errorf("Allowed public key file '%s' was read but contains no keys", *allowedPublicKeysFileFlag)
		}
		for _, allowedKey := range a.AllowedPublicKeys {
			logger.Infof("Allowing key name '%s'", allowedKey.Name)
		}
	} else {
		logger.Errorf("No allowed public keys file specified, no keys will be allowed")
	}

	h := &handler{logger, a.privateKey, a.AllowedPublicKeys}

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

	a.logger.Infof("Now serving on '%s'", *addressFlag)

	a.srv = &graceful.Server{
		Timeout: a.gracefulTimeout,
		Server:  e.Server(*addressFlag),
	}

	a.srv.ListenAndServe()
}

func (a *app) OnStop() {
	a.srv.Stop(a.gracefulTimeout)
	<-a.srv.StopChan()
}

package main

import (
	"crypto/rsa"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"time"

	"github.com/ayufan/golang-kardianos-service"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"gopkg.in/fsnotify.v1"
	"gopkg.in/tylerb/graceful.v1"

	"os/user"

	apex "github.com/apex/log"
	"github.com/apex/log/handlers/json"
	"github.com/golang-devops/go-psexec/services/encoding/checksums"
	"github.com/golang-devops/go-psexec/services/filepath_summary"
	"github.com/golang-devops/go-psexec/shared"
	"github.com/labstack/echo/engine/standard"
	"github.com/natefinch/lumberjack"
)

type app struct {
	debugMode    bool
	accessLogger bool

	logger          service.Logger
	gracefulTimeout time.Duration
	privateKey      *rsa.PrivateKey
	srv             *graceful.Server

	h *handler

	watcherPublicKeys *fsnotify.Watcher
}

func (a *app) registerInterruptSignal() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		defer recover()
		for range c { //for sig := range c {
			if a.srv != nil {
				a.logger.Warning("Interrupt signal received, stopping gracefully")

				c2 := make(chan bool, 1)
				go func() {
					defer recover()
					a.srv.Stop(a.gracefulTimeout)
					c2 <- true
				}()

				select {
				case <-c2:
					a.logger.Warning("Graceful shutdown complete")
					time.Sleep(time.Second) //Sleep a second to give log time to write out
					os.Exit(0)
				case <-time.After(a.gracefulTimeout):
					a.logger.Warning("Normal timeout, forcefully exiting")
					time.Sleep(time.Second) //Sleep a second to give log time to write out
					os.Exit(1)
				}
			}
		}
	}()
}

func (a *app) Run(logger service.Logger) {
	var logFilePath string
	usr, err := user.Current()
	if err == nil {
		logFilePath = filepath.Join(usr.HomeDir, ".config/go-psexec/server-hidden.log")
	} else {
		logFilePath = "server-hidden.log"
	}

	rollingFile := &lumberjack.Logger{
		Filename:   logFilePath,
		MaxSize:    20, // megabytes
		MaxBackups: 20,
		MaxAge:     28, //days
	}

	apex.SetLevel(apex.DebugLevel) //Global level
	apex.SetHandler(json.New(rollingFile))

	tmpLogger := &defaultLogger{
		logger,
		apex.WithField("exe", filepath.Base(os.Args[0])),
	}
	a.logger = tmpLogger

	defer func() {
		if r := recover(); r != nil {
			a.logger.Errorf("Panic recovery in service RUN function: %T %+v", r, r)
		}
	}()

	a.logger.Infof("Running server version %s", TempVersion)

	a.gracefulTimeout = 30 * time.Second //Because a command could be busy executing

	pvtKey, err := shared.ReadPemKey(*serverPemFlag)
	if err != nil {
		logger.Errorf("Cannot read server pem file, error: %s. Exiting server.", err.Error())
		return
	}
	a.privateKey = pvtKey

	checksumsSvc := checksums.New()
	handlerServices := &HandlerServices{
		FilePathSummaries: filepath_summary.New(checksumsSvc),
	}

	a.h = &handler{logger, a.privateKey, handlerServices}

	allowedKeys, err := shared.LoadAllowedPublicKeysFile(*allowedPublicKeysFileFlag)
	if err != nil {
		logger.Errorf("Cannot read allowed public keys, error: %s. Exiting server.", err.Error())
		return
	}
	if len(allowedKeys) == 0 {
		logger.Errorf("Allowed public key file '%s' was read but contains no keys. Exiting server.", *allowedPublicKeysFileFlag)
		return
	}

	a.setAllowedPublicKeys(allowedKeys)

	watcher, err := shared.StartWatcher(*allowedPublicKeysFileFlag, a)
	if err != nil {

	} else {
		a.watcherPublicKeys = watcher
	}

	e := echo.New()
	if a.debugMode {
		e.Debug()
	}

	e.Use(middleware.Recover())
	if a.accessLogger {
		loggerCfg := middleware.DefaultLoggerConfig
		loggerCfg.Output = tmpLogger
		e.Use(middleware.LoggerWithConfig(loggerCfg))
	}

	t := &htmlTemplateRenderer{}
	e.SetRenderer(t)

	// Unrestricted group
	e.POST("/token", a.h.handleGenerateTokenFunc)
	e.GET("/webui", a.h.handleWebUIFunc)

	// Restricted group
	r := e.Group("/auth")
	r.Use(GetClientPubkey())
	r.GET("/ping", a.h.handlePingFunc)
	r.GET("/version", a.h.handleVersionFunc)
	r.POST("/stream", a.h.handleStreamFunc)
	r.POST("/start", a.h.handleStartFunc)
	r.POST("/upload-tar", a.h.handleUploadTarFunc)
	r.GET("/download-tar", a.h.handleDownloadTarFunc)
	r.POST("/delete", a.h.handleDeleteFunc)
	r.POST("/move", a.h.handleMoveFunc)
	r.POST("/copy", a.h.handleCopyFunc)
	r.GET("/stats", a.h.handleStatsFunc)
	r.GET("/path-summary", a.h.handlePathSummaryFunc)
	r.GET("/get-temp-dir", a.h.handleGetTempDirFunc)
	r.GET("/get-os-type", a.h.handleGetOsTypeFunc)

	a.logger.Infof("Now serving on '%s'", *addressFlag)

	server := standard.New(*addressFlag)
	server.SetHandler(e)
	server.SetLogger(e.Logger())
	if e.Debug() {
		e.Logger().Debug("running in debug mode")
	}

	a.srv = &graceful.Server{
		Timeout: a.gracefulTimeout,
		Server:  server.Server,
	}

	a.registerInterruptSignal()

	err = a.srv.ListenAndServe()
	if err != nil {
		if !strings.Contains(err.Error(), "closed network connection") {
			logger.Errorf("Unable to ListenAndServe, error: %s", err.Error())
			time.Sleep(time.Second) //Sleep a second to give log time to write out
		}
	}
}

func (a *app) OnStop() {
	if a.watcherPublicKeys != nil {
		a.watcherPublicKeys.Close()
	}

	a.srv.Stop(a.gracefulTimeout)
	<-a.srv.StopChan()
}

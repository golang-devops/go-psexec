package service

import (
	"github.com/ayufan/golang-kardianos-service"
)

type program struct {
	Logger        service.Logger
	RunHandler    RunHandler
	OnStopHandler OnStopHandler
}

func (p *program) Start(s service.Service) error {
	if service.Interactive() {
		p.Logger.Info("Running in terminal.")
	} else {
		p.Logger.Info("Running under service manager.")
	}

	// Start should not block. Do the actual work async.
	go p.RunHandler.Run(p.Logger)
	return nil
}

func (p *program) Stop(s service.Service) error {
	// Any work in Stop should be quick, usually a few seconds at most.
	if p.OnStopHandler != nil {
		p.OnStopHandler.OnStop()
	}
	p.Logger.Info("Stopped")
	return nil
}

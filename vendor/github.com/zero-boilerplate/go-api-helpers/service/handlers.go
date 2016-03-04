package service

import (
	"github.com/ayufan/golang-kardianos-service"
)

type RunHandler interface {
	Run(logger service.Logger)
}

type OnStopHandler interface {
	OnStop()
}

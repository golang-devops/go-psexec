package main

import (
	"fmt"

	apex "github.com/apex/log"
	"github.com/ayufan/golang-kardianos-service"
)

type defaultLogger struct {
	svcLogger service.Logger
	apexEntry *apex.Entry
}

func (d *defaultLogger) Error(v ...interface{}) error {
	d.apexEntry.Error(fmt.Sprint(v...))
	return d.svcLogger.Error(v...)
}

func (d *defaultLogger) Warning(v ...interface{}) error {
	d.apexEntry.Warn(fmt.Sprint(v...))
	return d.svcLogger.Warning(v...)
}

func (d *defaultLogger) Info(v ...interface{}) error {
	d.apexEntry.Info(fmt.Sprint(v...))
	return d.svcLogger.Info(v...)
}

func (d *defaultLogger) Errorf(format string, a ...interface{}) error {
	d.apexEntry.Errorf(format, a...)
	return d.svcLogger.Errorf(format, a...)
}

func (d *defaultLogger) Warningf(format string, a ...interface{}) error {
	d.apexEntry.Warnf(format, a...)
	return d.svcLogger.Warningf(format, a...)
}

func (d *defaultLogger) Infof(format string, a ...interface{}) error {
	d.apexEntry.Infof(format, a...)
	return d.svcLogger.Infof(format, a...)
}

func (d *defaultLogger) Write(p []byte) (n int, err error) {
	return len(p), d.Infof("%s", string(p))
}

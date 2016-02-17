package shared

import (
	"errors"
	"gopkg.in/fsnotify.v1"
)

type WatcherEventHandler interface {
	OnWatcherEvent(event fsnotify.Event)
	OnWatcherError(err error)
}

func StartWatcher(path string, eventHandler WatcherEventHandler) (*fsnotify.Watcher, error) {
	if eventHandler == nil {
		return nil, errors.New("Event handler cannot be NIL for StartWatcher")
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	go func() {
		for {
			select {
			case event := <-watcher.Events:
				eventHandler.OnWatcherEvent(event)
			case err := <-watcher.Errors:
				eventHandler.OnWatcherError(err)
			}
		}
	}()

	err = watcher.Add(path)
	if err != nil {
		return nil, err
	}

	return watcher, nil
}

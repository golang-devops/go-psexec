package main

import (
	"gopkg.in/fsnotify.v1"
	"path/filepath"
	"strings"

	"github.com/golang-devops/go-psexec/shared"
)

func (a *app) setAllowedPublicKeys(allowedKeys []*shared.AllowedPublicKey) {
	for _, allowedKey := range allowedKeys {
		a.logger.Infof("Allowing key name '%s'", allowedKey.Name)
	}

	a.h.AllowedPublicKeys = allowedKeys
}

func (a *app) OnWatcherEvent(event fsnotify.Event) {
	if !strings.EqualFold(filepath.Base(event.Name), filepath.Base(*allowedPublicKeysFileFlag)) {
		return
	}
	a.logger.Infof("Watcher event: %s", event.String())

	allowedKeys, err := shared.LoadAllowedPublicKeysFile(*allowedPublicKeysFileFlag)
	if err != nil {
		a.logger.Warningf("Cannot read allowed public keys, error: %s. Skipping reload of allowed public keys.", err.Error())
		return
	}
	if len(allowedKeys) == 0 {
		a.logger.Warningf("Allowed public key file '%s' was read but contains no keys. Skipping reload of allowed public keys.", *allowedPublicKeysFileFlag)
		return
	}

	a.setAllowedPublicKeys(allowedKeys)
}

func (a *app) OnWatcherError(err error) {
	if err != nil {
		a.logger.Warningf("Watcher error occurred: %s", err.Error())
	}
}

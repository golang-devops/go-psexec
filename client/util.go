package client

import (
	"strings"
)

func combineServerUrl(baseServerUrl, relUrl string) string {
	return strings.TrimRight(baseServerUrl, "/") + "/" + strings.TrimLeft(relUrl, "/")
}

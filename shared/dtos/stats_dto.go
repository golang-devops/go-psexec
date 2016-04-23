package dtos

import (
	"os"
	"time"
)

type StatsDto struct {
	Path    string
	IsDir   bool
	ModTime time.Time
	Mode    os.FileMode
	Size    int64
}

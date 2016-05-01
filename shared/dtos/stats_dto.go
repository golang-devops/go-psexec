package dtos

import (
	"os"
	"time"
)

type StatsDto struct {
	Path    string
	Exists  bool
	IsDir   bool
	ModTime time.Time
	Mode    os.FileMode
	Size    int64
}

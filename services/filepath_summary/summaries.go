package filepath_summary

import (
	"strings"
	"time"

	"github.com/golang-devops/go-psexec/services/encoding/checksums"
)

//DirSummary holds the summary used for synchronizing across
type DirSummary struct {
	FlattenedFileSummaries []*FileSummary
}

//FileSummary holds the info for a single file
type FileSummary struct {
	FullPath string
	ModTime  time.Time
	Checksum *checksums.ChecksumResult
}

//HaveSamePath compares it to the path of another
func (f *FileSummary) HaveSamePath(other *FileSummary) bool {
	thisCleaned := strings.TrimSpace(f.FullPath)
	otherCleaned := strings.TrimSpace(other.FullPath)
	return strings.EqualFold(thisCleaned, otherCleaned)
}

//NewFileSummary will create a new instance of FileSummary
func NewFileSummary(fullPath string, modTime time.Time, checksum *checksums.ChecksumResult) *FileSummary {
	return &FileSummary{
		FullPath: fullPath,
		ModTime:  modTime,
		Checksum: checksum,
	}
}

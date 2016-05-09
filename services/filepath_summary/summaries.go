package filepath_summary

import (
	"strings"
	"time"

	"github.com/golang-devops/go-psexec/services/encoding/checksums"
)

//DirSummary holds the summary used for synchronizing across
type DirSummary struct {
	BaseDir                string
	FlattenedFileSummaries []*FileSummary
}

//NewDirSummary will create a new instance of DirSummary
func NewDirSummary(baseDir string, flattenedFileSummaries []*FileSummary) *DirSummary {
	return &DirSummary{
		BaseDir:                baseDir,
		FlattenedFileSummaries: flattenedFileSummaries,
	}
}

//FileSummary holds the info for a single file
type FileSummary struct {
	RelativePath string
	ModTime      time.Time
	Checksum     *checksums.ChecksumResult
}

//HaveSamePath compares it to the path of another
func (f *FileSummary) HaveSamePath(other *FileSummary) bool {
	thisCleaned := strings.TrimSpace(f.RelativePath)
	otherCleaned := strings.TrimSpace(other.RelativePath)
	return strings.EqualFold(thisCleaned, otherCleaned)
}

//NewFileSummary will create a new instance of FileSummary
func NewFileSummary(relativePath string, modTime time.Time, checksum *checksums.ChecksumResult) *FileSummary {
	return &FileSummary{
		RelativePath: relativePath,
		ModTime:      modTime,
		Checksum:     checksum,
	}
}

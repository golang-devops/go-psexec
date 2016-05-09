package dtos

import "time"

//FilesystemSummaryDto is the DTO used to transfer file system summaries to compare files before syncing
type FilesystemSummaryDto struct {
	BaseDir        string
	FlattenedFiles []*FileSummary
}

//FileSummary holds the info for a single file
type FileSummary struct {
	RelativePath string
	ModTime      time.Time
	ChecksumHex  string
}

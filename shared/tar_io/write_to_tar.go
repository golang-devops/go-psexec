package tar_io

import (
	"archive/tar"
	"fmt"

	"github.com/golang-devops/go-psexec/shared/io_throttler"
)

func WriteToTar(tarWriter *tar.Writer, file *TarFile) error {
	hdr, err := tar.FileInfoHeader(file.Info, "")
	if err != nil {
		return fmt.Errorf("Unable to get tar FileInfoHeader of tar file '%s', error: %s", file.FileName, err.Error())
	}

	hdr.Name = file.FileName

	if hdr.Xattrs == nil {
		hdr.Xattrs = map[string]string{}
	}
	hdr.Xattrs["SIZE"] = fmt.Sprintf("%d", file.Info.Size())

	if file.IsOnlyFile {
		hdr.Xattrs["SINGLE_FILE_ONLY"] = "1"
	}

	err = tarWriter.WriteHeader(hdr)
	if err != nil {
		return fmt.Errorf("Unable to write tar header for file '%s', error: %s", file.FileName, err.Error())
	}

	if !file.HasContent() {
		return nil
	}

	defer file.ContentReadCloser.Close()
	_, err = io_throttler.CopyThrottled(io_throttler.DefaultIOThrottlingBandwidth, tarWriter, file.ContentReadCloser)
	return err
}

package tar_io

import (
	"archive/tar"
	"fmt"
	"io"

	"github.com/golang-devops/go-psexec/shared"
)

func SaveTarToReceiver(reader io.Reader, tarReceiver TarReceiver) error {
	tarReader := tar.NewReader(reader)

	foundEndOfTar := false
	for {
		hdr, err := tarReader.Next()
		if err != nil {
			if err == io.EOF {
				// end of tar archive
				break
			} else {
				return fmt.Errorf("Unable to read next tar chunk, error: %s", err.Error())
			}
		}

		if hdr.Name == shared.END_OF_TAR_FILENAME {
			foundEndOfTar = true
			continue
		}

		err = tarReceiver.OnEntry(hdr, tarReader)
		if err != nil {
			return fmt.Errorf("Tar reader error: %s", err.Error())
		}
	}

	if !foundEndOfTar {
		return fmt.Errorf("TAR EOF not found. Stream validation failed, something has gone wrong during the transfer.")
	}

	return nil
}

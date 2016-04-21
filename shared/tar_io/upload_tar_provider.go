package tar_io

import (
	"archive/tar"
	"fmt"
	"io"
	"sync"

	"github.com/golang-devops/go-psexec/shared"
)

type UploadHandler interface {
	Read(reader io.Reader) error
	Done() error
}

func UploadProvider(tarProvider TarProvider, handler UploadHandler) error {
	pipeReader, pipeWriter := io.Pipe()
	tarWriter := tar.NewWriter(pipeWriter)
	defer pipeReader.Close() //TODO: Is this necessary?

	wg := &sync.WaitGroup{}
	wg.Add(1)

	var tarProviderError error
	go func() {
		defer wg.Done()
		defer pipeWriter.Close()
		defer tarWriter.Close()

		filesChan := tarProvider.Files()

		for file := range filesChan {
			err := WriteToTar(tarWriter, file)
			if err != nil {
				tarProviderError = err
			}
		}

		hdr := &tar.Header{
			Name: shared.END_OF_TAR_FILENAME,
		}
		err := tarWriter.WriteHeader(hdr)
		if err != nil {
			tarProviderError = err
			return
		}
	}()

	err := handler.Read(pipeReader)
	if err != nil {
		return fmt.Errorf("Error uploading tar stream to pipe reader, error: %s", err.Error())
	}

	//TODO: When an error occurs right above here we do not wait, what happens with that scenario?
	wg.Wait()
	err = handler.Done()
	if err != nil {
		return fmt.Errorf("Done failed for UploadTar, error: %s", err.Error())
	}

	if tarProviderError != nil {
		return fmt.Errorf("Error in reading tar: %s", tarProviderError.Error())
	}

	return nil
}

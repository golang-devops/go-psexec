package tar_io

import (
	"archive/tar"
	"io"
)

type TarReceiver interface {
	OnEntry(tarHeader *tar.Header, tarFileReader io.Reader) error
}

package testing_utils

import (
	"github.com/francoishill/afero"
)

func NewTestingRemoteFileSystem() *RemoteFileSystem {
	// TODO: Make it more obvious that this directory is "populated" with the server which might be another machine
	// but for tests we plan to make the server startup locally on a different port so we can interact with it
	return &RemoteFileSystem{
		Fs: afero.NewOsFs(),
	}
}

type RemoteFileSystem struct {
	afero.Fs
}

func (r *RemoteFileSystem) TempDir() (name string, err error) {
	//TODO: If the 'dir' argument passed into `afero.TempDir` is blank afero implicitly uses `os.TempDir` which I guess complete breaks the abstraction layer - see issue https://github.com/francoishill/afero/issues/84
	return afero.TempDir(r, "", "gopsexec-client-test-")
}

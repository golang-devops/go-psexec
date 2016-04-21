package tar_io

type TarProvider interface {
	IsDir() bool
	RemoteBasePath() string
	Files() <-chan *TarFile
}

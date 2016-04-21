package tar_io

type TarProvider interface {
	Files() <-chan *TarFile
}

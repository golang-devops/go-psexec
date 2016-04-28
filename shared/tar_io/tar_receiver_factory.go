package tar_io

type TarReceiverFactory interface {
	Dir(dir string) TarReceiver
	File(file string) TarReceiver
}

func NewTarReceiverFactory() TarReceiverFactory {
	return &tarReceiverFactory{}
}

type tarReceiverFactory struct{}

func (t *tarReceiverFactory) Dir(dir string) TarReceiver {
	return &dirTarReceiver{
		dir: dir,
	}
}

func (t *tarReceiverFactory) File(file string) TarReceiver {
	return &fileTarReceiver{
		file: file,
	}
}

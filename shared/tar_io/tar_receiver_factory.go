package tar_io

var DefaultTarReceiverFactory = &tarReceiverFactory{}

type TarReceiverFactory interface {
	BasePath(basePath string) TarReceiver
}

func NewTarReceiverFactory() TarReceiverFactory {
	return &tarReceiverFactory{}
}

type tarReceiverFactory struct{}

func (t *tarReceiverFactory) BasePath(basePath string) TarReceiver {
	return &basePathTarReceiver{
		basePath: basePath,
	}
}

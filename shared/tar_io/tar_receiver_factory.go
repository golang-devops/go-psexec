package tar_io

import "io"

type TarReceiverFactory interface {
	Dir(dir string) TarReceiver
	File(file string) TarReceiver
	Writer(writer io.Writer) TarReceiver
}

func NewTarReceiverFactory() TarReceiverFactory {
	return &tarReceiverFactory{}
}

type tarReceiverFactory struct{}

func (t *tarReceiverFactory) Dir(dir string) TarReceiver {
	return &dirTarReceiver{dir: dir}
}

func (t *tarReceiverFactory) File(file string) TarReceiver {
	return &fileTarReceiver{file: file}
}

func (t *tarReceiverFactory) Writer(writer io.Writer) TarReceiver {
	return &writerTarReceiver{writer: writer}
}

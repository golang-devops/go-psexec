package tar_io

var Factories = struct {
	TarReceiver *tarReceiverFactory
	TarProvider *tarProviderFactory
}{
	TarReceiver: &tarReceiverFactory{},
	TarProvider: &tarProviderFactory{},
}

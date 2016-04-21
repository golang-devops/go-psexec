package tar_io

var DefaultTarProviderFactory = &tarProviderFactory{}

type TarProviderFactory interface {
	Dir(fullDirPath, filePattern string) TarProvider
	File(fullFilePath string) TarProvider
}

func NewTarProviderFactory() TarProviderFactory {
	return &tarProviderFactory{}
}

type tarProviderFactory struct{}

func (t *tarProviderFactory) Dir(fullDirPath, filePattern string) TarProvider {
	return &directoryTarProvider{
		fullDirPath: fullDirPath,
		filePattern: filePattern,
	}
}

func (t *tarProviderFactory) File(fullFilePath string) TarProvider {
	return &fileTarProvider{
		fullFilePath: fullFilePath,
	}
}

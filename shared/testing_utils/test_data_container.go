package testing_utils

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/francoishill/afero"
)

func NewTestDataContainer(relativeTestDataDir string, expectedMinFileCount int) (*TestDataContainer, error) {
	fullTestDataDir, err := filepath.Abs(relativeTestDataDir)
	if err != nil {
		return nil, err
	}

	testDataBaseFs := afero.NewBasePathFs(afero.NewOsFs(), fullTestDataDir)
	if err != nil {
		return nil, err
	}

	//These are relative paths because we are using "Base Path FS"
	relativeAllPaths, relativeDirs, relativeFiles, err := listPathsInFsDir(testDataBaseFs, "")
	if err != nil {
		return nil, err
	}
	if len(relativeFiles) < expectedMinFileCount {
		return nil, fmt.Errorf("Expects at least %d files in dir '%s'", expectedMinFileCount, fullTestDataDir)
	}

	return &TestDataContainer{
		FullDir:          fullTestDataDir,
		TestDataBaseFs:   testDataBaseFs,
		RelativeAllPaths: relativeAllPaths,
		RelativeDirs:     relativeDirs,
		RelativeFiles:    relativeFiles,
	}, nil
}

type TestDataContainer struct {
	FullDir          string
	TestDataBaseFs   afero.Fs
	RelativeAllPaths []string
	RelativeDirs     []string
	RelativeFiles    []string
}

func (t *TestDataContainer) ForeachRelativePath(onEach func(f string)) {
	for _, f := range t.RelativeAllPaths {
		onEach(f)
	}
}

func (t *TestDataContainer) ForeachRelativeDir(onEach func(f string)) {
	for _, f := range t.RelativeDirs {
		onEach(f)
	}
}

func (t *TestDataContainer) ForeachRelativeFile(onEach func(f string)) {
	for _, f := range t.RelativeFiles {
		onEach(f)
	}
}

func listPathsInFsDir(fs afero.Fs, dir string) (relativeAllPaths, relativeDirs, relativeFiles []string, returnErr error) {
	relativeAllPaths = []string{}
	relativeDirs = []string{}
	relativeFiles = []string{}

	walkErr := afero.Walk(fs, dir, func(path string, info os.FileInfo, errParam error) error {
		if errParam != nil {
			return errParam
		}
		relativeAllPaths = append(relativeAllPaths, path)
		if info.IsDir() {
			relativeDirs = append(relativeDirs, path)
		} else {
			relativeFiles = append(relativeFiles, path)
		}
		return nil
	})
	if walkErr != nil {
		relativeAllPaths, relativeDirs, relativeFiles = nil, nil, nil
		returnErr = walkErr
		return
	}

	returnErr = nil
	return
}

package monitoring

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"

	"github.com/google/go-jsonnet"
)

// FileImporter imports data from the filesystem.
type EmbedImporter struct {
	FS     embed.FS
	JPaths []string
}

func (importer *EmbedImporter) tryPath(dir, importedPath string) (found bool, contents jsonnet.Contents, foundHere string, err error) {
	var absPath string
	if filepath.IsAbs(importedPath) {
		absPath = importedPath
	} else {
		absPath = filepath.Join(dir, importedPath)
	}

	contentBytes, err := importer.FS.ReadFile(absPath)
	if err != nil && os.IsNotExist(err) {
		return false, jsonnet.Contents{}, "", nil
	} else if err != nil {
		return false, jsonnet.Contents{}, "", err
	}

	return true, jsonnet.MakeContentsRaw(contentBytes), absPath, nil
}

// Import imports file from the filesystem.
func (importer *EmbedImporter) Import(importedFrom, importedPath string) (contents jsonnet.Contents, foundAt string, err error) {
	// TODO(sbarzowski) Make sure that dir is absolute and resolving of ""
	// is independent from current CWD. The default path should be saved
	// in the importer.
	// We need to relativize the paths in the error formatter, so that the stack traces
	// don't have ugly absolute paths (less readable and messy with golden tests).
	dir, _ := filepath.Split(importedFrom)
	found, content, foundHere, err := importer.tryPath(dir, importedPath)
	if err != nil {
		return jsonnet.Contents{}, "", err
	}

	for i := len(importer.JPaths) - 1; !found && i >= 0; i-- {
		found, content, foundHere, err = importer.tryPath(importer.JPaths[i], importedPath)
		if err != nil {
			return jsonnet.Contents{}, "", err
		}
	}

	if !found {
		return jsonnet.Contents{}, "", fmt.Errorf("couldn't open import %#v: no match locally or in the Jsonnet library paths", importedPath)
	}
	return content, foundHere, nil
}

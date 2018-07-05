package xcworkspace

import (
	"encoding/xml"
	"path/filepath"

	"github.com/bitrise-io/go-utils/fileutil"
)

// FileRef ...
type FileRef struct {
	Location string `xml:"location,attr"`
}

// Group ...
type Group struct {
	FileRefs []FileRef `xml:"FileRef"`
}

// Workspace ...
type Workspace struct {
	FileRefs []FileRef `xml:"FileRef"`
	Groups   []Group   `xml:"Group"`
}

// Open ...
func Open(pth string) (Workspace, error) {
	contentsPth := filepath.Join(pth, "contents.xcworkspacedata")
	b, err := fileutil.ReadBytesFromFile(contentsPth)
	if err != nil {
		return Workspace{}, err
	}

	var workspace Workspace
	if err := xml.Unmarshal(b, &workspace); err != nil {
		return Workspace{}, err
	}
	return workspace, nil
}

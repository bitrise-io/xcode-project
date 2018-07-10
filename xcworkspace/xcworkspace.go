package xcworkspace

import (
	"encoding/xml"
	"path/filepath"
	"strings"

	"github.com/bitrise-io/go-utils/fileutil"
	"github.com/bitrise-tools/xcode-project/xcodeproj"
)

// Workspace ...
type Workspace struct {
	FileRefs []FileRef `xml:"FileRef"`
	Groups   []Group   `xml:"Group"`

	Name string
	Path string
}

// FileLocations ...
func (w Workspace) FileLocations() ([]string, error) {
	var fileLocations []string

	for _, fileRef := range w.FileRefs {
		pth, err := fileRef.AbsPath(filepath.Dir(w.Path))
		if err != nil {
			return nil, err
		}

		fileLocations = append(fileLocations, pth)
	}

	for _, group := range w.Groups {
		groupFileLocations, err := group.FileLocations(filepath.Dir(w.Path))
		if err != nil {
			return nil, err
		}

		fileLocations = append(fileLocations, groupFileLocations...)
	}

	return fileLocations, nil
}

// ProjectFileLocations ...
func (w Workspace) ProjectFileLocations() ([]string, error) {
	var projectLocations []string
	fileLocations, err := w.FileLocations()
	if err != nil {
		return nil, err
	}
	for _, fileLocation := range fileLocations {
		if xcodeproj.IsXcodeProj(fileLocation) {
			projectLocations = append(projectLocations, fileLocation)
		}
	}
	return projectLocations, nil
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

	workspace.Name = strings.TrimSuffix(filepath.Base(pth), filepath.Ext(pth))
	workspace.Path = pth

	return workspace, nil
}

// IsWorkspace ...
func IsWorkspace(pth string) bool {
	return filepath.Ext(pth) == ".xcworkspace"
}

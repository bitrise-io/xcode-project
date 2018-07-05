package xcworkspace

import (
	"encoding/xml"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/bitrise-io/go-utils/pathutil"

	"github.com/bitrise-io/go-utils/fileutil"
)

// FileRef ...
type FileRef struct {
	Location string `xml:"location,attr"`
}

// FileRefType ...
type FileRefType string

// Known FileRefTypes
const (
	AbsoluteFileRefType  FileRefType = "absolute"
	GroupFileRefType     FileRefType = "group"
	ContainerFileRefType FileRefType = "container"
	UnknownFileRefType   FileRefType = "unknown"
)

// TypeAndPath ...
func (f FileRef) TypeAndPath() (FileRefType, string, error) {
	s := strings.Split(f.Location, ":")
	if len(s) != 2 {
		return UnknownFileRefType, "", fmt.Errorf("unknown file reference location (%s)", f.Location)
	}

	fileRefType := UnknownFileRefType
	switch s[0] {
	case "absolute":
		fileRefType = AbsoluteFileRefType
	case "group":
		fileRefType = GroupFileRefType
	case "container":
		fileRefType = ContainerFileRefType
	}
	return fileRefType, s[1], nil
}

// AbsPath ...
func (f FileRef) AbsPath(workspaceDir string) (string, error) {
	t, pth, err := f.TypeAndPath()
	if err != nil {
		return "", err
	}

	var absPth string
	switch t {
	case AbsoluteFileRefType:
		absPth = pth
	case GroupFileRefType, ContainerFileRefType:
		absPth = filepath.Join(workspaceDir, pth)
	}

	return pathutil.AbsPath(absPth)
}

// Group ...
type Group struct {
	FileRefs []FileRef `xml:"FileRef"`
}

// Workspace ...
type Workspace struct {
	FileRefs []FileRef `xml:"FileRef"`
	Groups   []Group   `xml:"Group"`

	Name string
	Path string
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

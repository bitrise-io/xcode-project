package xcworkspace

import (
	"encoding/xml"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/bitrise-io/go-utils/fileutil"
	"github.com/bitrise-io/go-utils/pathutil"
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
)

// TypeAndPath ...
func (f FileRef) TypeAndPath() (FileRefType, string, error) {
	s := strings.Split(f.Location, ":")
	if len(s) != 2 {
		return "", "", fmt.Errorf("unknown file reference location (%s)", f.Location)
	}

	switch s[0] {
	case "absolute":
		return AbsoluteFileRefType, s[1], nil
	case "group":
		return GroupFileRefType, s[1], nil
	case "container":
		return ContainerFileRefType, s[1], nil
	default:
		return "", "", fmt.Errorf("")
	}
}

// AbsPath ...
func (f FileRef) AbsPath(dir string) (string, error) {
	t, pth, err := f.TypeAndPath()
	if err != nil {
		return "", err
	}

	var absPth string
	switch t {
	case AbsoluteFileRefType:
		absPth = pth
	case GroupFileRefType, ContainerFileRefType:
		absPth = filepath.Join(dir, pth)
	}

	return pathutil.AbsPath(absPth)
}

// Group ...
type Group struct {
	Location string    `xml:"location,attr"`
	FileRefs []FileRef `xml:"FileRef"`
}

// AbsPath ...
func (g Group) AbsPath(dir string) (string, error) {
	s := strings.Split(g.Location, ":")
	if len(s) != 2 {
		return "", fmt.Errorf("unknown group location (%s)", g.Location)
	}
	pth := filepath.Join(dir, s[1])
	return pathutil.AbsPath(pth)
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

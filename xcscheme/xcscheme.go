package xcscheme

import (
	"encoding/xml"
	"path/filepath"
	"strings"

	"github.com/bitrise-io/go-utils/fileutil"
)

// BuildableReference ...
type BuildableReference struct {
	BlueprintIdentifier string `xml:"BlueprintIdentifier,attr"`
	BlueprintName       string `xml:"BlueprintName,attr"`
	BuildableName       string `xml:"BuildableName,attr"`
	ReferencedContainer string `xml:"ReferencedContainer,attr"`
}

// BuildActionEntry ...
type BuildActionEntry struct {
	BuildForTesting    string `xml:"buildForTesting,attr"`
	BuildForArchiving  string `xml:"buildForArchiving,attr"`
	BuildableReference BuildableReference
}

// BuildAction ...
type BuildAction struct {
	BuildActionEntries []BuildActionEntry `xml:"BuildActionEntries>BuildActionEntry"`
}

// ArchiveAction ...
type ArchiveAction struct {
	BuildConfiguration string `xml:"buildConfiguration,attr"`
}

// Scheme ...
type Scheme struct {
	Name          string
	BuildAction   BuildAction
	ArchiveAction ArchiveAction
}

// Open ...
func Open(pth string) (Scheme, error) {
	b, err := fileutil.ReadBytesFromFile(pth)
	if err != nil {
		return Scheme{}, err
	}

	var scheme Scheme
	if err := xml.Unmarshal(b, &scheme); err != nil {
		return Scheme{}, err
	}
	scheme.Name = strings.TrimSuffix(filepath.Base(pth), filepath.Ext(pth))
	return scheme, nil
}

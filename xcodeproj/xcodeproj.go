package xcodeproj

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/bitrise-io/go-utils/fileutil"
	"github.com/bitrise-tools/xcode-project/serialized"
	"github.com/bitrise-tools/xcode-project/xcscheme"
	"howett.net/plist"
)

// XcodeProj ...
type XcodeProj struct {
	Proj Proj

	Name string
	Path string
}

// Scheme ...
func (p XcodeProj) Scheme(name string) (xcscheme.Scheme, bool) {
	schemes, err := p.Schemes()
	if err != nil {
		return xcscheme.Scheme{}, false
	}

	for _, scheme := range schemes {
		if scheme.Name == name {
			return scheme, true
		}
	}

	return xcscheme.Scheme{}, false
}

// Schemes ...
func (p XcodeProj) Schemes() ([]xcscheme.Scheme, error) {
	pattern := filepath.Join(p.Path, "xcshareddata", "xcschemes", "*.xcscheme")
	pths, err := filepath.Glob(pattern)
	if err != nil {
		return nil, err
	}

	var schemes []xcscheme.Scheme
	for _, pth := range pths {
		scheme, err := xcscheme.Open(pth)
		if err != nil {
			return nil, err
		}
		schemes = append(schemes, scheme)
	}

	return schemes, nil
}

// Open ...
func Open(pth string) (XcodeProj, error) {
	pbxProjPth := filepath.Join(pth, "project.pbxproj")

	b, err := fileutil.ReadBytesFromFile(pbxProjPth)
	if err != nil {
		return XcodeProj{}, err
	}

	var raw serialized.Object
	if _, err := plist.Unmarshal(b, &raw); err != nil {
		return XcodeProj{}, fmt.Errorf("failed to generate json from Pbxproj - error: %s", err)
	}

	objects, err := raw.Object("objects")
	if err != nil {
		return XcodeProj{}, err
	}

	projectID := ""
	for id := range objects {
		object, err := objects.Object(id)
		if err != nil {
			return XcodeProj{}, err
		}

		objectISA, err := object.String("isa")
		if err != nil {
			return XcodeProj{}, err
		}

		if objectISA == "PBXProject" {
			projectID = id
			break
		}
	}

	p, err := parseProj(projectID, objects)
	if err != nil {
		return XcodeProj{}, nil
	}

	return XcodeProj{
		Proj: p,
		Path: pth,
		Name: strings.TrimSuffix(filepath.Base(pth), filepath.Ext(pth)),
	}, nil
}

// IsXcodeProj ...
func IsXcodeProj(pth string) bool {
	return filepath.Ext(pth) == ".xcodeproj"
}

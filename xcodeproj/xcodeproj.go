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
	raw serialized.Object

	nativeTargetByID       map[string]NativeTarget
	targetDependencyByID   map[string]TargetDependency
	configurationListByID  map[string]ConfigurationList
	buildConfigurationByID map[string]BuildConfiguration

	Name string
	Path string

	ISA     string
	ID      string
	Targets []NativeTarget
}

// NativeTarget ...
func (p XcodeProj) NativeTarget(id string) (NativeTarget, error) {
	target, ok := p.nativeTargetByID[id]
	if !ok {
		return NativeTarget{}, fmt.Errorf("native target not found with id: %s", id)
	}
	return target, nil
}

// SharedSchemes ...
func (p XcodeProj) SharedSchemes() ([]xcscheme.Scheme, error) {
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

	proj, err := ParsePBXProj(objects)
	if err != nil {
		return XcodeProj{}, nil
	}
	proj.Path = pth
	proj.Name = strings.TrimSuffix(filepath.Base(pth), filepath.Ext(pth))

	return proj, nil
}

// IsXcodeProj ...
func IsXcodeProj(pth string) bool {
	return strings.HasSuffix(pth, ".xcodeproj")
}

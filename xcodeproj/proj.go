package xcodeproj

import (
	"fmt"

	"github.com/bitrise-tools/xcode-project/serialized"
)

// ParsePBXProj ...
func ParsePBXProj(raw serialized.Object) (XcodeProj, error) {
	proj := XcodeProj{
		raw:                    raw,
		nativeTargetByID:       map[string]NativeTarget{},
		targetDependencyByID:   map[string]TargetDependency{},
		configurationListByID:  map[string]ConfigurationList{},
		buildConfigurationByID: map[string]BuildConfiguration{},
		ISA: "PBXProject",
	}

	rawPBXProj, id, err := proj.firstObject("PBXProject")
	if err != nil {
		return XcodeProj{}, err
	}
	proj.ID = id

	rawTargets, err := rawPBXProj.StringSlice("targets")
	if err != nil {
		return XcodeProj{}, err
	}

	targets := []NativeTarget{}
	for _, targetID := range rawTargets {
		target, err := proj.nativeTarget(targetID)
		if err != nil {
			return XcodeProj{}, err
		}
		targets = append(targets, target)
	}
	proj.Targets = targets

	return proj, nil
}

func (p XcodeProj) firstObject(isa string) (serialized.Object, string, error) {
	for id := range p.raw {
		object, err := p.raw.Object(id)
		if err != nil {
			return serialized.Object{}, "", err
		}

		objectISA, err := object.String("isa")
		if err != nil {
			return nil, "", err
		}

		if objectISA == isa {
			return object, id, nil
		}
	}

	return nil, "", fmt.Errorf("object not found with isa: %s", isa)
}

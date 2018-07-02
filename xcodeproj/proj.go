package xcodeproj

import (
	"github.com/bitrise-tools/xcode-project/serialized"
)

func parseProj(id string, objects serialized.Object) (XcodeProj, error) {
	rawPBXProj, err := objects.Object(id)
	if err != nil {
		return XcodeProj{}, err
	}

	rawTargets, err := rawPBXProj.StringSlice("targets")
	if err != nil {
		return XcodeProj{}, err
	}

	targets := []NativeTarget{}
	for _, targetID := range rawTargets {
		target, err := parseNativeTarget(targetID, objects)
		if err != nil {
			return XcodeProj{}, err
		}
		targets = append(targets, target)
	}

	return XcodeProj{
		ISA:     "PBXProject",
		ID:      id,
		Targets: targets,
	}, nil
}

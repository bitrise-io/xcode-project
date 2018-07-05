package xcodeproj

import (
	"github.com/bitrise-tools/xcode-project/serialized"
)

// Proj ...
type Proj struct {
	ID            string
	NativeTargets []NativeTarget
}

func parseProj(id string, objects serialized.Object) (Proj, error) {
	rawPBXProj, err := objects.Object(id)
	if err != nil {
		return Proj{}, err
	}

	rawTargets, err := rawPBXProj.StringSlice("targets")
	if err != nil {
		return Proj{}, err
	}

	targets := []NativeTarget{}
	for _, targetID := range rawTargets {
		target, err := parseNativeTarget(targetID, objects)
		if err != nil {
			return Proj{}, err
		}
		targets = append(targets, target)
	}

	return Proj{
		ID:            id,
		NativeTargets: targets,
	}, nil
}

// NativeTarget ...
func (p Proj) NativeTarget(id string) (NativeTarget, bool) {
	for _, target := range p.NativeTargets {
		if target.ID == id {
			return target, true
		}
	}
	return NativeTarget{}, false
}

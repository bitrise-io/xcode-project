package xcodeproj

import "github.com/bitrise-tools/xcode-project/serialized"

// TargetDependency ...
type TargetDependency struct {
	ISA    string
	ID     string
	Target NativeTarget
}

func parseTargetDependency(id string, objects serialized.Object) (TargetDependency, error) {
	rawTargetDependency, err := objects.Object(id)
	if err != nil {
		return TargetDependency{}, err
	}

	targetID, err := rawTargetDependency.String("target")
	if err != nil {
		return TargetDependency{}, err
	}

	target, err := parseNativeTarget(targetID, objects)
	if err != nil {
		return TargetDependency{}, err
	}

	return TargetDependency{
		ISA:    "PBXTargetDependency",
		ID:     id,
		Target: target,
	}, nil
}

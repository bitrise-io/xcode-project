package xcodeproj

// TargetDependency ...
type TargetDependency struct {
	ISA    string
	ID     string
	Target NativeTarget
}

func (p XcodeProj) targetDependency(id string) (TargetDependency, error) {
	if targetDependency, ok := p.targetDependencyByID[id]; ok {
		return targetDependency, nil
	}

	rawTargetDependency, err := p.raw.Object(id)
	if err != nil {
		return TargetDependency{}, err
	}

	targetID, err := rawTargetDependency.String("target")
	if err != nil {
		return TargetDependency{}, err
	}

	target, err := p.nativeTarget(targetID)
	if err != nil {
		return TargetDependency{}, err
	}

	targetDependency := TargetDependency{
		ISA:    "PBXTargetDependency",
		ID:     id,
		Target: target,
	}

	p.targetDependencyByID[id] = targetDependency

	return targetDependency, nil
}

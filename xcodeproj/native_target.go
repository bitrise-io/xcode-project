package xcodeproj

// NativeTarget ...
type NativeTarget struct {
	ISA string
	ID  string

	Name                   string
	BuildConfigurationList ConfigurationList
	Dependencies           []TargetDependency
}

func (p XcodeProj) nativeTarget(id string) (NativeTarget, error) {
	if nativeTarget, ok := p.nativeTargetByID[id]; ok {
		return nativeTarget, nil
	}

	rawTarget, err := p.raw.Object(id)
	if err != nil {
		return NativeTarget{}, err
	}

	name, err := rawTarget.String("name")
	if err != nil {
		return NativeTarget{}, err
	}

	buildConfigurationListID, err := rawTarget.String("buildConfigurationList")
	if err != nil {
		return NativeTarget{}, err
	}

	buildConfigurationList, err := p.configurationList(buildConfigurationListID)
	if err != nil {
		return NativeTarget{}, err
	}

	target := NativeTarget{
		ISA:  "PBXNativeTarget",
		ID:   id,
		Name: name,
		BuildConfigurationList: buildConfigurationList,
	}

	p.nativeTargetByID[id] = target

	return target, nil
}

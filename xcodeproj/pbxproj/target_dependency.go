package pbxproj

// PBXTargetDependency ...
type PBXTargetDependency struct {
	ISA         string
	ID          string
	Target      PBXNativeTarget
	TargetProxy PBXContainerItemProxy
}

var pbxTargetDependencyByID = map[string]PBXTargetDependency{}

// GetPBXTargetDependency ...
func GetPBXTargetDependency(id string, raw map[string]interface{}) PBXTargetDependency {
	if targetDependency, ok := pbxTargetDependencyByID[id]; ok {
		return targetDependency
	}

	rawPBXTargetDependency := raw[id].(map[string]interface{})

	target := GetPBXNativeTarget(rawPBXTargetDependency["target"].(string), raw)
	targetProxy := GetPBXContainerItemProxy(rawPBXTargetDependency["targetProxy"].(string), raw)

	targetDependency := PBXTargetDependency{
		ISA:         "PBXTargetDependency",
		ID:          id,
		Target:      target,
		TargetProxy: targetProxy,
	}
	pbxTargetDependencyByID[id] = targetDependency
	return targetDependency
}

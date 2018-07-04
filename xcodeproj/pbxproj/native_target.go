package pbxproj

// PBXNativeTarget ...
type PBXNativeTarget struct {
	ISA                    string
	ID                     string
	Name                   string
	ProductName            string
	ProductReference       PBXFileReference
	ProductType            string
	BuildConfigurationList XCConfigurationList
	BuildPhases            []PBXSourcesBuildPhase
	BuildRules             []string
	Dependencies           []PBXTargetDependency
}

var pbxNativeTargetByID = map[string]PBXNativeTarget{}

// GetPBXNativeTarget ...
func GetPBXNativeTarget(id string, raw map[string]interface{}) PBXNativeTarget {
	if nativeTarget, ok := pbxNativeTargetByID[id]; ok {
		return nativeTarget
	}

	rawTarget := raw[id].(map[string]interface{})

	name := rawTarget["name"].(string)
	productName := rawTarget["productName"].(string)
	productReference := GetPBXFileReference(rawTarget["productReference"].(string), raw)
	productType := rawTarget["productType"].(string)
	buildConfigurationList := GetXCConfigurationList(rawTarget["buildConfigurationList"].(string), raw)

	target := PBXNativeTarget{
		ISA:                    "PBXNativeTarget",
		ID:                     id,
		Name:                   name,
		ProductName:            productName,
		ProductReference:       productReference,
		ProductType:            productType,
		BuildConfigurationList: buildConfigurationList,
	}
	pbxNativeTargetByID[id] = target
	return target
}

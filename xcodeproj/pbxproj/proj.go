package pbxproj

// PBXProj ...
type PBXProj struct {
	ISA                    string
	ID                     string
	Attributes             map[string]interface{}
	BuildConfigurationList XCConfigurationList
	CompatibilityVersion   string
	DevelopmentRegion      string
	HasScannedForEncodings string
	KnownRegions           []string
	MainGroup              PBXGroup
	ProductRefGroup        PBXGroup
	ProjectDirPath         string
	ProjectRoot            string
	Targets                []PBXNativeTarget
}

// GetPBXProj ...
func GetPBXProj(id string, raw map[string]interface{}) PBXProj {
	rawPBXProj := raw[id].(map[string]interface{})

	attributes := rawPBXProj["attributes"].(map[string]interface{})
	buildConfigurationList := GetXCConfigurationList(rawPBXProj["buildConfigurationList"].(string), raw)
	compatibilityVersion := rawPBXProj["compatibilityVersion"].(string)
	developmentRegion := rawPBXProj["developmentRegion"].(string)
	hasScannedForEncodings := rawPBXProj["hasScannedForEncodings"].(string)
	knownRegions := []string{}
	if _, ok := rawPBXProj["knownRegions"]; ok {
		for _, knownRegion := range rawPBXProj["knownRegions"].([]interface{}) {
			knownRegions = append(knownRegions, knownRegion.(string))
		}
	}
	mainGroup := GetPBXGroup(rawPBXProj["mainGroup"].(string), raw)
	productRefGroup := GetPBXGroup(rawPBXProj["productRefGroup"].(string), raw)
	projectDirPath := rawPBXProj["projectDirPath"].(string)
	projectRoot := rawPBXProj["projectRoot"].(string)
	targets := []PBXNativeTarget{}
	if _, ok := rawPBXProj["targets"]; ok {
		rawTargets := rawPBXProj["targets"].([]interface{})
		for _, rawtarget := range rawTargets {
			targetID := rawtarget.(string)
			target := GetPBXNativeTarget(targetID, raw)
			targets = append(targets, target)
		}
	}

	return PBXProj{
		ISA:                    "PBXProject",
		ID:                     id,
		Attributes:             attributes,
		BuildConfigurationList: buildConfigurationList,
		CompatibilityVersion:   compatibilityVersion,
		DevelopmentRegion:      developmentRegion,
		HasScannedForEncodings: hasScannedForEncodings,
		KnownRegions:           knownRegions,
		MainGroup:              mainGroup,
		ProductRefGroup:        productRefGroup,
		ProjectDirPath:         projectDirPath,
		ProjectRoot:            projectRoot,
		Targets:                targets,
	}
}

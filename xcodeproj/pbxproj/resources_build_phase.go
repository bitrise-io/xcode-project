package pbxproj

// PBXResourcesBuildPhase ...
type PBXResourcesBuildPhase struct {
	ISA                                string
	ID                                 string
	RunOnlyForDeploymentPostprocessing string
	BuildActionMask                    string
	Files                              []PBXBuildFile
}

var pbxResourcesBuildPhaseByID = map[string]PBXResourcesBuildPhase{}

// GetPBXResourcesBuildPhase ...
func GetPBXResourcesBuildPhase(id string, raw map[string]interface{}) PBXResourcesBuildPhase {
	if sourcesBuildPhase, ok := pbxResourcesBuildPhaseByID[id]; ok {
		return sourcesBuildPhase
	}

	rawPBXResourcesBuildPhase := raw[id].(map[string]interface{})

	runOnlyForDeploymentPostprocessing := rawPBXResourcesBuildPhase["runOnlyForDeploymentPostprocessing"].(string)
	buildActionMask := rawPBXResourcesBuildPhase["BuildActionMask"].(string)

	sourcesBuildPhase := PBXResourcesBuildPhase{
		ISA: "PBXResourcesBuildPhase",
		ID:  id,
		RunOnlyForDeploymentPostprocessing: runOnlyForDeploymentPostprocessing,
		BuildActionMask:                    buildActionMask,
	}
	pbxResourcesBuildPhaseByID[id] = sourcesBuildPhase
	return sourcesBuildPhase
}

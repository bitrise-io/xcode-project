package pbxproj

// PBXSourcesBuildPhase ...
type PBXSourcesBuildPhase struct {
	ISA                                string
	ID                                 string
	RunOnlyForDeploymentPostprocessing string
	BuildActionMask                    string
	Files                              []PBXBuildFile
}

var pbxSourcesBuildPhaseByID = map[string]PBXSourcesBuildPhase{}

// GetPBXSourcesBuildPhase ...
func GetPBXSourcesBuildPhase(id string, raw map[string]interface{}) PBXSourcesBuildPhase {
	if sourcesBuildPhase, ok := pbxSourcesBuildPhaseByID[id]; ok {
		return sourcesBuildPhase
	}

	rawPBXSourcesBuildPhase := raw[id].(map[string]interface{})

	runOnlyForDeploymentPostprocessing := rawPBXSourcesBuildPhase["runOnlyForDeploymentPostprocessing"].(string)
	buildActionMask := rawPBXSourcesBuildPhase["BuildActionMask"].(string)

	sourcesBuildPhase := PBXSourcesBuildPhase{
		ISA: "PBXSourcesBuildPhase",
		ID:  id,
		RunOnlyForDeploymentPostprocessing: runOnlyForDeploymentPostprocessing,
		BuildActionMask:                    buildActionMask,
	}
	pbxSourcesBuildPhaseByID[id] = sourcesBuildPhase
	return sourcesBuildPhase
}

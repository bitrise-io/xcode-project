package pbxproj

// PBXFrameworksBuildPhase ...
type PBXFrameworksBuildPhase struct {
	ISA                                string
	ID                                 string
	RunOnlyForDeploymentPostprocessing string
	BuildActionMask                    string
	Files                              []PBXBuildFile
}

var pbxFrameworksBuildPhaseByID = map[string]PBXFrameworksBuildPhase{}

// GetPBXFrameworksBuildPhase ...
func GetPBXFrameworksBuildPhase(id string, raw map[string]interface{}) PBXFrameworksBuildPhase {
	if sourcesBuildPhase, ok := pbxFrameworksBuildPhaseByID[id]; ok {
		return sourcesBuildPhase
	}

	rawPBXFrameworksBuildPhase := raw[id].(map[string]interface{})

	runOnlyForDeploymentPostprocessing := rawPBXFrameworksBuildPhase["runOnlyForDeploymentPostprocessing"].(string)
	buildActionMask := rawPBXFrameworksBuildPhase["BuildActionMask"].(string)

	sourcesBuildPhase := PBXFrameworksBuildPhase{
		ISA: "PBXFrameworksBuildPhase",
		ID:  id,
		RunOnlyForDeploymentPostprocessing: runOnlyForDeploymentPostprocessing,
		BuildActionMask:                    buildActionMask,
	}
	pbxFrameworksBuildPhaseByID[id] = sourcesBuildPhase
	return sourcesBuildPhase
}

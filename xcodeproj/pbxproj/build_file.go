package pbxproj

// PBXBuildFile ...
type PBXBuildFile struct {
	ISA     string
	ID      string
	FileRef PBXFileReference
}

var pbxBuildFileByID = map[string]PBXBuildFile{}

// GetPBXBuildFile ...
func GetPBXBuildFile(id string, raw map[string]interface{}) PBXBuildFile {
	if sourcesBuildPhase, ok := pbxBuildFileByID[id]; ok {
		return sourcesBuildPhase
	}

	rawPBXBuildFile := raw[id].(map[string]interface{})

	fileRef := GetPBXFileReference(rawPBXBuildFile["fileRef"].(string), raw)

	sourcesBuildPhase := PBXBuildFile{
		ISA:     "PBXBuildFile",
		ID:      id,
		FileRef: fileRef,
	}
	pbxBuildFileByID[id] = sourcesBuildPhase
	return sourcesBuildPhase
}

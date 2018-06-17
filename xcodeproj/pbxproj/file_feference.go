package pbxproj

// PBXFileReference ...
type PBXFileReference struct {
	ISA               string
	ID                string
	LastKnownFileType string
	Path              string
	SourceTree        string
	ExplicitFileType  string
	IncludeInIndex    int
}

var pbxFileReferenceByID = map[string]PBXFileReference{}

// GetPBXFileReference ...
func GetPBXFileReference(id string, raw map[string]interface{}) PBXFileReference {
	if fileReference, ok := pbxFileReferenceByID[id]; ok {
		return fileReference
	}

	rawPBXFileReference := raw[id].(map[string]interface{})

	path := rawPBXFileReference["path"].(string)
	sourceTree := rawPBXFileReference["sourceTree"].(string)

	fileReference := PBXFileReference{
		ISA:        "PBXFileReference",
		ID:         id,
		Path:       path,
		SourceTree: sourceTree,
	}
	pbxFileReferenceByID[id] = fileReference
	return fileReference
}

package pbxproj

// PBXGroup ...
type PBXGroup struct {
	ISA        string
	ID         string
	SourceTree string
	Children   []PBXGroup
}

var pbxGroupByID = map[string]PBXGroup{}

// GetPBXGroup ...
func GetPBXGroup(id string, raw map[string]interface{}) PBXGroup {
	if group, ok := pbxGroupByID[id]; ok {
		return group
	}

	rawPBXGroup := raw[id].(map[string]interface{})

	sourceTree := ""
	if _, ok := rawPBXGroup["sourceTree"]; ok {
		sourceTree = rawPBXGroup["sourceTree"].(string)
	}
	children := []PBXGroup{}
	if _, ok := rawPBXGroup["children"]; ok {
		rawChildren := rawPBXGroup["children"].([]interface{})
		for _, rawChild := range rawChildren {
			child := GetPBXGroup(rawChild.(string), raw)
			children = append(children, child)
		}
	}

	group := PBXGroup{
		ISA:        "PBXGroup",
		ID:         id,
		SourceTree: sourceTree,
		Children:   children,
	}
	pbxGroupByID[id] = group
	return group
}

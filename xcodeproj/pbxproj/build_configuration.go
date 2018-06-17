package pbxproj

// XCBuildConfiguration ..
type XCBuildConfiguration struct {
	ISA           string
	ID            string
	Name          string
	BuildSettings map[string]interface{}
}

var xcBuildConfigurationByID = map[string]XCBuildConfiguration{}

// GetXCBuildConfiguration ...
func GetXCBuildConfiguration(id string, raw map[string]interface{}) XCBuildConfiguration {
	if buildConfiguration, ok := xcBuildConfigurationByID[id]; ok {
		return buildConfiguration
	}

	rawXCBuildConfiguration := raw[id].(map[string]interface{})

	buildConfiguration := XCBuildConfiguration{
		ISA:           "XCBuildConfiguration",
		ID:            id,
		Name:          rawXCBuildConfiguration["name"].(string),
		BuildSettings: rawXCBuildConfiguration["buildSettings"].(map[string]interface{}),
	}
	xcBuildConfigurationByID[id] = buildConfiguration
	return buildConfiguration
}

package pbxproj

// XCConfigurationList ...
type XCConfigurationList struct {
	ISA                           string
	ID                            string
	DefaultConfigurationName      string
	DefaultConfigurationIsVisible string
	BuildConfigurations           []XCBuildConfiguration
}

var xcConfigurationListByID = map[string]XCConfigurationList{}

// GetXCConfigurationList ...
func GetXCConfigurationList(id string, raw map[string]interface{}) XCConfigurationList {
	if configurationList, ok := xcConfigurationListByID[id]; ok {
		return configurationList
	}

	rawXCConfigurationList := raw[id].(map[string]interface{})

	buildConfigurations := []XCBuildConfiguration{}
	if _, ok := rawXCConfigurationList["buildConfigurations"]; ok {
		rawBuildConfigurations := rawXCConfigurationList["buildConfigurations"].([]interface{})
		for _, rawID := range rawBuildConfigurations {
			buildConfiguration := GetXCBuildConfiguration(rawID.(string), raw)
			buildConfigurations = append(buildConfigurations, buildConfiguration)
		}
	}

	defaultConfigurationName := ""
	if _, ok := rawXCConfigurationList["defaultConfigurationName"]; ok {
		defaultConfigurationName = rawXCConfigurationList["defaultConfigurationName"].(string)
	}

	defaultConfigurationIsVisible := ""
	if _, ok := rawXCConfigurationList["defaultConfigurationIsVisible"]; ok {
		defaultConfigurationIsVisible = rawXCConfigurationList["defaultConfigurationIsVisible"].(string)
	}

	configurationList := XCConfigurationList{
		ISA: "XCConfigurationList",
		ID:  id,
		DefaultConfigurationName:      defaultConfigurationName,
		DefaultConfigurationIsVisible: defaultConfigurationIsVisible,
		BuildConfigurations:           buildConfigurations,
	}
	xcConfigurationListByID[id] = configurationList
	return configurationList
}

package xcodeproj

// ConfigurationList ...
type ConfigurationList struct {
	ISA                      string
	ID                       string
	DefaultConfigurationName string
	BuildConfigurations      []BuildConfiguration
}

func (p XcodeProj) configurationList(id string) (ConfigurationList, error) {
	if configurationList, ok := p.configurationListByID[id]; ok {
		return configurationList, nil
	}

	raw, err := p.raw.Object(id)
	if err != nil {
		return ConfigurationList{}, err
	}

	rawBuildConfigurations, err := raw.StringSlice("buildConfigurations")
	if err != nil {
		return ConfigurationList{}, err
	}

	buildConfigurations := []BuildConfiguration{}
	for _, rawID := range rawBuildConfigurations {
		buildConfiguration, err := p.buildConfiguration(rawID)
		if err != nil {
			return ConfigurationList{}, err
		}

		buildConfigurations = append(buildConfigurations, buildConfiguration)
	}

	defaultConfigurationName := ""
	if aDefaultConfigurationName, err := raw.String("defaultConfigurationName"); err == nil {
		defaultConfigurationName = aDefaultConfigurationName
	}

	configurationList := ConfigurationList{
		ISA: "XCConfigurationList",
		ID:  id,
		DefaultConfigurationName: defaultConfigurationName,
		BuildConfigurations:      buildConfigurations,
	}

	p.configurationListByID[id] = configurationList

	return configurationList, nil
}

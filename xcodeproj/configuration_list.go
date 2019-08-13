package xcodeproj

import (
	"fmt"

	"github.com/bitrise-io/xcode-project/serialized"
)

// ConfigurationList ...
type ConfigurationList struct {
	ID                       string
	DefaultConfigurationName string
	BuildConfigurations      []BuildConfiguration
}

func parseConfigurationList(id string, objects serialized.Object) (ConfigurationList, error) {
	raw, err := objects.Object(id)
	if err != nil {
		return ConfigurationList{}, err
	}

	rawBuildConfigurations, err := raw.StringSlice("buildConfigurations")
	if err != nil {
		return ConfigurationList{}, err
	}

	var buildConfigurations []BuildConfiguration
	for _, rawID := range rawBuildConfigurations {
		buildConfiguration, err := parseBuildConfiguration(rawID, objects)
		if err != nil {
			return ConfigurationList{}, err
		}

		buildConfigurations = append(buildConfigurations, buildConfiguration)
	}

	var defaultConfigurationName string
	if aDefaultConfigurationName, err := raw.String("defaultConfigurationName"); err == nil {
		defaultConfigurationName = aDefaultConfigurationName
	}

	return ConfigurationList{
		ID:                       id,
		DefaultConfigurationName: defaultConfigurationName,
		BuildConfigurations:      buildConfigurations,
	}, nil
}

// BuildConfigurationList ...
func (p XcodeProj) BuildConfigurationList(targetID string) (serialized.Object, error) {
	objects, err := p.RawProj.Object("objects")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch target buildConfigurationList, the objects of the project are not found, error: %s", err)
	}

	object, err := objects.Object(p.Proj.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch target buildConfigurationList, the project objects with ID (%s) is not found, error: %s", p.Proj.ID, err)
	}
	buildConfigurationListID, err := object.String("buildConfigurationList")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch target's (%s) buildConfigurationList, error: %s", targetID, err)
	}

	buildConfigurationList, err := objects.Object(buildConfigurationListID)
	if err != nil {
		return nil, err
	}

	return buildConfigurationList, nil
}

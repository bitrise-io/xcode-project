package xcodeproj

import (
	"github.com/bitrise-tools/xcode-project/serialized"
)

// Proj ...
type Proj struct {
	ID                     string
	BuildConfigurationList ConfigurationList
	Targets                []Target
}

func parseProj(id string, objects serialized.Object) (Proj, error) {
	rawPBXProj, err := objects.Object(id)
	if err != nil {
		return Proj{}, err
	}

	buildConfigurationListID, err := rawPBXProj.String("buildConfigurationList")
	if err != nil {
		return Proj{}, err
	}

	buildConfigurationList, err := parseConfigurationList(buildConfigurationListID, objects)
	if err != nil {
		return Proj{}, err
	}

	rawTargets, err := rawPBXProj.StringSlice("targets")
	if err != nil {
		return Proj{}, err
	}

	var targets []Target
	for i := range rawTargets {
		target, err := parseTarget(rawTargets[i], objects)
		if err != nil {
			return Proj{}, err
		}
		targets = append(targets, target)
	}

	return Proj{
		ID: id,
		BuildConfigurationList: buildConfigurationList,
		Targets:                targets,
	}, nil
}

// BuildSettings ...
func (p Proj) BuildSettings(configuration string) (serialized.Object, error) {
	for _, buildConfigurationList := range p.BuildConfigurationList.BuildConfigurations {
		if buildConfigurationList.Name == configuration {
			return buildConfigurationList.BuildSettings, nil
		}
	}
	return nil, NewConfigurationNotFoundError(configuration)
}

// Target ...
func (p Proj) Target(id string) (Target, bool) {
	for _, target := range p.Targets {
		if target.ID == id {
			return target, true
		}
	}
	return Target{}, false
}

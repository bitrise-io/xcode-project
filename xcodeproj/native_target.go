package xcodeproj

import "github.com/bitrise-tools/xcode-project/serialized"

// NativeTarget ...
type NativeTarget struct {
	ID                     string
	Name                   string
	BuildConfigurationList ConfigurationList
	Dependencies           []TargetDependency
}

func parseNativeTarget(id string, objects serialized.Object) (NativeTarget, error) {
	rawTarget, err := objects.Object(id)
	if err != nil {
		return NativeTarget{}, err
	}

	name, err := rawTarget.String("name")
	if err != nil {
		return NativeTarget{}, err
	}

	buildConfigurationListID, err := rawTarget.String("buildConfigurationList")
	if err != nil {
		return NativeTarget{}, err
	}

	buildConfigurationList, err := parseConfigurationList(buildConfigurationListID, objects)
	if err != nil {
		return NativeTarget{}, err
	}

	dependencyIDs, err := rawTarget.StringSlice("dependencies")
	if err != nil {
		return NativeTarget{}, err
	}

	var dependencies []TargetDependency
	for _, dependencyID := range dependencyIDs {
		dependency, err := parseTargetDependency(dependencyID, objects)
		if err != nil {
			return NativeTarget{}, err
		}

		dependencies = append(dependencies, dependency)
	}

	return NativeTarget{
		ID:   id,
		Name: name,
		BuildConfigurationList: buildConfigurationList,
		Dependencies:           dependencies,
	}, nil
}

package xcodeproj

import (
	"github.com/bitrise-io/xcode-project/serialized"
)

// Proj ...
type Proj struct {
	ID                     string
	BuildConfigurationList ConfigurationList
	Targets                []Target
	TargetToAssetCatalogs  map[string]map[string]string
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

	targetToAssetCatalogs := make(map[string]map[string]string)
	for _, target := range targets {
		resourcesBuildPhase, err := filterResourcesBuildPhase(target.buildPhaseIDs, objects)
		if err != nil {
			return Proj{}, err
		}
		assetCatalogs, err := filterAssetCatalogs(resourcesBuildPhase, objects)
		if err != nil {
			return Proj{}, err
		}
		targetToAssetCatalogs[target.ID] = assetCatalogs
	}

	return Proj{
		ID:                     id,
		BuildConfigurationList: buildConfigurationList,
		Targets:                targets,
		TargetToAssetCatalogs:  targetToAssetCatalogs,
	}, nil
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

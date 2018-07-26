package xcodeproj

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/bitrise-io/go-utils/fileutil"
	"github.com/bitrise-tools/xcode-project/serialized"
	"howett.net/plist"
)

// TargetType ...
type TargetType string

// TargetTypes
const (
	NativeTargetType    TargetType = "PBXNativeTarget"
	AggregateTargetType TargetType = "PBXAggregateTarget"
	LegacyTargetType    TargetType = "PBXLegacyTarget"
)

// Target ...
type Target struct {
	Type                   TargetType
	ID                     string
	Name                   string
	BuildConfigurationList ConfigurationList
	Dependencies           []TargetDependency
}

// BuildSettings ...
func (t Target) BuildSettings(configuration string) (serialized.Object, error) {
	for _, buildConfigurationList := range t.BuildConfigurationList.BuildConfigurations {
		if buildConfigurationList.Name == configuration {
			return buildConfigurationList.BuildSettings, nil
		}
	}
	return nil, NewConfigurationNotFoundError(configuration)
}

// InformationPropertyListPath ...
func (t Target) InformationPropertyListPath(configuration string, containerDir string) (string, error) {
	buildSettings, err := t.BuildSettings(configuration)
	if err != nil {
		return "", err
	}

	relPth, err := buildSettings.String("INFOPLIST_FILE")
	if err != nil {
		return "", err
	}

	return filepath.Join(containerDir, relPth), nil
}

// InformationPropertyList ...
func (t Target) InformationPropertyList(configuration string, containerDir string) (serialized.Object, error) {
	informationPropertyListPth, err := t.InformationPropertyListPath(configuration, containerDir)
	if err != nil {
		return nil, err
	}

	informationPropertyListContent, err := fileutil.ReadBytesFromFile(informationPropertyListPth)
	if err != nil {
		return nil, err
	}

	var informationPropertyList serialized.Object
	if _, err := plist.Unmarshal([]byte(informationPropertyListContent), &informationPropertyList); err != nil {
		return nil, err
	}

	return informationPropertyList, nil
}

// BundleID ...
func (t Target) BundleID(configuration, containerDir string) (BundleID, error) {
	buildSettings, err := t.BuildSettings(configuration)
	if err != nil {
		return "", err
	}

	bundleID, err := buildSettings.String("PRODUCT_BUNDLE_IDENTIFIER")
	if err != nil && !serialized.IsKeyNotFoundError(err) {
		return "", err
	}

	if bundleID != "" {
		return BundleID(bundleID), nil
	}

	informationPropertyList, err := t.InformationPropertyList(configuration, containerDir)
	if err != nil {
		return "", err
	}

	bundleID, err = informationPropertyList.String("CFBundleIdentifier")
	if err != nil {
		return "", err
	}

	if bundleID == "" {
		return "", errors.New("no PRODUCT_BUNDLE_IDENTIFIER build settings nor CFBundleIdentifier information property found")
	}

	return BundleID(bundleID), nil
}

// DependentTargets ...
func (t Target) DependentTargets() []Target {
	var targets []Target
	for _, targetDependency := range t.Dependencies {
		childTarget := targetDependency.Target
		targets = append(targets, childTarget)

		childDependentTargets := childTarget.DependentTargets()
		targets = append(targets, childDependentTargets...)
	}

	return targets
}

func parseTarget(id string, objects serialized.Object) (Target, error) {
	rawTarget, err := objects.Object(id)
	if err != nil {
		return Target{}, err
	}

	isa, err := rawTarget.String("isa")
	if err != nil {
		return Target{}, err
	}

	var targetType TargetType
	switch isa {
	case "PBXNativeTarget":
		targetType = NativeTargetType
	case "PBXAggregateTarget":
		targetType = AggregateTargetType
	case "PBXLegacyTarget":
		targetType = LegacyTargetType
	default:
		return Target{}, fmt.Errorf("unknown target type: %s", isa)
	}

	name, err := rawTarget.String("name")
	if err != nil {
		return Target{}, err
	}

	buildConfigurationListID, err := rawTarget.String("buildConfigurationList")
	if err != nil {
		return Target{}, err
	}

	buildConfigurationList, err := parseConfigurationList(buildConfigurationListID, objects)
	if err != nil {
		return Target{}, err
	}

	dependencyIDs, err := rawTarget.StringSlice("dependencies")
	if err != nil {
		return Target{}, err
	}

	var dependencies []TargetDependency
	for _, dependencyID := range dependencyIDs {
		dependency, err := parseTargetDependency(dependencyID, objects)
		if err != nil {
			return Target{}, err
		}

		dependencies = append(dependencies, dependency)
	}

	return Target{
		Type: targetType,
		ID:   id,
		Name: name,
		BuildConfigurationList: buildConfigurationList,
		Dependencies:           dependencies,
	}, nil
}

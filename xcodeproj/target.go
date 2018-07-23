package xcodeproj

import (
	"errors"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

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

// BundleID ...
func (t Target) BundleID(configuration, containerDir string) (string, error) {
	var buildConfiguration BuildConfiguration
	for i := range t.BuildConfigurationList.BuildConfigurations {
		if t.BuildConfigurationList.BuildConfigurations[i].Name == configuration {
			buildConfiguration = t.BuildConfigurationList.BuildConfigurations[i]
			break
		}
	}
	if buildConfiguration.Name == "" {
		return "", fmt.Errorf("configuration: %s not found", configuration)
	}

	bundleID, err := buildConfiguration.BuildSettings.String("PRODUCT_BUNDLE_IDENTIFIER")
	if err != nil && !serialized.IsKeyNotFoundError(err) {
		return "", err
	}

	if bundleID == "" {
		infoPlistRelPth, err := buildConfiguration.BuildSettings.String("INFOPLIST_FILE")
		if err != nil {
			return "", fmt.Errorf("no PRODUCT_BUNDLE_IDENTIFIER build settings defined and failed to read target INFOPLIST_FILE: %s", err)
		}

		infoPlistPth := filepath.Join(containerDir, infoPlistRelPth)
		infoPlistContent, err := fileutil.ReadBytesFromFile(infoPlistPth)
		if err != nil {
			return "", fmt.Errorf("no PRODUCT_BUNDLE_IDENTIFIER build settings defined and failed to read the Info.plist file: %s", err)
		}

		var infoPlistObject serialized.Object
		if _, err := plist.Unmarshal([]byte(infoPlistContent), &infoPlistObject); err != nil {
			return "", fmt.Errorf("no PRODUCT_BUNDLE_IDENTIFIER build settings defined and failed to unmarshal the Info.plist file: %s", err)
		}

		bundleID, err = infoPlistObject.String("CFBundleIdentifier")
		if err != nil {
			return "", fmt.Errorf("no PRODUCT_BUNDLE_IDENTIFIER build settings defined and failed to read CFBundleIdentifier from the Info.plist file: %s", err)
		}
	}

	if bundleID == "" {
		return "", errors.New("no PRODUCT_BUNDLE_IDENTIFIER build settings nor CFBundleIdentifier manifest defined")
	}

	for strings.Contains(bundleID, "$") {
		resolvedBundleID, err := resolveBundleID(bundleID, buildConfiguration.BuildSettings)
		if err != nil {
			return "", fmt.Errorf("no PRODUCT_BUNDLE_IDENTIFIER build settings defined and %s", err)
		}

		if resolvedBundleID == bundleID {
			return "", fmt.Errorf("no PRODUCT_BUNDLE_IDENTIFIER build settings defined and failed to resolve CFBundleIdentifier (%s)", bundleID)
		}

		bundleID = resolvedBundleID
	}

	if strings.Contains(bundleID, "$") {
		resolvedBundleID, err := resolveBundleID(bundleID, buildConfiguration.BuildSettings)
		if err != nil {
			return "", fmt.Errorf("no PRODUCT_BUNDLE_IDENTIFIER build settings defined and failed to resolve CFBundleIdentifier (%s): %s", bundleID, err)
		}
		bundleID = resolvedBundleID
	}

	return bundleID, nil
}

func resolveBundleID(bundleID string, buildSettings serialized.Object) (string, error) {
	re := regexp.MustCompile(`(.*)\$\((.*)\)(.*)`)
	matches := re.FindStringSubmatch(bundleID)
	if len(matches) != 4 {
		return "", fmt.Errorf(`failed to resolve bundle id (%s): does not conforms to: (.*)$\(.*\)(.*)`, bundleID)
	}

	prefix := matches[1]
	suffix := matches[3]
	envKey := matches[2]

	split := strings.Split(envKey, ":")
	envKey = split[0]
	envValue, err := buildSettings.String(envKey)
	if err != nil {
		return "", fmt.Errorf("failed to resolve bundle id (%s): build settings not found with key: %s", bundleID, envKey)
	}

	return prefix + envValue + suffix, nil
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

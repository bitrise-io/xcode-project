package xcodeproj

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/bitrise-tools/xcode-project/serialized"
)

// BundleID ...
type BundleID string

// Resolve ...
func (b BundleID) Resolve(buildSettings serialized.Object) (string, error) {
	bundleID := string(b)
	if !strings.Contains(bundleID, "$") {
		return bundleID, nil
	}

	var resolvedBundleIDs []string
	for strings.Contains(bundleID, "$") {
		resolvedBundleID, err := BundleID(bundleID).resolve(buildSettings)
		if err != nil {
			return "", err
		}

		for _, id := range resolvedBundleIDs {
			if id == resolvedBundleID {
				return "", fmt.Errorf("failed to resolve bundle id (%s): reference cycle found", bundleID)
			}
		}
		resolvedBundleIDs = append(resolvedBundleIDs, resolvedBundleID)

		if resolvedBundleID == bundleID {
			return "", fmt.Errorf("failed to resolve bundle id (%s)", bundleID)
		}

		bundleID = resolvedBundleID
	}

	return bundleID, nil
}

func (b BundleID) resolve(buildSettings serialized.Object) (string, error) {
	bundleID := string(b)

	if !strings.Contains(bundleID, "$") {
		return bundleID, nil
	}

	if strings.Contains(bundleID, "$(") {
		// Bitrise.$(PRODUCT_NAME:rfc1034identifier).watch
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

	// Bitrise.$PRODUCT_NAME.watch
	split := strings.Split(bundleID, "$")
	var longestMatchingEnvKey string
	for _, key := range buildSettings.Keys() {
		if strings.HasPrefix(split[1], key) {
			if len(key) > len(longestMatchingEnvKey) {
				longestMatchingEnvKey = key
			}
		}
	}

	if longestMatchingEnvKey == "" {
		return "", fmt.Errorf("failed to resolve bundle id (%s), build settings not found", bundleID)
	}

	envValue, err := buildSettings.String(longestMatchingEnvKey)
	if err != nil {
		return "", err
	}

	return strings.Replace(bundleID, "$"+longestMatchingEnvKey, envValue, -1), nil
}

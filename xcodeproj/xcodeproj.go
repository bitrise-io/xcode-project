package xcodeproj

import (
	"errors"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/bitrise-io/go-utils/fileutil"
	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/bitrise-io/xcode-project/serialized"
	"github.com/bitrise-io/xcode-project/xcodebuild"
	"github.com/bitrise-io/xcode-project/xcscheme"
	"golang.org/x/text/unicode/norm"
	"howett.net/plist"
)

// XcodeProj ...
type XcodeProj struct {
	Proj Proj

	Name string
	Path string
}

func (p XcodeProj) buildSettingsFilePath(target, configuration, key string) (string, error) {
	buildSettings, err := p.TargetBuildSettings(target, configuration)
	if err != nil {
		return "", err
	}

	pth, err := buildSettings.String(key)
	if err != nil {
		return "", err
	}

	if pathutil.IsRelativePath(pth) {
		pth = filepath.Join(filepath.Dir(p.Path), pth)
	}

	return pth, nil
}

// TargetCodeSignEntitlementsPath ...
func (p XcodeProj) TargetCodeSignEntitlementsPath(target, configuration string) (string, error) {
	return p.buildSettingsFilePath(target, configuration, "CODE_SIGN_ENTITLEMENTS")
}

// TargetCodeSignEntitlements ...
func (p XcodeProj) TargetCodeSignEntitlements(target, configuration string) (serialized.Object, error) {
	codeSignEntitlementsPth, err := p.TargetCodeSignEntitlementsPath(target, configuration)
	if err != nil {
		return nil, err
	}

	codeSignEntitlementsContent, err := fileutil.ReadBytesFromFile(codeSignEntitlementsPth)
	if err != nil {
		return nil, err
	}

	var codeSignEntitlements serialized.Object
	if _, err := plist.Unmarshal([]byte(codeSignEntitlementsContent), &codeSignEntitlements); err != nil {
		return nil, err
	}

	return codeSignEntitlements, nil
}

// TargetInformationPropertyListPath ...
func (p XcodeProj) TargetInformationPropertyListPath(target, configuration string) (string, error) {
	return p.buildSettingsFilePath(target, configuration, "INFOPLIST_FILE")
}

// TargetInformationPropertyList ...
func (p XcodeProj) TargetInformationPropertyList(target, configuration string) (serialized.Object, error) {
	informationPropertyListPth, err := p.TargetInformationPropertyListPath(target, configuration)
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

// TargetBundleID ...
func (p XcodeProj) TargetBundleID(target, configuration string) (string, error) {
	buildSettings, err := p.TargetBuildSettings(target, configuration)
	if err != nil {
		return "", err
	}

	bundleID, err := buildSettings.String("PRODUCT_BUNDLE_IDENTIFIER")
	if err != nil && !serialized.IsKeyNotFoundError(err) {
		return "", err
	}

	if bundleID != "" {
		return Resolve(bundleID, buildSettings)
	}

	informationPropertyList, err := p.TargetInformationPropertyList(target, configuration)
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

	return Resolve(bundleID, buildSettings)
}

// Resolve returns the resolved the bundleID. We need this, becaue the bundle ID is not exposed in the .pbxproj file ( raw ).
// If the raw BundleID contains an environment variable we have to replace it.
//
//**Example:**
//BundleID in the .pbxproj: Bitrise.Test.$(PRODUCT_NAME:rfc1034identifier).Suffix
//BundleID after the env is expanded: Bitrise.Test.Sample.Suffix
func Resolve(bundleID string, buildSettings serialized.Object) (string, error) {
	resolvedBundleIDs := map[string]bool{}
	resolved := bundleID
	for true {
		if !strings.Contains(resolved, "$") {
			return resolved, nil
		}

		var err error
		resolved, err = expand(resolved, buildSettings)
		fmt.Printf("\n\nresolved: %s\n\n", resolved)
		if err != nil {
			return "", err
		}

		_, ok := resolvedBundleIDs[resolved]
		if ok {
			return "", fmt.Errorf("bundle id reference cycle found")
		}
		resolvedBundleIDs[resolved] = true
	}
	return "", fmt.Errorf("failed to resolve bundle id: %s", bundleID)
}

func expand(bundleID string, buildSettings serialized.Object) (string, error) {
	// Get the raw env key: $(PRODUCT_NAME:rfc1034identifier) || $(PRODUCT_NAME) || ${PRODUCT_NAME:rfc1034identifier} || ${PRODUCT_NAME}
	r, err := regexp.Compile("[$][{(]?[^.]+[)}]?")
	if err != nil {
		return "", err
	}
	if !r.MatchString(bundleID) {
		return "", fmt.Errorf("failed to match regex [$][{(]?[^.]+[)}]? to %s bundleID", bundleID)
	}

	rawEnvKey := r.FindString(bundleID)

	replacer := strings.NewReplacer("$", "", "(", "", ")", "", "{", "", "}", "")
	envKey := strings.Split(replacer.Replace(rawEnvKey), ":")[0]

	var envValue string
	var removedChar int
	for len(envKey) > 1 {
		for settingsKey := range buildSettings {
			fmt.Printf("envKey: %s\n", envKey)
			if settingsKey == envKey {
				envValue, err = buildSettings.String(envKey)
				if err != nil {
					return "", fmt.Errorf("%s build settings not found", envKey)
				}
				goto END
			}

		}
		envKey = envKey[:len(envKey)-1]
		removedChar++
	}
END:

	if string(rawEnvKey[len(rawEnvKey)-1]) == ")" && removedChar > 0 {
		rawEnvKey = rawEnvKey[:len(rawEnvKey)-removedChar] + ")"
	} else if string(rawEnvKey[len(rawEnvKey)-1]) == "}" && removedChar > 0 {
		rawEnvKey = rawEnvKey[:len(rawEnvKey)-removedChar] + "}"
	} else {
		rawEnvKey = rawEnvKey[:len(rawEnvKey)-removedChar]
	}

	fmt.Printf("rawEnvKey: %s", rawEnvKey)

	// Fetch the env value for the env key
	return strings.Replace(bundleID, rawEnvKey, envValue, -1), nil
}

// TargetBuildSettings ...
func (p XcodeProj) TargetBuildSettings(target, configuration string, customOptions ...string) (serialized.Object, error) {
	return xcodebuild.ShowProjectBuildSettings(p.Path, target, configuration, customOptions...)
}

// Scheme ...
func (p XcodeProj) Scheme(name string) (xcscheme.Scheme, bool) {
	schemes, err := p.Schemes()
	if err != nil {
		return xcscheme.Scheme{}, false
	}

	normName := norm.NFC.String(name)
	for _, scheme := range schemes {
		if norm.NFC.String(scheme.Name) == normName {
			return scheme, true
		}
	}

	return xcscheme.Scheme{}, false
}

// Schemes ...
func (p XcodeProj) Schemes() ([]xcscheme.Scheme, error) {
	return xcscheme.FindSchemesIn(p.Path)
}

// Open ...
func Open(pth string) (XcodeProj, error) {
	absPth, err := pathutil.AbsPath(pth)
	if err != nil {
		return XcodeProj{}, err
	}

	objects, projectID, err := open(pth)

	p, err := parseProj(projectID, objects)
	if err != nil {
		return XcodeProj{}, err
	}

	return XcodeProj{
		Proj: p,
		Path: absPth,
		Name: strings.TrimSuffix(filepath.Base(absPth), filepath.Ext(absPth)),
	}, nil
}

func open(absPth string) (serialized.Object, string, error) {
	pbxProjPth := filepath.Join(absPth, "project.pbxproj")

	b, err := fileutil.ReadBytesFromFile(pbxProjPth)
	if err != nil {
		return serialized.Object{}, "", err
	}

	var raw serialized.Object
	if _, err := plist.Unmarshal(b, &raw); err != nil {
		return serialized.Object{}, "", fmt.Errorf("failed to generate json from Pbxproj - error: %s", err)
	}

	objects, err := raw.Object("objects")
	if err != nil {
		return serialized.Object{}, "", err
	}

	projectID := ""
	for id := range objects {
		object, err := objects.Object(id)
		if err != nil {
			return serialized.Object{}, "", err
		}

		objectISA, err := object.String("isa")
		if err != nil {
			return serialized.Object{}, "", err
		}

		if objectISA == "PBXProject" {
			projectID = id
			break
		}
	}
	return objects, projectID, nil
}

// IsXcodeProj ...
func IsXcodeProj(pth string) bool {
	return filepath.Ext(pth) == ".xcodeproj"
}

package xcodeproj

import (
	"path/filepath"
	"reflect"
	"testing"

	"github.com/bitrise-io/xcode-project/pretty"
	"github.com/bitrise-io/xcode-project/serialized"
	"github.com/bitrise-io/xcode-project/testhelper"
	"github.com/bitrise-io/xcode-project/xcscheme"
	"github.com/stretchr/testify/require"
	"golang.org/x/text/unicode/norm"
)

func TestResolve(t *testing.T) {

	t.Log("resolves bundle id in format: prefix.${ENV_KEY}.$ENV_KEY_2")
	{
		bundleID := `prefix.${PRODUCT_NAME}.$VERSION`
		buildSettings := serialized.Object{
			"PRODUCT_NAME": "ios-simple-objc",
			"VERSION":      "beta",
		}
		resolved, err := Resolve(bundleID, buildSettings)
		require.NoError(t, err)
		require.Equal(t, "prefix.ios-simple-objc.beta", resolved)
	}

	t.Log("resolves bundle id in format: prefix.{text.${ENV_KEY}.text}")
	{
		bundleID := `prefix.{text.${PRODUCT_NAME}.text}`
		buildSettings := serialized.Object{
			"PRODUCT_NAME": "ios-simple-objc",
		}
		resolved, err := Resolve(bundleID, buildSettings)
		require.NoError(t, err)
		require.Equal(t, "prefix.{text.ios-simple-objc.text}", resolved)
	}

	t.Log("resolves bundle id in format: prefix.$ENV_KEY")
	{
		bundleID := `auto_provision.$PRODUCT_NAME`
		buildSettings := serialized.Object{
			"PRODUCT_NAME": "ios-simple-objc",
		}
		resolved, err := Resolve(bundleID, buildSettings)
		require.NoError(t, err)
		require.Equal(t, "auto_provision.ios-simple-objc", resolved)
	}

	t.Log("resolves bundle id in format: prefix.$ENV_KEYsuffix")
	{
		bundleID := `auto_provision.$PRODUCT_NAMEsuffix`
		buildSettings := serialized.Object{
			"PRODUCT_NAME": "ios-simple-objc",
		}
		resolved, err := Resolve(bundleID, buildSettings)
		require.NoError(t, err)
		require.Equal(t, "auto_provision.ios-simple-objcsuffix", resolved)
	}

	t.Log("resolves bundle id in format: prefix.$ENV_KEYsuffix$ENV_KEY")
	{
		bundleID := `auto_provision.$PRODUCT_NAMEsuffix$PRODUCT_NAME`
		buildSettings := serialized.Object{
			"PRODUCT_NAME": "ios-simple-objc",
		}
		resolved, err := Resolve(bundleID, buildSettings)
		require.NoError(t, err)
		require.Equal(t, "auto_provision.ios-simple-objcsuffixios-simple-objc", resolved)
	}

	t.Log("resolves bundle id in format: prefix.$ENV_KEY$ENV_KEY_2")
	{
		bundleID := `auto_provision.$PRODUCT_NAME$VERSION`
		buildSettings := serialized.Object{
			"PRODUCT_NAME": "ios-simple-objc",
			"VERSION":      "beta",
		}
		resolved, err := Resolve(bundleID, buildSettings)
		require.NoError(t, err)
		require.Equal(t, "auto_provision.ios-simple-objcbeta", resolved)
	}

	t.Log("resolves bundle id in format: prefix.$(ENV_KEY:rfc1034identifier)")
	{
		bundleID := `auto_provision.$(PRODUCT_NAME:rfc1034identifier)`
		buildSettings := serialized.Object{
			"PRODUCT_NAME": "ios-simple-objc",
		}
		resolved, err := Resolve(bundleID, buildSettings)
		require.NoError(t, err)
		require.Equal(t, "auto_provision.ios-simple-objc", resolved)
	}

	t.Log("resolves bundle id in format: prefix.$ENV_KEYtest.suffix")
	{
		bundleID := `auto_provision.$PRODUCT_NAMEtest.suffix`
		buildSettings := serialized.Object{
			"PRODUCT_NAME": "ios-simple-objc",
		}
		resolved, err := Resolve(bundleID, buildSettings)
		require.NoError(t, err)
		require.Equal(t, "auto_provision.ios-simple-objctest.suffix", resolved)
	}

	t.Log("resolves bundle id with cross env reference")
	{
		bundleID := `auto_provision.$(BUNDLE_ID:rfc1034identifier)`
		buildSettings := serialized.Object{
			"PRODUCT_NAME": "ios-simple-objc",
			"BUNDLE_ID":    "$(PRODUCT_NAME:rfc1034identifier)",
		}
		resolved, err := Resolve(bundleID, buildSettings)
		require.NoError(t, err)
		require.Equal(t, "auto_provision.ios-simple-objc", resolved)
	}

	t.Log("detects env refernce cycle")
	{
		bundleID := `auto_provision.$(BUNDLE_ID:rfc1034identifier)`
		buildSettings := serialized.Object{
			"PRODUCT_NAME": "$(BUNDLE_ID:rfc1034identifier)",
			"BUNDLE_ID":    "$(PRODUCT_NAME:rfc1034identifier)",
		}
		resolved, err := Resolve(bundleID, buildSettings)
		require.EqualError(t, err, "bundle id reference cycle found")
		require.Equal(t, "", resolved)
	}

	t.Log("resolves bundle id in format: prefix.$(ENV_KEY:rfc1034identifier).suffix")
	{
		bundleID := `auto_provision.$(PRODUCT_NAME:rfc1034identifier).suffix`
		buildSettings := serialized.Object{
			"PRODUCT_NAME": "ios-simple-objc",
		}
		resolved, err := Resolve(bundleID, buildSettings)
		require.NoError(t, err)
		require.Equal(t, "auto_provision.ios-simple-objc.suffix", resolved)
	}

	t.Log("resolves bundle id in format: $(ENV_KEY:rfc1034identifier)")
	{
		bundleID := `$(PRODUCT_NAME:rfc1034identifier)`
		buildSettings := serialized.Object{
			"PRODUCT_NAME": "ios-simple-objc",
		}
		resolved, err := Resolve(bundleID, buildSettings)
		require.NoError(t, err)
		require.Equal(t, "ios-simple-objc", resolved)
	}

	t.Log("resolves bundle id in format: prefix.$(ENV_KEY:rfc1034identifier)")
	{
		bundleID := `auto_provision.$(PRODUCT_NAME:rfc1034identifier)`
		buildSettings := serialized.Object{
			"PRODUCT_NAME": "ios-simple-objc",
		}
		resolved, err := Resolve(bundleID, buildSettings)
		require.NoError(t, err)
		require.Equal(t, "auto_provision.ios-simple-objc", resolved)
	}

	t.Log("resolves bundle id with cross env reference")
	{
		bundleID := `auto_provision.$(BUNDLE_ID:rfc1034identifier)`
		buildSettings := serialized.Object{
			"PRODUCT_NAME": "ios-simple-objc",
			"BUNDLE_ID":    "$(PRODUCT_NAME:rfc1034identifier)",
		}
		resolved, err := Resolve(bundleID, buildSettings)
		require.NoError(t, err)
		require.Equal(t, "auto_provision.ios-simple-objc", resolved)
	}

	t.Log("detects env refernce cycle")
	{
		bundleID := `auto_provision.${BUNDLE_ID:rfc1034identifier}`
		buildSettings := serialized.Object{
			"PRODUCT_NAME": "${BUNDLE_ID:rfc1034identifier}",
			"BUNDLE_ID":    "${PRODUCT_NAME:rfc1034identifier}",
		}
		resolved, err := Resolve(bundleID, buildSettings)
		require.EqualError(t, err, "bundle id reference cycle found")
		require.Equal(t, "", resolved)
	}

	t.Log("resolves bundle id in format: prefix.${ENV_KEY:rfc1034identifier}.suffix")
	{
		bundleID := `auto_provision.${PRODUCT_NAME:rfc1034identifier}.suffix`
		buildSettings := serialized.Object{
			"PRODUCT_NAME": "ios-simple-objc",
		}
		resolved, err := Resolve(bundleID, buildSettings)
		require.NoError(t, err)
		require.Equal(t, "auto_provision.ios-simple-objc.suffix", resolved)
	}

	t.Log("resolves bundle id in format: ${ENV_KEY:rfc1034identifier}")
	{
		bundleID := `${PRODUCT_NAME:rfc1034identifier}`
		buildSettings := serialized.Object{
			"PRODUCT_NAME": "ios-simple-objc",
		}
		resolved, err := Resolve(bundleID, buildSettings)
		require.NoError(t, err)
		require.Equal(t, "ios-simple-objc", resolved)
	}

	t.Log("resolves bundle id in format: prefix.$(ENV_KEY:rfc1034identifier).suffix.$(ENV_KEY:rfc1034identifier)")
	{
		bundleID := `auto_provision.$(PRODUCT_NAME:rfc1034identifier).suffix.$(PRODUCT_NAME:rfc1034identifier)`
		buildSettings := serialized.Object{
			"PRODUCT_NAME": "ios-simple-objc",
		}
		resolved, err := Resolve(bundleID, buildSettings)
		require.NoError(t, err)
		require.Equal(t, "auto_provision.ios-simple-objc.suffix.ios-simple-objc", resolved)
	}

	t.Log("resolves bundle id in format: prefix.$(ENV_KEY:rfc1034identifier).suffix.$(ENV_KEY_2:rfc1034identifier)")
	{
		bundleID := `auto_provision.$(PRODUCT_NAME:rfc1034identifier).suffix.$(VERSION:rfc1034identifier)`
		buildSettings := serialized.Object{
			"PRODUCT_NAME": "ios-simple-objc",
			"VERSION":      "beta",
		}
		resolved, err := Resolve(bundleID, buildSettings)
		require.NoError(t, err)
		require.Equal(t, "auto_provision.ios-simple-objc.suffix.beta", resolved)
	}

	t.Log("resolves bundle id in format: prefix.$ENV_KEY.suffix")
	{
		bundleID := `auto_provision.$PRODUCT_NAME.suffix`
		buildSettings := serialized.Object{
			"PRODUCT_NAME": "ios-simple-objc",
		}
		resolved, err := Resolve(bundleID, buildSettings)
		require.NoError(t, err)
		require.Equal(t, "auto_provision.ios-simple-objc.suffix", resolved)
	}
	t.Log("resolves bundle id in format: prefix.second.${ENV_KEY}.suffix")
	{
		bundleID := `prefix.second.${PRODUCT_NAME}.suffix`
		buildSettings := serialized.Object{
			"PRODUCT_NAME": "ios-simple-objc",
		}
		resolved, err := Resolve(bundleID, buildSettings)
		require.NoError(t, err)
		require.Equal(t, "prefix.second.ios-simple-objc.suffix", resolved)
	}
	t.Log("resolves bundle id in format: prefix.second.third.${ENV_KEY}")
	{
		bundleID := `prefix.second.third.${PRODUCT_NAME}`
		buildSettings := serialized.Object{
			"PRODUCT_NAME": "ios-simple-objc",
		}
		resolved, err := Resolve(bundleID, buildSettings)
		require.NoError(t, err)
		require.Equal(t, "prefix.second.third.ios-simple-objc", resolved)
	}
	t.Log("resolves bundle id in format: prefix.second.third.fourth.${ENV_KEY}")
	{
		bundleID := `prefix.second.third.fourth.${PRODUCT_NAME}`
		buildSettings := serialized.Object{
			"PRODUCT_NAME": "ios-simple-objc",
		}
		resolved, err := Resolve(bundleID, buildSettings)
		require.NoError(t, err)
		require.Equal(t, "prefix.second.third.fourth.ios-simple-objc", resolved)
	}
}

func TestExpand(t *testing.T) {

	// Complex env
	t.Log("resolves bundle id in format: prefix.$(ENV_KEY:rfc1034identifier).suffix.$(ENV_KEY:rfc1034identifier)")
	{
		bundleID := `auto_provision.$(PRODUCT_NAME:rfc1034identifier).suffix.$(PRODUCT_NAME:rfc1034identifier)`
		buildSettings := serialized.Object{
			"PRODUCT_NAME": "ios-simple-objc",
		}
		resolved, err := expand(bundleID, buildSettings)
		require.NoError(t, err)
		require.Equal(t, "auto_provision.ios-simple-objc.suffix.ios-simple-objc", resolved)
	}

	t.Log("resolves bundle id in format: prefix.$(ENV_KEY:rfc1034identifier).suffix")
	{
		bundleID := `auto_provision.$(PRODUCT_NAME:rfc1034identifier).suffix`
		buildSettings := serialized.Object{
			"PRODUCT_NAME": "ios-simple-objc",
		}
		resolved, err := expand(bundleID, buildSettings)
		require.NoError(t, err)
		require.Equal(t, "auto_provision.ios-simple-objc.suffix", resolved)
	}

	// Simple env
	t.Log("resolves bundle id in format: prefix.$ENV_KEY")
	{
		bundleID := `auto_provision.$PRODUCT_NAME`
		buildSettings := serialized.Object{
			"PRODUCT_NAME": "ios-simple-objc",
		}
		resolved, err := expand(bundleID, buildSettings)
		require.NoError(t, err)
		require.Equal(t, "auto_provision.ios-simple-objc", resolved)
	}

	t.Log("resolves bundle id in format: prefix.$ENV_KEYsuffix")
	{
		bundleID := `auto_provision.$PRODUCT_NAMEsuffix`
		buildSettings := serialized.Object{
			"PRODUCT_NAME": "ios-simple-objc",
		}
		resolved, err := expand(bundleID, buildSettings)
		require.NoError(t, err)
		require.Equal(t, "auto_provision.ios-simple-objcsuffix", resolved)
	}
}

func TestExpandComplexEnv(t *testing.T) {
	t.Log("resolves bundle id in format: prefix.{text.${ENV_KEY}.text}")
	{
		bundleID := `prefix.{text.${PRODUCT_NAME}.text}`
		buildSettings := serialized.Object{
			"PRODUCT_NAME": "ios-simple-objc",
		}
		resolved, err := expandComplexEnv(bundleID, buildSettings)
		require.NoError(t, err)
		require.Equal(t, "prefix.{text.ios-simple-objc.text}", resolved)
	}

	t.Log("resolves bundle id in format: prefix.$(ENV_KEY:rfc1034identifier).suffix.$(ENV_KEY:rfc1034identifier)")
	{
		bundleID := `auto_provision.$(PRODUCT_NAME:rfc1034identifier).suffix.$(PRODUCT_NAME:rfc1034identifier)`
		buildSettings := serialized.Object{
			"PRODUCT_NAME": "ios-simple-objc",
		}
		resolved, err := expandComplexEnv(bundleID, buildSettings)
		require.NoError(t, err)
		require.Equal(t, "auto_provision.ios-simple-objc.suffix.ios-simple-objc", resolved)
	}

	t.Log("resolves bundle id in format: prefix.$(ENV_KEY:rfc1034identifier).suffix")
	{
		bundleID := `auto_provision.$(PRODUCT_NAME:rfc1034identifier).suffix`
		buildSettings := serialized.Object{
			"PRODUCT_NAME": "ios-simple-objc",
		}
		resolved, err := expandComplexEnv(bundleID, buildSettings)
		require.NoError(t, err)
		require.Equal(t, "auto_provision.ios-simple-objc.suffix", resolved)
	}
}

func TestExpandSimpleEnv(t *testing.T) {
	t.Log("resolves bundle id in format: prefix.$ENV_KEY")
	{
		bundleID := `auto_provision.$PRODUCT_NAME`
		buildSettings := serialized.Object{
			"PRODUCT_NAME": "ios-simple-objc",
		}
		resolved, err := expandSimpleEnv(bundleID, buildSettings)
		require.NoError(t, err)
		require.Equal(t, "auto_provision.ios-simple-objc", resolved)
	}

	t.Log("resolves bundle id in format: prefix.$ENV_KEYsuffix")
	{
		bundleID := `auto_provision.$PRODUCT_NAMEsuffix`
		buildSettings := serialized.Object{
			"PRODUCT_NAME": "ios-simple-objc",
		}
		resolved, err := expandSimpleEnv(bundleID, buildSettings)
		require.NoError(t, err)
		require.Equal(t, "auto_provision.ios-simple-objcsuffix", resolved)
	}
}

func TestTargets(t *testing.T) {
	dir := testhelper.GitCloneIntoTmpDir(t, "https://github.com/bitrise-io/xcode-project-test.git")
	project, err := Open(filepath.Join(dir, "Group/SubProject/SubProject.xcodeproj"))
	require.NoError(t, err)

	{
		target, ok := project.Proj.Target("7D0342D720F4B5AD0050B6A6")
		require.True(t, ok)

		dependentTargets := target.DependentTargets()
		require.Equal(t, 2, len(dependentTargets))
		require.Equal(t, "WatchKitApp", dependentTargets[0].Name)
		require.Equal(t, "WatchKitApp Extension", dependentTargets[1].Name)

		dependentExecutableTarget := target.DependentExecutableProductTargets(false)
		require.Equal(t, 2, len(dependentExecutableTarget))
		require.Equal(t, "WatchKitApp", dependentExecutableTarget[0].Name)
		require.Equal(t, "WatchKitApp Extension", dependentExecutableTarget[1].Name)
	}

	{
		settings, err := project.TargetBuildSettings("SubProject", "Debug")
		require.NoError(t, err)
		require.True(t, len(settings) > 0)

		bundleID, err := settings.String("PRODUCT_BUNDLE_IDENTIFIER")
		require.NoError(t, err)
		require.Equal(t, "com.bitrise.SubProject", bundleID)

		infoPlist, err := settings.String("INFOPLIST_PATH")
		require.NoError(t, err)
		require.Equal(t, "SubProject.app/Info.plist", infoPlist)
	}

	{
		bundleID, err := project.TargetBundleID("SubProject", "Debug")
		require.NoError(t, err)
		require.Equal(t, "com.bitrise.SubProject", bundleID)
	}

	{
		properties, err := project.TargetInformationPropertyList("SubProject", "Debug")
		require.NoError(t, err)
		require.Equal(t, serialized.Object{"CFBundlePackageType": "APPL",
			"UISupportedInterfaceOrientations":      []interface{}{"UIInterfaceOrientationPortrait", "UIInterfaceOrientationLandscapeLeft", "UIInterfaceOrientationLandscapeRight"},
			"CFBundleInfoDictionaryVersion":         "6.0",
			"CFBundleName":                          "$(PRODUCT_NAME)",
			"UISupportedInterfaceOrientations~ipad": []interface{}{"UIInterfaceOrientationPortrait", "UIInterfaceOrientationPortraitUpsideDown", "UIInterfaceOrientationLandscapeLeft", "UIInterfaceOrientationLandscapeRight"},
			"CFBundleDevelopmentRegion":             "$(DEVELOPMENT_LANGUAGE)",
			"CFBundleExecutable":                    "$(EXECUTABLE_NAME)",
			"CFBundleShortVersionString":            "1.0",
			"CFBundleVersion":                       "1",
			"LSRequiresIPhoneOS":                    true,
			"UIMainStoryboardFile":                  "Main",
			"UIRequiredDeviceCapabilities":          []interface{}{"armv7"},
			"CFBundleIdentifier":                    "$(PRODUCT_BUNDLE_IDENTIFIER)",
			"UILaunchStoryboardName":                "LaunchScreen"}, properties)
	}

	{
		entitlements, err := project.TargetCodeSignEntitlements("WatchKitApp", "Debug")
		require.NoError(t, err)
		require.Equal(t, serialized.Object{"com.apple.security.application-groups": []interface{}{}}, entitlements)

	}
}

func TestScheme(t *testing.T) {
	dir := testhelper.GitCloneIntoTmpDir(t, "https://github.com/bitrise-io/xcode-project-test.git")
	pth := filepath.Join(dir, "XcodeProj.xcodeproj")
	project, err := Open(pth)
	require.NoError(t, err)

	{
		scheme, container, err := project.Scheme("ProjectTodayExtensionScheme")
		require.NoError(t, err)
		require.Equal(t, "ProjectTodayExtensionScheme", scheme.Name)
		require.Equal(t, pth, container)
	}

	{
		scheme, container, err := project.Scheme("NotExistScheme")
		require.EqualError(t, err, "scheme NotExistScheme not found in XcodeProj")
		require.Equal(t, (*xcscheme.Scheme)(nil), scheme)
		require.Equal(t, "", container)
	}

	{
		// Gdańsk represented in High Sierra
		b := []byte{71, 100, 97, 197, 132, 115, 107}
		scheme, container, err := project.Scheme(string(b))
		require.NoError(t, err)
		require.Equal(t, norm.NFC.String(string(b)), norm.NFC.String(scheme.Name))
		require.Equal(t, pth, container)
	}

	{
		// Gdańsk represented in Mojave
		b := []byte{71, 100, 97, 110, 204, 129, 115, 107}
		scheme, container, err := project.Scheme(string(b))
		require.NoError(t, err)
		require.Equal(t, norm.NFC.String(string(b)), norm.NFC.String(scheme.Name))
		require.Equal(t, pth, container)
	}
}

func TestSchemes(t *testing.T) {
	dir := testhelper.GitCloneIntoTmpDir(t, "https://github.com/bitrise-io/xcode-project-test.git")
	project, err := Open(filepath.Join(dir, "XcodeProj.xcodeproj"))
	require.NoError(t, err)

	schemes, err := project.Schemes()
	require.NoError(t, err)
	require.Equal(t, 3, len(schemes))

	// Gdańsk represented in High Sierra
	b := []byte{71, 100, 97, 197, 132, 115, 107}
	require.Equal(t, norm.NFC.String(string(b)), norm.NFC.String(schemes[0].Name))
	require.Equal(t, "ProjectScheme", schemes[1].Name)

	// Gdańsk represented in Mojave
	b = []byte{71, 100, 97, 110, 204, 129, 115, 107}
	require.Equal(t, norm.NFC.String(string(b)), norm.NFC.String(schemes[0].Name))
	require.Equal(t, "ProjectScheme", schemes[1].Name)
}

func TestOpenXcodeproj(t *testing.T) {
	t.Log("Opening Pods.xcodeproj in sample-apps-ios-workspace-swift.git")
	{
		dir := testhelper.GitCloneIntoTmpDir(t, "https://github.com/bitrise-io/sample-apps-ios-workspace-swift.git")
		project, err := Open(filepath.Join(dir, "Pods", "Pods.xcodeproj"))
		require.NoError(t, err)
		require.Equal(t, filepath.Join(dir, "Pods", "Pods.xcodeproj"), project.Path)
		require.Equal(t, "Pods", project.Name)
	}
	t.Log("Opening XcodeProj.xcodeproj in xcode-project-test.git")
	{
		dir := testhelper.GitCloneIntoTmpDir(t, "https://github.com/bitrise-io/xcode-project-test.git")
		project, err := Open(filepath.Join(dir, "XcodeProj.xcodeproj"))
		require.NoError(t, err)
		require.Equal(t, filepath.Join(dir, "XcodeProj.xcodeproj"), project.Path)
		require.Equal(t, "XcodeProj", project.Name)
	}
}

func TestIsXcodeProj(t *testing.T) {
	require.True(t, IsXcodeProj("./BitriseSample.xcodeproj"))
	require.False(t, IsXcodeProj("./BitriseSample.xcworkspace"))
}

func TestXcodeProj_forceCodeSign(t *testing.T) {
	dir := testhelper.GitCloneIntoTmpDir(t, "https://github.com/bitrise-io/xcode-project-test.git")
	project, err := Open(filepath.Join(dir, "XcodeProj.xcodeproj"))
	if err != nil {
		t.Fatalf("Failed to init project for test case, error: %s", err)
	}
	tests := []struct {
		name                    string
		configuration           string
		developmentTeam         string
		targetName              string
		codeSignIdentity        string
		provisioningProfileUUID string
		wantErr                 bool
	}{
		{
			name:                    "Force code sign - XcodeProj",
			configuration:           "Release",
			developmentTeam:         "72SA8V3WYL",
			targetName:              "XcodeProj",
			codeSignIdentity:        "iPhone Developer: Test",
			provisioningProfileUUID: "",
			wantErr:                 false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := project.ForceCodeSign(tt.configuration, tt.targetName, tt.developmentTeam, tt.codeSignIdentity, tt.provisioningProfileUUID); (err != nil) != tt.wantErr {
				t.Errorf("XcodeProj.forceCodeSign() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}

	if err := project.Save(); err != nil {
		t.Errorf("Failed to save project, error: %s", err)
	}
}

func TestXcodeProj_foreceCodeSignOnTargetAttributes(t *testing.T) {
	dir := testhelper.GitCloneIntoTmpDir(t, "https://github.com/bitrise-io/xcode-project-test.git")
	project, err := Open(filepath.Join(dir, "XcodeProj.xcodeproj"))
	if err != nil {
		t.Fatalf("Failed to init project for test case, error: %s", err)
	}
	tests := []struct {
		name             string
		targetAttributes serialized.Object
		developmentTeam  string
		targetID         string
		want             serialized.Object
		wantErr          bool
	}{
		{
			name: "Force code sign - XcodeProj",
			targetAttributes: func() serialized.Object {
				targetAttributes, err := project.TargetAttributes()
				if err != nil {
					t.Fatalf("Failed to fetch TargetAttributes for test case, error: %s", err)
				}
				return targetAttributes
			}(),
			developmentTeam: "72SA8V3WYL",
			targetID:        "7D5B35FB20E28EE80022BAE6",
			want: map[string]interface{}{
				"7D0342F020F4BA280050B6A6": map[string]interface{}{
					"CreatedOnToolsVersion": "9.4.1",
					"TestTargetID":          "7D5B35FB20E28EE80022BAE6",
				},
				"7D03430C20F4BB070050B6A6": map[string]interface{}{
					"CreatedOnToolsVersion": "9.4.1",
					"SystemCapabilities": map[string]interface{}{
						"com.apple.Push": map[string]interface{}{
							"enabled": "1",
						},
						"com.apple.iCloud": map[string]interface{}{
							"enabled": "1",
						},
					},
				},
				"7D5B35FB20E28EE80022BAE6": map[string]interface{}{
					"CreatedOnToolsVersion": "9.4.1",
					"DevelopmentTeam":       "72SA8V3WYL",
					"DevelopmentTeamName":   "",
					"ProvisioningStyle":     "Manual",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := foreceCodeSignOnTargetAttributes(tt.targetAttributes, tt.targetID, tt.developmentTeam)
			if (err != nil) != tt.wantErr {
				t.Errorf("XcodeProj.foreceCodeSignOnTargetAttributes() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(tt.targetAttributes, tt.want) {
				t.Errorf("XcodeProj.foreceCodeSignOnTargetAttributes() got = %s, wantErr %s", pretty.Object(tt.targetAttributes), pretty.Object(tt.want))
				return
			}
		})
	}
}

func TestXcodeProj_forceBundleID(t *testing.T) {
	dir := testhelper.GitCloneIntoTmpDir(t, "https://github.com/bitrise-io/xcode-project-test.git")
	project, err := Open(filepath.Join(dir, "XcodeProj.xcodeproj"))
	if err != nil {
		t.Fatalf("Failed to init project for test case, error: %s", err)
	}

	tests := []struct {
		name          string
		target        string
		configuration string
		bundleID      string
		wantErr       bool
	}{
		{
			name:          "Update bundle ID for target and configuration",
			target:        "XcodeProj",
			configuration: "Release",
			bundleID:      "io.bitrise.test.XcodeProj",
			wantErr:       false,
		},
		{
			name:    "Target not found",
			target:  "NON_EXISTENT_TARGET",
			wantErr: true,
		},
		{
			name:          "Configuration not found",
			configuration: "NON_EXISTENT_CONFIGURATION",
			wantErr:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			err := project.ForceTargetBundleID(tt.target, tt.configuration, tt.bundleID)
			if (err != nil) != tt.wantErr {
				t.Fatalf("got error: %s", err)
			}

			if got, err := project.TargetBundleID(tt.target, tt.configuration); (err != nil) != tt.wantErr {
				t.Fatalf("error validating test: %s", err)

			} else if err == nil && got != tt.bundleID {
				t.Fatalf("%s, %s", got, tt.bundleID)
			}

		})
	}
}

func TestXcodePrj_forceTargetCodeSignEntitlement(t *testing.T) {
	dir := testhelper.GitCloneIntoTmpDir(t, "https://github.com/bitrise-io/xcode-project-test.git")
	project, err := Open(filepath.Join(dir, "XcodeProj.xcodeproj"))
	if err != nil {
		t.Fatalf("Failed to init project for test case, error: %s", err)
	}

	tests := []struct {
		name          string
		target        string
		configuration string
		entitlement   string
		value         string
		wantErr       bool
	}{
		{
			name:          "Update entitlement",
			target:        "TodayExtension",
			configuration: "Release",
			entitlement:   "com.apple.security.application-groups",
			value:         "io.bitrise.test",
			wantErr:       false,
		},
		{
			name:    "Target not found",
			target:  "NON_EXISTENT_TARGET",
			wantErr: true,
		},
		{
			name:          "Configuration not found",
			configuration: "NON_EXISTENT_CONFIGURATION",
			wantErr:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			err := project.ForceTargetCodeSignEntitlement(tt.target, tt.configuration, tt.entitlement, tt.value)
			if (err != nil) != tt.wantErr {
				t.Fatalf("got error: %s, but wanErr = %t", err, tt.wantErr)
			}

			if got, err := project.TargetCodeSignEntitlements(tt.target, tt.configuration); (err != nil) != tt.wantErr {
				t.Fatalf("error validating test: %s", err)
			} else if err == nil && got[tt.entitlement] != tt.value {
				t.Fatalf("got %s, want %s", got[tt.entitlement], tt.value)
			}
		})
	}

}

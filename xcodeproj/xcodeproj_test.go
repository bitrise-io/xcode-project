package xcodeproj

import (
	"path/filepath"
	"reflect"
	"testing"

	"github.com/bitrise-io/xcode-project/pretty"

	"github.com/bitrise-io/xcode-project/serialized"
	"github.com/bitrise-io/xcode-project/testhelper"
	"github.com/stretchr/testify/require"
	"golang.org/x/text/unicode/norm"
)

func TestResolve(t *testing.T) {
	t.Log("resolves bundle id in format: prefix.$(ENV_KEY:rfc1034identifier).suffix")
	{
		bundleID := `auto_provision.$(PRODUCT_NAME:rfc1034identifier)`
		buildSettings := serialized.Object{
			"PRODUCT_NAME": "ios-simple-objc",
		}
		resolved, err := resolve(bundleID, buildSettings)
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
		resolved, err := resolve(bundleID, buildSettings)
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
		resolved, err := resolve(bundleID, buildSettings)
		require.EqualError(t, err, "bundle id reference cycle found")
		require.Equal(t, "", resolved)
	}
}

func TestExpand(t *testing.T) {
	t.Log("resolves bundle id in format: prefix.$(ENV_KEY:rfc1034identifier).suffix")
	{
		bundleID := `auto_provision.$(PRODUCT_NAME:rfc1034identifier)`
		buildSettings := serialized.Object{
			"PRODUCT_NAME": "ios-simple-objc",
		}
		resolved, err := expand(bundleID, buildSettings)
		require.NoError(t, err)
		require.Equal(t, "auto_provision.ios-simple-objc", resolved)
	}

	t.Log("works without env")
	{
		bundleID := `auto_provision.ios-simple-objc`
		buildSettings := serialized.Object{
			"PRODUCT_NAME": "ios-simple-objc",
		}
		resolved, err := expand(bundleID, buildSettings)
		require.NoError(t, err)
		require.Equal(t, "auto_provision.ios-simple-objc", resolved)
	}

	t.Log("fails if bundle id does not conforms to: (.*)\\$\\((.*)\\)(.*)")
	{
		bundleID := `auto_provision.$PRODUCT_NAME`
		buildSettings := serialized.Object{
			"PRODUCT_NAME": "ios-simple-objc",
		}
		resolved, err := expand(bundleID, buildSettings)
		require.EqualError(t, err, "auto_provision.$PRODUCT_NAME does not match to pattern: (.*)\\$\\((.*)\\)(.*)")
		require.Equal(t, "", resolved)
	}

	t.Log("fails if env not found")
	{
		bundleID := `auto_provision.$(PRODUCT_NAME:rfc1034identifier)`
		buildSettings := serialized.Object{}
		resolved, err := expand(bundleID, buildSettings)
		require.EqualError(t, err, "PRODUCT_NAME build settings not found")
		require.Equal(t, "", resolved)
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
	project, err := Open(filepath.Join(dir, "XcodeProj.xcodeproj"))
	require.NoError(t, err)

	{
		scheme, ok := project.Scheme("ProjectTodayExtensionScheme")
		require.True(t, ok)
		require.Equal(t, "ProjectTodayExtensionScheme", scheme.Name)
	}

	{
		scheme, ok := project.Scheme("NotExistScheme")
		require.False(t, ok)
		require.Equal(t, "", scheme.Name)
	}

	{
		// Gdańsk represented in High Sierra
		b := []byte{71, 100, 97, 197, 132, 115, 107}
		scheme, ok := project.Scheme(string(b))
		require.True(t, ok)
		require.Equal(t, norm.NFC.String(string(b)), norm.NFC.String(scheme.Name))
	}

	{
		// Gdańsk represented in Mojave
		b := []byte{71, 100, 97, 110, 204, 129, 115, 107}
		scheme, ok := project.Scheme(string(b))
		require.True(t, ok)
		require.Equal(t, norm.NFC.String(string(b)), norm.NFC.String(scheme.Name))
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
	dir := testhelper.GitCloneIntoTmpDir(t, "https://github.com/bitrise-io/xcode-project-test.git")
	project, err := Open(filepath.Join(dir, "XcodeProj.xcodeproj"))
	require.NoError(t, err)
	require.Equal(t, filepath.Join(dir, "XcodeProj.xcodeproj"), project.Path)
	require.Equal(t, "XcodeProj", project.Name)
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
		developmentTeam         string
		targetName              string
		codeSignIdentity        string
		provisioningProfileUUID string
		wantErr                 bool
	}{
		{
			name:                    "Force code sign - XcodeProj",
			developmentTeam:         "72SA8V3WYL",
			targetName:              "XcodeProj",
			codeSignIdentity:        "",
			provisioningProfileUUID: "",
			wantErr:                 false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := project.ForceCodeSign(tt.targetName, tt.developmentTeam, tt.codeSignIdentity, tt.provisioningProfileUUID); (err != nil) != tt.wantErr {
				t.Errorf("XcodeProj.forceCodeSign() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}

	if err := project.Save(); err != nil {
		t.Errorf("Failed to save project, error: %s", err)
	}
	t.Errorf("Project path: %s", project.Path)
}

func TestXcodeProj_foreceCodeSignOnTargetAttributes(t *testing.T) {
	dir := testhelper.GitCloneIntoTmpDir(t, "https://github.com/bitrise-io/xcode-project-test.git")
	project, err := Open(filepath.Join(dir, "XcodeProj.xcodeproj"))
	if err != nil {
		t.Fatalf("Failed to init project for test case, error: %s", err)
	}
	tests := []struct {
		name                    string
		targetAttributes        serialized.Object
		developmentTeam         string
		targetID                string
		provisioningProfileUUID string
		want                    serialized.Object
		wantErr                 bool
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
			developmentTeam:         "72SA8V3WYL",
			targetID:                "7D5B35FB20E28EE80022BAE6",
			provisioningProfileUUID: "",
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
			got, err := foreceCodeSignOnTargetAttributes(tt.targetAttributes, tt.targetID, tt.developmentTeam, tt.provisioningProfileUUID)
			if (err != nil) != tt.wantErr {
				t.Errorf("XcodeProj.foreceCodeSignOnTargetAttributes() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("XcodeProj.foreceCodeSignOnTargetAttributes() got = %s, wantErr %s", pretty.Object(got), pretty.Object(tt.want))
				return
			}
		})
	}
}

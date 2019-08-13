package xcodeproj

import (
	"path/filepath"
	"reflect"
	"testing"

	"github.com/bitrise-io/xcode-project/pretty"
	"github.com/bitrise-io/xcode-project/serialized"
	"github.com/bitrise-io/xcode-project/testhelper"
	"github.com/stretchr/testify/require"
	"howett.net/plist"
)

func TestParseConfigurationList(t *testing.T) {
	var raw serialized.Object
	_, err := plist.Unmarshal([]byte(rawConfigurationList), &raw)
	require.NoError(t, err)

	configurationList, err := parseConfigurationList("13E76E3A1F4AC90A0028096E", raw)
	require.NoError(t, err)
	// fmt.Printf("configurationList:\n%s\n", pretty.Object(configurationList))
	require.Equal(t, expectedConfigurationList, pretty.Object(configurationList))
}

const rawConfigurationList = `
{
	13E76E3A1F4AC90A0028096E /* Build configuration list for PBXNativeTarget "code-sign-test" */ = {
		isa = XCConfigurationList;
		buildConfigurations = (
			13E76E3B1F4AC90A0028096E /* Debug */,
			13E76E3C1F4AC90A0028096E /* Release */,
		);
		defaultConfigurationIsVisible = 0;
		defaultConfigurationName = Release;
	};

	13E76E3B1F4AC90A0028096E /* Debug */ = {
		isa = XCBuildConfiguration;
		buildSettings = {
			ASSETCATALOG_COMPILER_APPICON_NAME = AppIcon;
			"CODE_SIGN_IDENTITY[sdk=iphoneos*]" = "iPhone Developer";
			CODE_SIGN_STYLE = Automatic;
			DEVELOPMENT_TEAM = 72SA8V3WYL;
			INFOPLIST_FILE = "code-sign-test/Info.plist";
			LD_RUNPATH_SEARCH_PATHS = "$(inherited) @executable_path/Frameworks";
			PRODUCT_BUNDLE_IDENTIFIER = "com.bitrise.code-sign-test";
			PRODUCT_NAME = "$(TARGET_NAME)";
			PROVISIONING_PROFILE = "";
			PROVISIONING_PROFILE_SPECIFIER = "";
			TARGETED_DEVICE_FAMILY = "1,2";
		};
		name = Debug;
	};

	13E76E3C1F4AC90A0028096E /* Release */ = {
		isa = XCBuildConfiguration;
		buildSettings = {
			ASSETCATALOG_COMPILER_APPICON_NAME = AppIcon;
			"CODE_SIGN_IDENTITY[sdk=iphoneos*]" = "iPhone Developer";
			CODE_SIGN_STYLE = Automatic;
			DEVELOPMENT_TEAM = 72SA8V3WYL;
			INFOPLIST_FILE = "code-sign-test/Info.plist";
			LD_RUNPATH_SEARCH_PATHS = "$(inherited) @executable_path/Frameworks";
			PRODUCT_BUNDLE_IDENTIFIER = "com.bitrise.code-sign-test";
			PRODUCT_NAME = "$(TARGET_NAME)";
			PROVISIONING_PROFILE = "";
			PROVISIONING_PROFILE_SPECIFIER = "";
			TARGETED_DEVICE_FAMILY = "1,2";
		};
		name = Release;
	};
}`

const expectedConfigurationList = `{
	"ID": "13E76E3A1F4AC90A0028096E",
	"DefaultConfigurationName": "Release",
	"BuildConfigurations": [
		{
			"ID": "13E76E3B1F4AC90A0028096E",
			"Name": "Debug",
			"BuildSettings": {
				"ASSETCATALOG_COMPILER_APPICON_NAME": "AppIcon",
				"CODE_SIGN_IDENTITY[sdk=iphoneos*]": "iPhone Developer",
				"CODE_SIGN_STYLE": "Automatic",
				"DEVELOPMENT_TEAM": "72SA8V3WYL",
				"INFOPLIST_FILE": "code-sign-test/Info.plist",
				"LD_RUNPATH_SEARCH_PATHS": "$(inherited) @executable_path/Frameworks",
				"PRODUCT_BUNDLE_IDENTIFIER": "com.bitrise.code-sign-test",
				"PRODUCT_NAME": "$(TARGET_NAME)",
				"PROVISIONING_PROFILE": "",
				"PROVISIONING_PROFILE_SPECIFIER": "",
				"TARGETED_DEVICE_FAMILY": "1,2"
			}
		},
		{
			"ID": "13E76E3C1F4AC90A0028096E",
			"Name": "Release",
			"BuildSettings": {
				"ASSETCATALOG_COMPILER_APPICON_NAME": "AppIcon",
				"CODE_SIGN_IDENTITY[sdk=iphoneos*]": "iPhone Developer",
				"CODE_SIGN_STYLE": "Automatic",
				"DEVELOPMENT_TEAM": "72SA8V3WYL",
				"INFOPLIST_FILE": "code-sign-test/Info.plist",
				"LD_RUNPATH_SEARCH_PATHS": "$(inherited) @executable_path/Frameworks",
				"PRODUCT_BUNDLE_IDENTIFIER": "com.bitrise.code-sign-test",
				"PRODUCT_NAME": "$(TARGET_NAME)",
				"PROVISIONING_PROFILE": "",
				"PROVISIONING_PROFILE_SPECIFIER": "",
				"TARGETED_DEVICE_FAMILY": "1,2"
			}
		}
	]
}`

func TestXcodeProjBuildConfigurationList(t *testing.T) {
	dir := testhelper.GitCloneIntoTmpDir(t, "https://github.com/bitrise-io/xcode-project-test.git")
	project, err := Open(filepath.Join(dir, "XcodeProj.xcodeproj"))
	if err != nil {
		t.Fatalf("Failed to init project for test case, error: %s", err)
	}
	tests := []struct {
		name     string
		targetID string
		want     map[string]interface{}
		wantErr  bool
	}{
		{
			name:     "Fetch xcode-project-test sample's buildConfiguration list",
			targetID: "7D5B35FB20E28EE80022BAE6",
			want: map[string]interface{}{
				"buildConfigurations": []string{
					"7D5B360F20E28EEA0022BAE6",
					"7D5B361020E28EEA0022BAE6",
				},
				"defaultConfigurationIsVisible": "0",
				"defaultConfigurationName":      "Release",
				"isa":                           "XCConfigurationList",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := project.BuildConfigurationList(tt.targetID)
			if (err != nil) != tt.wantErr {
				t.Errorf("XcodeProj.BuildConfigurations() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			slice, _ := got.StringSlice("buildConfigurations")
			if !reflect.DeepEqual(slice, tt.want["buildConfigurations"]) {
				t.Errorf("XcodeProj.BuildConfigurations() = %s, want %s", pretty.Object(got), pretty.Object(tt.want))
			}
		})
	}
}

func TestXcodeProjBuildConfigurations(t *testing.T) {
	dir := testhelper.GitCloneIntoTmpDir(t, "https://github.com/bitrise-io/xcode-project-test.git")
	project, err := Open(filepath.Join(dir, "XcodeProj.xcodeproj"))
	if err != nil {
		t.Fatalf("Failed to init project for test case, error: %s", err)
	}

	buildConfigurationList, err := project.BuildConfigurationList("7D5B35FB20E28EE80022BAE6")
	if err != nil {
		t.Fatalf("Failed to init buildConfigurationList for test case, error: %s", err)
	}

	tests := []struct {
		name                   string
		buildConfigurationList serialized.Object
		want                   map[string]interface{}
		wantErr                bool
	}{
		{
			name:                   "Fetch xcode-project-test sample's buildConfigurations",
			buildConfigurationList: buildConfigurationList,
			want:                   nil,
			wantErr:                false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := project.BuildConfigurations(tt.buildConfigurationList)
			if (err != nil) != tt.wantErr {
				t.Errorf("XcodeProj.BuildConfigurations() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("XcodeProj.BuildConfigurations() = %v, want %v", got, tt.want)
			}
		})
	}
}

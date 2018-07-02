package xcodeproj

import (
	"testing"

	"github.com/bitrise-tools/xcode-project/pretty"
	"github.com/bitrise-tools/xcode-project/serialized"
	"github.com/stretchr/testify/require"
	"howett.net/plist"
)

const rawNativeTarget = `{
	13E76E0D1F4AC90A0028096E /* code-sign-test */ = {
		isa = PBXNativeTarget;
		buildConfigurationList = 13E76E3A1F4AC90A0028096E /* Build configuration list for PBXNativeTarget "code-sign-test" */;
		buildPhases = (
			13E76E0A1F4AC90A0028096E /* Sources */,
			13E76E0B1F4AC90A0028096E /* Frameworks */,
			13E76E0C1F4AC90A0028096E /* Resources */,
			13E76E561F4AC94F0028096E /* Embed App Extensions */,
			13E76E811F4AC9800028096E /* Embed Watch Content */,
		);
		buildRules = (
		);
		dependencies = (
			13E76E511F4AC94F0028096E /* PBXTargetDependency */,
		);
		name = "code-sign-test";
		productName = "code-sign-test";
		productReference = 13E76E0E1F4AC90A0028096E /* code-sign-test.app */;
		productType = "com.apple.product-type.application";
	};

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

	13E76E511F4AC94F0028096E /* PBXTargetDependency */ = {
		isa = PBXTargetDependency;
		target = 13E76E461F4AC94F0028096E /* share-extension */;
		targetProxy = 13E76E501F4AC94F0028096E /* PBXContainerItemProxy */;
	};

	13E76E461F4AC94F0028096E /* share-extension */ = {
		isa = PBXNativeTarget;
		buildConfigurationList = 13E76E3A1F4AC90A0028096E /* Build configuration list for PBXNativeTarget "share-extension" */;
		buildPhases = (
			13E76E431F4AC94F0028096E /* Sources */,
			13E76E441F4AC94F0028096E /* Frameworks */,
			13E76E451F4AC94F0028096E /* Resources */,
		);
		buildRules = (
		);
		dependencies = (
		);
		name = "share-extension";
		productName = "share-extension";
		productReference = 13E76E471F4AC94F0028096E /* share-extension.appex */;
		productType = "com.apple.product-type.app-extension";
	};
}`

const expectedNativeTarget = `{
	"ISA": "PBXNativeTarget",
	"ID": "13E76E0D1F4AC90A0028096E",
	"Name": "code-sign-test",
	"BuildConfigurationList": {
		"ISA": "XCConfigurationList",
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
	},
	"Dependencies": [
		{
			"ISA": "PBXTargetDependency",
			"ID": "13E76E511F4AC94F0028096E",
			"Target": {
				"ISA": "PBXNativeTarget",
				"ID": "13E76E461F4AC94F0028096E",
				"Name": "share-extension",
				"BuildConfigurationList": {
					"ISA": "XCConfigurationList",
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
				},
				"Dependencies": null
			}
		}
	]
}`

func TestParseNativeTarget(t *testing.T) {
	var raw serialized.Object
	_, err := plist.Unmarshal([]byte(rawNativeTarget), &raw)
	require.NoError(t, err)

	nativeTarget, err := parseNativeTarget("13E76E0D1F4AC90A0028096E", raw)
	require.NoError(t, err)
	// fmt.Printf("nativeTarget:\n%s\n", pretty.Object(nativeTarget))
	require.Equal(t, expectedNativeTarget, pretty.Object(nativeTarget))
}

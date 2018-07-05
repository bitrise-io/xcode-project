package xcodeproj

import (
	"testing"

	"github.com/bitrise-tools/xcode-project/pretty"
	"github.com/bitrise-tools/xcode-project/serialized"
	"github.com/stretchr/testify/require"
	"howett.net/plist"
)

const rawProj = `
{
	13E76E061F4AC90A0028096E /* Project object */ = {
		isa = PBXProject;
		attributes = {
			LastUpgradeCheck = 0900;
			ORGANIZATIONNAME = "Gödrei Krisztián";
			TargetAttributes = {
				13E76E0D1F4AC90A0028096E = {
					CreatedOnToolsVersion = 9.0;
					ProvisioningStyle = Automatic;
				};
				13E76E251F4AC90A0028096E = {
					CreatedOnToolsVersion = 9.0;
					ProvisioningStyle = Automatic;
					TestTargetID = 13E76E0D1F4AC90A0028096E;
				};
				13E76E301F4AC90A0028096E = {
					CreatedOnToolsVersion = 9.0;
					ProvisioningStyle = Automatic;
					TestTargetID = 13E76E0D1F4AC90A0028096E;
				};
				13E76E461F4AC94F0028096E = {
					CreatedOnToolsVersion = 9.0;
					ProvisioningStyle = Automatic;
				};
				13E76E581F4AC9800028096E = {
					CreatedOnToolsVersion = 9.0;
					ProvisioningStyle = Automatic;
				};
				13E76E641F4AC9800028096E = {
					CreatedOnToolsVersion = 9.0;
					ProvisioningStyle = Automatic;
				};
			};
		};
		buildConfigurationList = 13E76E091F4AC90A0028096E /* Build configuration list for PBXProject "code-sign-test" */;
		compatibilityVersion = "Xcode 8.0";
		developmentRegion = en;
		hasScannedForEncodings = 0;
		knownRegions = (
			en,
			Base,
		);
		mainGroup = 13E76E051F4AC90A0028096E;
		productRefGroup = 13E76E0F1F4AC90A0028096E /* Products */;
		projectDirPath = "";
		projectRoot = "";
		targets = (
			13E76E301F4AC90A0028096E /* code-sign-testUITests */,
		);
	};

	13E76E301F4AC90A0028096E /* code-sign-testUITests */ = {
		isa = PBXNativeTarget;
		buildConfigurationList = 13E76E401F4AC90A0028096E /* Build configuration list for PBXNativeTarget "code-sign-testUITests" */;
		buildPhases = (
			13E76E2D1F4AC90A0028096E /* Sources */,
			13E76E2E1F4AC90A0028096E /* Frameworks */,
			13E76E2F1F4AC90A0028096E /* Resources */,
		);
		buildRules = (
		);
		dependencies = ();
		name = "code-sign-testUITests";
		productName = "code-sign-testUITests";
		productReference = 13E76E311F4AC90A0028096E /* code-sign-testUITests.xctest */;
		productType = "com.apple.product-type.bundle.ui-testing";
	};

	13E76E401F4AC90A0028096E /* Build configuration list for PBXNativeTarget "code-sign-testUITests" */ = {
		isa = XCConfigurationList;
		buildConfigurations = (
			13E76E411F4AC90A0028096E /* Debug */,
		);
		defaultConfigurationIsVisible = 0;
		defaultConfigurationName = Release;
	};

	13E76E411F4AC90A0028096E /* Debug */ = {
		isa = XCBuildConfiguration;
		buildSettings = {
			"CODE_SIGN_IDENTITY[sdk=iphoneos*]" = "iPhone Developer";
			CODE_SIGN_STYLE = Automatic;
			DEVELOPMENT_TEAM = 72SA8V3WYL;
			INFOPLIST_FILE = "code-sign-testUITests/Info.plist";
			LD_RUNPATH_SEARCH_PATHS = "$(inherited) @executable_path/Frameworks @loader_path/Frameworks";
			PRODUCT_BUNDLE_IDENTIFIER = "com.bitrise.code-sign-testUITests";
			PRODUCT_NAME = "$(TARGET_NAME)";
			PROVISIONING_PROFILE = "";
			PROVISIONING_PROFILE_SPECIFIER = "";
			TARGETED_DEVICE_FAMILY = "1,2";
			TEST_TARGET_NAME = "code-sign-test";
		};
		name = Debug;
	};
}`

const expectedProj = `{
	"ID": "13E76E061F4AC90A0028096E",
	"NativeTargets": [
		{
			"ID": "13E76E301F4AC90A0028096E",
			"Name": "code-sign-testUITests",
			"BuildConfigurationList": {
				"ID": "13E76E401F4AC90A0028096E",
				"DefaultConfigurationName": "Release",
				"BuildConfigurations": [
					{
						"ID": "13E76E411F4AC90A0028096E",
						"Name": "Debug",
						"BuildSettings": {
							"CODE_SIGN_IDENTITY[sdk=iphoneos*]": "iPhone Developer",
							"CODE_SIGN_STYLE": "Automatic",
							"DEVELOPMENT_TEAM": "72SA8V3WYL",
							"INFOPLIST_FILE": "code-sign-testUITests/Info.plist",
							"LD_RUNPATH_SEARCH_PATHS": "$(inherited) @executable_path/Frameworks @loader_path/Frameworks",
							"PRODUCT_BUNDLE_IDENTIFIER": "com.bitrise.code-sign-testUITests",
							"PRODUCT_NAME": "$(TARGET_NAME)",
							"PROVISIONING_PROFILE": "",
							"PROVISIONING_PROFILE_SPECIFIER": "",
							"TARGETED_DEVICE_FAMILY": "1,2",
							"TEST_TARGET_NAME": "code-sign-test"
						}
					}
				]
			},
			"Dependencies": null
		}
	]
}`

func TestParseProj(t *testing.T) {
	var raw serialized.Object
	_, err := plist.Unmarshal([]byte(rawProj), &raw)
	require.NoError(t, err)

	proj, err := parseProj("13E76E061F4AC90A0028096E", raw)
	require.NoError(t, err)
	require.Equal(t, expectedProj, pretty.Object(proj))
}

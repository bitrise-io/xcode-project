package xcodeproj

import (
	"fmt"
	"testing"

	"github.com/bitrise-tools/xcode-project/pretty"
	"github.com/bitrise-tools/xcode-project/serialized"
	"github.com/stretchr/testify/require"
	"howett.net/plist"
)

func TestParseTarget(t *testing.T) {
	t.Log("PBXNativeTarget")
	{
		var raw serialized.Object
		_, err := plist.Unmarshal([]byte(rawNativeTarget), &raw)
		require.NoError(t, err)

		target, err := parseTarget("13E76E0D1F4AC90A0028096E", raw)
		require.NoError(t, err)
		// fmt.Printf("target:\n%s\n", pretty.Object(target))
		require.Equal(t, expectedNativeTarget, pretty.Object(target))
	}

	t.Log("PBXAggregateTarget")
	{
		var raw serialized.Object
		_, err := plist.Unmarshal([]byte(rawAggregateTarget), &raw)
		require.NoError(t, err)

		target, err := parseTarget("FD55DAD914CE0B0000F84D24", raw)
		require.NoError(t, err)
		fmt.Printf("target:\n%s\n", pretty.Object(target))
		require.Equal(t, expectedAggregateTarget, pretty.Object(target))
	}

	t.Log("PBXLegacyTarget")
	{
		var raw serialized.Object
		_, err := plist.Unmarshal([]byte(rawLegacyTarget), &raw)
		require.NoError(t, err)

		target, err := parseTarget("407952600CEA391500E202DC", raw)
		require.NoError(t, err)
		// fmt.Printf("target:\n%s\n", pretty.Object(target))
		require.Equal(t, expectedLegacyTarget, pretty.Object(target))
	}

	t.Log("Invalid Target ID")
	{
		var raw serialized.Object
		_, err := plist.Unmarshal([]byte(rawLegacyTarget), &raw)
		require.NoError(t, err)

		target, err := parseTarget("INVALID_TARGET_ID", raw)
		require.Error(t, err)
		require.Equal(t, Target{}, target)
	}
}

const rawLegacyTarget = `{
	407952600CEA391500E202DC /* build */ = {
		isa = PBXLegacyTarget;
		buildArgumentsString = all;
		buildConfigurationList = 407952610CEA393300E202DC /* Build configuration list for PBXLegacyTarget "build" */;
		buildPhases = (
		);
		buildToolPath = /usr/bin/make;
		buildWorkingDirectory = firmware;
		dependencies = (
		);
		name = build;
		passBuildSettingsInEnvironment = 1;
		productName = "Build All";
	};

	407952610CEA393300E202DC /* Build configuration list for PBXLegacyTarget "build" */ = {
		isa = XCConfigurationList;
		buildConfigurations = (
			407952630CEA393300E202DC /* Release */,
		);
		defaultConfigurationIsVisible = 0;
		defaultConfigurationName = Release;
	};

	407952630CEA393300E202DC /* Release */ = {
		isa = XCBuildConfiguration;
		buildSettings = {
			PATH = "$(PATH):/usr/local/CrossPack-AVR/bin";
		};
		name = Release;
	};
}`

const expectedLegacyTarget = `{
	"Type": "PBXLegacyTarget",
	"ID": "407952600CEA391500E202DC",
	"Name": "build",
	"BuildConfigurationList": {
		"ID": "407952610CEA393300E202DC",
		"DefaultConfigurationName": "Release",
		"BuildConfigurations": [
			{
				"ID": "407952630CEA393300E202DC",
				"Name": "Release",
				"BuildSettings": {
					"PATH": "$(PATH):/usr/local/CrossPack-AVR/bin"
				}
			}
		]
	},
	"Dependencies": null,
	"ProductReference": {
		"Path": ""
	}
}`

const rawAggregateTarget = `{
	FD55DAD914CE0B0000F84D24 /* rpcsvc */ = {
		isa = PBXAggregateTarget;
		buildConfigurationList = FD55DADA14CE0B0000F84D24 /* Build configuration list for PBXAggregateTarget "rpcsvc" */;
		buildPhases = (
			FD55DADC14CE0B0700F84D24 /* Run Script */,
		);
		dependencies = (
		);
		name = rpcsvc;
		productName = rpcsvc;
	};

	FD55DADA14CE0B0000F84D24 /* Build configuration list for PBXAggregateTarget "rpcsvc" */ = {
		isa = XCConfigurationList;
		buildConfigurations = (
			FD55DADB14CE0B0000F84D24 /* Release */,
		);
		defaultConfigurationIsVisible = 0;
		defaultConfigurationName = Release;
	};

	FD55DADB14CE0B0000F84D24 /* Release */ = {
		isa = XCBuildConfiguration;
		buildSettings = {
			INSTALLHDRS_SCRIPT_PHASE = YES;
			PRODUCT_NAME = "$(TARGET_NAME)";
			PRODUCT_BUNDLE_IDENTIFIER = "Bitrise.$(PRODUCT_NAME:rfc1034identifier).watch";
		};
		name = Release;
	};
}`

const expectedAggregateTarget = `{
	"Type": "PBXAggregateTarget",
	"ID": "FD55DAD914CE0B0000F84D24",
	"Name": "rpcsvc",
	"BuildConfigurationList": {
		"ID": "FD55DADA14CE0B0000F84D24",
		"DefaultConfigurationName": "Release",
		"BuildConfigurations": [
			{
				"ID": "FD55DADB14CE0B0000F84D24",
				"Name": "Release",
				"BuildSettings": {
					"INSTALLHDRS_SCRIPT_PHASE": "YES",
					"PRODUCT_BUNDLE_IDENTIFIER": "Bitrise.$(PRODUCT_NAME:rfc1034identifier).watch",
					"PRODUCT_NAME": "$(TARGET_NAME)"
				}
			}
		]
	},
	"Dependencies": null,
	"ProductReference": {
		"Path": ""
	}
}`

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

	13E76E0E1F4AC90A0028096E /* code-sign-test.app */ = {isa = PBXFileReference; explicitFileType = wrapper.application; includeInIndex = 0; path = "code-sign-test.app"; sourceTree = BUILT_PRODUCTS_DIR; };
	13E76E111F4AC90A0028096E /* AppDelegate.h */ = {isa = PBXFileReference; lastKnownFileType = sourcecode.c.h; path = AppDelegate.h; sourceTree = "<group>"; };
	13E76E121F4AC90A0028096E /* AppDelegate.m */ = {isa = PBXFileReference; lastKnownFileType = sourcecode.c.objc; path = AppDelegate.m; sourceTree = "<group>"; };
	13E76E141F4AC90A0028096E /* ViewController.h */ = {isa = PBXFileReference; lastKnownFileType = sourcecode.c.h; path = ViewController.h; sourceTree = "<group>"; };
	13E76E151F4AC90A0028096E /* ViewController.m */ = {isa = PBXFileReference; lastKnownFileType = sourcecode.c.objc; path = ViewController.m; sourceTree = "<group>"; };
	13E76E181F4AC90A0028096E /* Base */ = {isa = PBXFileReference; lastKnownFileType = file.storyboard; name = Base; path = Base.lproj/Main.storyboard; sourceTree = "<group>"; };
	13E76E1A1F4AC90A0028096E /* Assets.xcassets */ = {isa = PBXFileReference; lastKnownFileType = folder.assetcatalog; path = Assets.xcassets; sourceTree = "<group>"; };
	13E76E1D1F4AC90A0028096E /* Base */ = {isa = PBXFileReference; lastKnownFileType = file.storyboard; name = Base; path = Base.lproj/LaunchScreen.storyboard; sourceTree = "<group>"; };
	13E76E1F1F4AC90A0028096E /* Info.plist */ = {isa = PBXFileReference; lastKnownFileType = text.plist.xml; path = Info.plist; sourceTree = "<group>"; };
	13E76E201F4AC90A0028096E /* main.m */ = {isa = PBXFileReference; lastKnownFileType = sourcecode.c.objc; path = main.m; sourceTree = "<group>"; };
	13E76E261F4AC90A0028096E /* code-sign-testTests.xctest */ = {isa = PBXFileReference; explicitFileType = wrapper.cfbundle; includeInIndex = 0; path = "code-sign-testTests.xctest"; sourceTree = BUILT_PRODUCTS_DIR; };
	13E76E2A1F4AC90A0028096E /* code_sign_testTests.m */ = {isa = PBXFileReference; lastKnownFileType = sourcecode.c.objc; path = code_sign_testTests.m; sourceTree = "<group>"; };
	13E76E2C1F4AC90A0028096E /* Info.plist */ = {isa = PBXFileReference; lastKnownFileType = text.plist.xml; path = Info.plist; sourceTree = "<group>"; };
	13E76E311F4AC90A0028096E /* code-sign-testUITests.xctest */ = {isa = PBXFileReference; explicitFileType = wrapper.cfbundle; includeInIndex = 0; path = "code-sign-testUITests.xctest"; sourceTree = BUILT_PRODUCTS_DIR; };
	13E76E351F4AC90A0028096E /* code_sign_testUITests.m */ = {isa = PBXFileReference; lastKnownFileType = sourcecode.c.objc; path = code_sign_testUITests.m; sourceTree = "<group>"; };
	13E76E371F4AC90A0028096E /* Info.plist */ = {isa = PBXFileReference; lastKnownFileType = text.plist.xml; path = Info.plist; sourceTree = "<group>"; };
	13E76E471F4AC94F0028096E /* share-extension.appex */ = {isa = PBXFileReference; explicitFileType = "wrapper.app-extension"; includeInIndex = 0; path = "share-extension.appex"; sourceTree = BUILT_PRODUCTS_DIR; };
	13E76E491F4AC94F0028096E /* ShareViewController.h */ = {isa = PBXFileReference; lastKnownFileType = sourcecode.c.h; path = ShareViewController.h; sourceTree = "<group>"; };
	13E76E4A1F4AC94F0028096E /* ShareViewController.m */ = {isa = PBXFileReference; lastKnownFileType = sourcecode.c.objc; path = ShareViewController.m; sourceTree = "<group>"; };
	13E76E4D1F4AC94F0028096E /* Base */ = {isa = PBXFileReference; lastKnownFileType = file.storyboard; name = Base; path = Base.lproj/MainInterface.storyboard; sourceTree = "<group>"; };
	13E76E4F1F4AC94F0028096E /* Info.plist */ = {isa = PBXFileReference; lastKnownFileType = text.plist.xml; path = Info.plist; sourceTree = "<group>"; };
	13E76E591F4AC9800028096E /* watchkit-app.app */ = {isa = PBXFileReference; explicitFileType = wrapper.application; includeInIndex = 0; path = "watchkit-app.app"; sourceTree = BUILT_PRODUCTS_DIR; };
	13E76E5C1F4AC9800028096E /* Base */ = {isa = PBXFileReference; lastKnownFileType = file.storyboard; name = Base; path = Base.lproj/Interface.storyboard; sourceTree = "<group>"; };
	13E76E5E1F4AC9800028096E /* Assets.xcassets */ = {isa = PBXFileReference; lastKnownFileType = folder.assetcatalog; path = Assets.xcassets; sourceTree = "<group>"; };
	13E76E601F4AC9800028096E /* Info.plist */ = {isa = PBXFileReference; lastKnownFileType = text.plist.xml; path = Info.plist; sourceTree = "<group>"; };
	13E76E651F4AC9800028096E /* watchkit-app Extension.appex */ = {isa = PBXFileReference; explicitFileType = "wrapper.app-extension"; includeInIndex = 0; path = "watchkit-app Extension.appex"; sourceTree = BUILT_PRODUCTS_DIR; };
	13E76E6A1F4AC9800028096E /* InterfaceController.h */ = {isa = PBXFileReference; lastKnownFileType = sourcecode.c.h; path = InterfaceController.h; sourceTree = "<group>"; };
	13E76E6B1F4AC9800028096E /* InterfaceController.m */ = {isa = PBXFileReference; lastKnownFileType = sourcecode.c.objc; path = InterfaceController.m; sourceTree = "<group>"; };
	13E76E6D1F4AC9800028096E /* ExtensionDelegate.h */ = {isa = PBXFileReference; lastKnownFileType = sourcecode.c.h; path = ExtensionDelegate.h; sourceTree = "<group>"; };
	13E76E6E1F4AC9800028096E /* ExtensionDelegate.m */ = {isa = PBXFileReference; lastKnownFileType = sourcecode.c.objc; path = ExtensionDelegate.m; sourceTree = "<group>"; };
	13E76E701F4AC9800028096E /* NotificationController.h */ = {isa = PBXFileReference; lastKnownFileType = sourcecode.c.h; path = NotificationController.h; sourceTree = "<group>"; };
	13E76E711F4AC9800028096E /* NotificationController.m */ = {isa = PBXFileReference; lastKnownFileType = sourcecode.c.objc; path = NotificationController.m; sourceTree = "<group>"; };
	13E76E731F4AC9800028096E /* Assets.xcassets */ = {isa = PBXFileReference; lastKnownFileType = folder.assetcatalog; path = Assets.xcassets; sourceTree = "<group>"; };
	13E76E751F4AC9800028096E /* Info.plist */ = {isa = PBXFileReference; lastKnownFileType = text.plist.xml; path = Info.plist; sourceTree = "<group>"; };
	13E76E761F4AC9800028096E /* PushNotificationPayload.apns */ = {isa = PBXFileReference; lastKnownFileType = text; path = PushNotificationPayload.apns; sourceTree = "<group>"; };
}`

const expectedNativeTarget = `{
	"Type": "PBXNativeTarget",
	"ID": "13E76E0D1F4AC90A0028096E",
	"Name": "code-sign-test",
	"BuildConfigurationList": {
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
			"ID": "13E76E511F4AC94F0028096E",
			"Target": {
				"Type": "PBXNativeTarget",
				"ID": "13E76E461F4AC94F0028096E",
				"Name": "share-extension",
				"BuildConfigurationList": {
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
				"Dependencies": null,
				"ProductReference": {
					"Path": "share-extension.appex"
				}
			}
		}
	],
	"ProductReference": {
		"Path": "code-sign-test.app"
	}
}`

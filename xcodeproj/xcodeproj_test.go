package xcodeproj

import (
	"path/filepath"
	"testing"

	"github.com/bitrise-tools/xcode-project/serialized"
	"github.com/bitrise-tools/xcode-project/testhelper"
	"github.com/stretchr/testify/require"
)

func TestTargets(t *testing.T) {
	dir := testhelper.GitCloneIntoTmpDir(t, "https://github.com/bitrise-samples/xcode-project-test.git")
	project, err := Open(filepath.Join(dir, "Group/SubProject/SubProject.xcodeproj"))
	require.NoError(t, err)

	target, ok := project.Proj.Target("7D0342D720F4B5AD0050B6A6")
	require.True(t, ok)

	dependentTargets := target.DependentTargets()
	require.Equal(t, 2, len(dependentTargets))
	require.Equal(t, "WatchKitApp", dependentTargets[0].Name)
	require.Equal(t, "WatchKitApp Extension", dependentTargets[1].Name)

	{
		properties, err := target.InformationPropertyList("Debug", filepath.Dir(project.Path))
		require.NoError(t, err)
		require.Equal(t, serialized.Object{
			"UISupportedInterfaceOrientations":      []interface{}{"UIInterfaceOrientationPortrait", "UIInterfaceOrientationLandscapeLeft", "UIInterfaceOrientationLandscapeRight"},
			"CFBundleExecutable":                    "$(EXECUTABLE_NAME)",
			"CFBundleInfoDictionaryVersion":         "6.0",
			"CFBundleVersion":                       "1",
			"UIRequiredDeviceCapabilities":          []interface{}{"armv7"},
			"UISupportedInterfaceOrientations~ipad": []interface{}{"UIInterfaceOrientationPortrait", "UIInterfaceOrientationPortraitUpsideDown", "UIInterfaceOrientationLandscapeLeft", "UIInterfaceOrientationLandscapeRight"},
			"CFBundleIdentifier":                    "$(PRODUCT_BUNDLE_IDENTIFIER)",
			"CFBundlePackageType":                   "APPL",
			"LSRequiresIPhoneOS":                    true,
			"CFBundleName":                          "$(PRODUCT_NAME)",
			"CFBundleShortVersionString":            "1.0",
			"CFBundleDevelopmentRegion":             "$(DEVELOPMENT_LANGUAGE)",
			"UILaunchStoryboardName":                "LaunchScreen",
			"UIMainStoryboardFile":                  "Main"}, properties)
	}

	{
		properties, err := target.InformationPropertyList("NotExist", filepath.Dir(project.Path))
		require.EqualError(t, err, "configuration not found with name: NotExist")
		require.Equal(t, serialized.Object(nil), properties)
	}
}

func TestProjectBuildSettings(t *testing.T) {
	dir := testhelper.GitCloneIntoTmpDir(t, "https://github.com/bitrise-samples/xcode-project-test.git")
	project, err := Open(filepath.Join(dir, "XcodeProj.xcodeproj"))
	require.NoError(t, err)

	{
		settings, err := project.Proj.BuildSettings("Debug")
		require.NoError(t, err)
		require.Equal(t, serialized.Object{
			"GCC_WARN_UNINITIALIZED_AUTOS":               "YES_AGGRESSIVE",
			"MTL_ENABLE_DEBUG_INFO":                      "YES",
			"ONLY_ACTIVE_ARCH":                           "YES",
			"SWIFT_OPTIMIZATION_LEVEL":                   "-Onone",
			"CLANG_ANALYZER_NUMBER_OBJECT_CONVERSION":    "YES_AGGRESSIVE",
			"CLANG_WARN_INFINITE_RECURSION":              "YES",
			"CLANG_WARN_RANGE_LOOP_ANALYSIS":             "YES",
			"ALWAYS_SEARCH_USER_PATHS":                   "NO",
			"CLANG_WARN_COMMA":                           "YES",
			"CLANG_WARN_DIRECT_OBJC_ISA_USAGE":           "YES_ERROR",
			"CLANG_WARN_EMPTY_BODY":                      "YES",
			"GCC_WARN_UNDECLARED_SELECTOR":               "YES",
			"GCC_WARN_UNUSED_FUNCTION":                   "YES",
			"SDKROOT":                                    "iphoneos",
			"CLANG_CXX_LIBRARY":                          "libc++",
			"CLANG_ENABLE_OBJC_ARC":                      "YES",
			"CLANG_WARN_CONSTANT_CONVERSION":             "YES",
			"CLANG_WARN_OBJC_ROOT_CLASS":                 "YES_ERROR",
			"DEBUG_INFORMATION_FORMAT":                   "dwarf",
			"GCC_DYNAMIC_NO_PIC":                         "NO",
			"GCC_OPTIMIZATION_LEVEL":                     "0",
			"GCC_PREPROCESSOR_DEFINITIONS":               []interface{}{"DEBUG=1", "$(inherited)"},
			"CLANG_ENABLE_MODULES":                       "YES",
			"CLANG_WARN_BOOL_CONVERSION":                 "YES",
			"GCC_WARN_ABOUT_RETURN_TYPE":                 "YES_ERROR",
			"CLANG_WARN_DOCUMENTATION_COMMENTS":          "YES",
			"CLANG_WARN_OBJC_IMPLICIT_RETAIN_SELF":       "YES",
			"CLANG_WARN_STRICT_PROTOTYPES":               "YES",
			"CLANG_WARN_UNREACHABLE_CODE":                "YES",
			"COPY_PHASE_STRIP":                           "NO",
			"SWIFT_ACTIVE_COMPILATION_CONDITIONS":        "DEBUG",
			"CLANG_WARN_BLOCK_CAPTURE_AUTORELEASING":     "YES",
			"CLANG_WARN_DEPRECATED_OBJC_IMPLEMENTATIONS": "YES",
			"CLANG_WARN_UNGUARDED_AVAILABILITY":          "YES_AGGRESSIVE",
			"CLANG_WARN__DUPLICATE_METHOD_MATCH":         "YES",
			"GCC_C_LANGUAGE_STANDARD":                    "gnu11",
			"GCC_WARN_UNUSED_VARIABLE":                   "YES",
			"CLANG_CXX_LANGUAGE_STANDARD":                "gnu++14",
			"CLANG_WARN_ENUM_CONVERSION":                 "YES",
			"CLANG_WARN_NON_LITERAL_NULL_CONVERSION":     "YES",
			"GCC_NO_COMMON_BLOCKS":                       "YES",
			"IPHONEOS_DEPLOYMENT_TARGET":                 "11.4",
			"CLANG_ENABLE_OBJC_WEAK":                     "YES",
			"CLANG_WARN_INT_CONVERSION":                  "YES",
			"CLANG_WARN_SUSPICIOUS_MOVE":                 "YES",
			"CODE_SIGN_IDENTITY":                         "iPhone Developer",
			"ENABLE_STRICT_OBJC_MSGSEND":                 "YES",
			"ENABLE_TESTABILITY":                         "YES",
			"GCC_WARN_64_TO_32_BIT_CONVERSION":           "YES",
			"CLANG_ANALYZER_NONNULL":                     "YES",
			"CLANG_WARN_OBJC_LITERAL_CONVERSION":         "YES"}, settings)
	}

	{
		settings, err := project.Proj.BuildSettings("NotExist")
		require.EqualError(t, err, "configuration not found with name: NotExist")
		require.Equal(t, serialized.Object(nil), settings)
	}
}

func TestScheme(t *testing.T) {
	dir := testhelper.GitCloneIntoTmpDir(t, "https://github.com/bitrise-samples/xcode-project-test.git")
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
}

func TestSchemes(t *testing.T) {
	dir := testhelper.GitCloneIntoTmpDir(t, "https://github.com/bitrise-samples/xcode-project-test.git")
	project, err := Open(filepath.Join(dir, "XcodeProj.xcodeproj"))
	require.NoError(t, err)

	schemes, err := project.Schemes()
	require.NoError(t, err)
	require.Equal(t, 2, len(schemes))

	require.Equal(t, "ProjectScheme", schemes[0].Name)
	require.Equal(t, "ProjectTodayExtensionScheme", schemes[1].Name)
}

func TestOpenXcodeproj(t *testing.T) {
	dir := testhelper.GitCloneIntoTmpDir(t, "https://github.com/bitrise-samples/xcode-project-test.git")
	project, err := Open(filepath.Join(dir, "XcodeProj.xcodeproj"))
	require.NoError(t, err)
	require.Equal(t, filepath.Join(dir, "XcodeProj.xcodeproj"), project.Path)
	require.Equal(t, "XcodeProj", project.Name)
}

func TestIsXcodeProj(t *testing.T) {
	require.True(t, IsXcodeProj("./BitriseSample.xcodeproj"))
	require.False(t, IsXcodeProj("./BitriseSample.xcworkspace"))
}

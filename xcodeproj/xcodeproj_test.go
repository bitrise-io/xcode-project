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

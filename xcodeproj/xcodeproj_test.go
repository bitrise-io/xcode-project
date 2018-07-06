package xcodeproj

import (
	"path/filepath"
	"testing"

	"github.com/bitrise-tools/xcode-project/testhelper"
	"github.com/stretchr/testify/require"
)

func TestOpenXcodeproj(t *testing.T) {
	dir := testhelper.GitCloneIntoTmpDir(t, "https://github.com/bitrise-samples/sample-apps-ios-simple-objc.git")
	project, err := Open(filepath.Join(dir, "ios-simple-objc/ios-simple-objc.xcodeproj"))
	require.NoError(t, err)
	require.Equal(t, 2, len(project.Proj.Targets))

	{
		target, ok := project.Proj.Target("BA3CBE7419F7A93800CED4D5")
		require.True(t, ok)
		require.Equal(t, "ios-simple-objc", target.Name)
	}

	{
		target, ok := project.Proj.Target("BA3CBE9019F7A93900CED4D5")
		require.True(t, ok)
		require.Equal(t, "ios-simple-objcTests", target.Name)
	}
}

func TestIsXcodeProj(t *testing.T) {
	require.True(t, IsXcodeProj("./BitriseSample.xcodeproj"))
	require.False(t, IsXcodeProj("./BitriseSample.xcworkspace"))
}

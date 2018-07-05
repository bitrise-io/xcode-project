package xcworkspace

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bitrise-tools/xcode-project/testhelper"
)

func TestOpen(t *testing.T) {
	pth := testhelper.CreateTmpFile(t, "contents.xcworkspacedata", workspaceContentsContent)
	dir := filepath.Dir(pth)

	workspace, err := Open(dir)
	require.NoError(t, err)

	require.Equal(t, filepath.Base(filepath.Dir(pth)), workspace.Name)
	require.Equal(t, filepath.Dir(pth), workspace.Path)

	require.Equal(t, 1, len(workspace.Groups))

	{
		require.Equal(t, 1, len(workspace.FileRefs))
		require.Equal(t, "group:XcodeProj.xcodeproj", workspace.FileRefs[0].Location)

		pth, err := workspace.FileRefs[0].AbsPath(dir)
		require.NoError(t, err)
		require.Equal(t, filepath.Join(dir, "XcodeProj.xcodeproj"), pth)
	}

	{
		group := workspace.Groups[0]
		require.Equal(t, 2, len(group.FileRefs))

		require.Equal(t, "group:../XcodeProj/AppDelegate.swift", group.FileRefs[0].Location)
		require.Equal(t, "group:XcodeProj.xcodeproj", group.FileRefs[1].Location)
	}
}

const workspaceContentsContent = `<?xml version="1.0" encoding="UTF-8"?>
<Workspace
   version = "1.0">
   <Group
      location = "group:Group"
      name = "Group">
      <FileRef
         location = "group:../XcodeProj/AppDelegate.swift">
      </FileRef>
      <FileRef
         location = "group:XcodeProj.xcodeproj">
      </FileRef>
   </Group>
   <FileRef
      location = "group:XcodeProj.xcodeproj">
   </FileRef>
</Workspace>
`

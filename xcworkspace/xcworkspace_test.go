package xcworkspace

import (
	"path/filepath"
	"testing"

	"github.com/bitrise-tools/xcode-project/testhelper"
	"github.com/stretchr/testify/require"
)

func TestOpen(t *testing.T) {
	workspaceContentsPth := testhelper.CreateTmpFile(t, "contents.xcworkspacedata", workspaceContentsContent)
	workspacePth := filepath.Dir(workspaceContentsPth)

	workspace, err := Open(workspacePth)
	require.NoError(t, err)

	require.Equal(t, filepath.Base(workspacePth), workspace.Name)
	require.Equal(t, workspacePth, workspace.Path)
	require.Equal(t, 1, len(workspace.FileRefs))
	require.Equal(t, "group:XcodeProj.xcodeproj", workspace.FileRefs[0].Location)
	require.Equal(t, 1, len(workspace.Groups))
	require.Equal(t, "group:Group", workspace.Groups[0].Location)
}

func TestIsWorkspace(t *testing.T) {
	require.True(t, IsWorkspace("./BitriseSample.xcworkspace"))
	require.False(t, IsWorkspace("./BitriseSample.xcodeproj"))
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

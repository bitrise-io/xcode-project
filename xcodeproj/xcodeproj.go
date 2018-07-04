package xcodeproj

import (
	"fmt"
	"path/filepath"

	"github.com/bitrise-io/go-utils/fileutil"
	"github.com/bitrise-tools/xcode-project/xcodeproj/pbxproj"
	"howett.net/plist"
)

// XcodeProj ...
type XcodeProj struct {
	raw map[string]interface{}

	ArchiveVersion string
	Classes        interface{}
	ObjectVersion  string

	Proj       pbxproj.PBXProj
	SharedData interface{}
	UserData   interface{}
	Workspace  XCWorkspace
}

// Open ...
func Open(pth string) (XcodeProj, error) {
	pbxProjPth := filepath.Join(pth, "project.pbxproj")

	b, err := fileutil.ReadBytesFromFile(pbxProjPth)
	if err != nil {
		return XcodeProj{}, err
	}

	var raw map[string]interface{}
	if _, err := plist.Unmarshal(b, &raw); err != nil {
		return XcodeProj{}, fmt.Errorf("failed to generate json from Pbxproj - error: %s", err)
	}

	objects := raw["objects"].(map[string]interface{})

	return XcodeProj{
		ArchiveVersion: raw["archiveVersion"].(string),
		Classes:        raw["classes"],
		ObjectVersion:  raw["objectVersion"].(string),
		Proj:           pbxproj.GetPBXProj(pbxProjID(objects), objects),
	}, nil
}

func pbxProjID(raw map[string]interface{}) string {
	for id, obj := range raw {
		rawObj := obj.(map[string]interface{})
		isa := rawObj["isa"].(string)

		if isa == "PBXProject" {
			return id
		}
	}
	return ""
}

package xcodeproj

import (
	"reflect"
	"testing"

	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/bitrise-io/xcode-project/serialized"
	"howett.net/plist"
)

func Test_parseResourcesBuildPhase(t *testing.T) {
	var raw serialized.Object
	_, err := plist.Unmarshal([]byte(rawResourcesBuildPhase), &raw)
	if err != nil {
		t.Errorf("setup: failed to parse raw object")
	}

	const id1 = "47C11A3D21FF63950084FD7F"

	tests := []struct {
		name string

		id      string
		objects serialized.Object

		want    resourcesBuildPhase
		wantErr bool
	}{
		{
			name:    "normal",
			id:      id1,
			objects: raw,
			want: resourcesBuildPhase{
				ID:    id1,
				files: []string{"47C11A4D21FF63970084FD7F", "47C11A4A21FF63970084FD7F", "47C11A4821FF63950084FD7F"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseResourcesBuildPhase(tt.id, tt.objects)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseResourcesBuildPhase() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseResourcesBuildPhase() = %v, want %v", got, tt.want)
			}
		})
	}
}

const rawResourcesBuildPhase = `
/* Begin PBXResourcesBuildPhase section */
		47C11A3D21FF63950084FD7F /* Resources */ = {
			isa = PBXResourcesBuildPhase;
			buildActionMask = 2147483647;
			files = (
				47C11A4D21FF63970084FD7F /* LaunchScreen.storyboard in Resources */,
				47C11A4A21FF63970084FD7F /* Assets.xcassets in Resources */,
				47C11A4821FF63950084FD7F /* Main.storyboard in Resources */,
			);
			runOnlyForDeploymentPostprocessing = 0;
		};
		47F01785221C4C1E00DF0B8B /* Resources */ = {
			isa = PBXResourcesBuildPhase;
			buildActionMask = 2147483647;
			files = (
				47F0178F221C4C1E00DF0B8B /* MainInterface.storyboard in Resources */,
			);
			runOnlyForDeploymentPostprocessing = 0;
		};
/* End PBXResourcesBuildPhase section */
`

func Test_parseFileReference(t *testing.T) {
	var raw serialized.Object
	_, err := plist.Unmarshal([]byte(rawFileReference), &raw)
	if err != nil {
		t.Errorf("setup: failed to parse raw object")
	}

	tests := []struct {
		name string

		id      string
		objects serialized.Object

		want    fileReference
		wantErr bool
	}{
		{
			name:    "Normal case",
			id:      "47C11A4921FF63970084FD7F",
			objects: raw,
			want: fileReference{
				id:         "47C11A4921FF63970084FD7F",
				path:       "Assets.xcassets",
				sourceTree: groupParent,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseFileReference(tt.id, tt.objects)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseFileReference() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseFileReference() = %v, want %v", got, tt.want)
			}
		})
	}
}

const rawFileReference = `
47C11A4921FF63970084FD7F /* Assets.xcassets */ = {isa = PBXFileReference; lastKnownFileType = folder.assetcatalog; path = Assets.xcassets; sourceTree = "<group>"; };
`

func Test_resolveFileReferenceAbsolutePath(t *testing.T) {
	const projectPth = "/Users/lpusok/Develop/_ios_github/OnionBrowser-2.X/OnionBrowser2.xcodeproj"
	absPth, err := pathutil.AbsPath(projectPth)
	if err != nil {
		t.Errorf("abs")
	}

	objects, projectID, err := open(absPth)

	_, err = parseProj(projectID, objects)
	if err != nil {
		t.Errorf("setup parse proj, error: %s", err)
	}

	type args struct {
		ref       fileReference
		projectID string
		objects   serialized.Object
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "dsf",
			args: args{
				ref: fileReference{
					id: "01801EA51A32CA2A002B4718",
				},
				projectID: projectID,
				objects:   objects,
			},
			want:    "",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := resolveFileReferenceAbsolutePath(tt.args.ref, tt.args.projectID, tt.args.objects)
			if (err != nil) != tt.wantErr {
				t.Errorf("resolveFileReferenceAbsolutePath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("resolveFileReferenceAbsolutePath() = %v, want %v", got, tt.want)
			}
		})
	}
}

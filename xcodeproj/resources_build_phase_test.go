package xcodeproj

import (
	"reflect"
	"testing"

	"github.com/bitrise-tools/xcode-project/serialized"
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

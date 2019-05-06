package xcodeproj

import (
	"io/ioutil"
	"os"
	"path"
	"reflect"
	"testing"

	"github.com/bitrise-io/xcode-project/serialized"
	"howett.net/plist"
)

func Test_assetCatalog(t *testing.T) {
	var objects serialized.Object
	_, err := plist.Unmarshal([]byte(rawProj), &objects)
	if err != nil {
		t.Fatalf("setup: failed to unmarshal project")
	}
	proj, err := parseProj("BA3CBE6D19F7A93800CED4D5", objects)
	if err != nil {
		t.Fatalf("setup: failed to parse project")
	}

	tests := []struct {
		name    string
		target  Target
		objects serialized.Object
		want    []fileReference
		wantErr bool
	}{
		{
			name:    "good path",
			target:  proj.Targets[0],
			objects: objects,
			want: []fileReference{{
				id:   "BA3CBE8819F7A93900CED4D5",
				path: "Images.xcassets",
			}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := assetCatalogs(tt.target, proj.ID, tt.objects)
			if (err != nil) != tt.wantErr {
				t.Errorf("AssetCatalogs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AssetCatalogs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_appIconSetPaths(t *testing.T) {
	projectDir, err := ioutil.TempDir("", "ios-dummy-project")
	if err != nil {
		t.Errorf("setup: failed to create temp dir")
	}
	defer func() {
		if err := os.RemoveAll(projectDir); err != nil {
			t.Logf("Failed to clean up after test, error: %s", err)
		}
	}()

	appIconSetPath := path.Join(projectDir, "ios-simple-objc", "Images.xcassets", "AppIcon.appiconset")
	if err := os.MkdirAll(appIconSetPath, 0755); err != nil {
		t.Errorf("setup: failed top create dir %s", appIconSetPath)
	}

	var objects serialized.Object
	_, err = plist.Unmarshal([]byte(rawProj), &objects)
	if err != nil {
		t.Fatalf("setup: failed to unmarshal project")
	}
	proj, err := parseProj("BA3CBE6D19F7A93800CED4D5", objects)
	if err != nil {
		t.Fatalf("setup: failed to parse project")
	}

	type args struct {
		project     Proj
		projectPath string
		objects     serialized.Object
	}
	tests := []struct {
		name    string
		args    args
		want    TargetsToAppIconSets
		wantErr bool
	}{
		{
			name: "happy case",
			args: args{
				project:     proj,
				projectPath: path.Join(projectDir, "ios-simple-objc.xcodeproj"),
				objects:     objects,
			},
			want: TargetsToAppIconSets{
				proj.Targets[0].ID: []string{appIconSetPath},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := appIconSetPaths(tt.args.project, tt.args.projectPath, tt.args.objects)
			if (err != nil) != tt.wantErr {
				t.Errorf("appIconSetPaths() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("appIconSetPaths() = %v, want %v", got, tt.want)
			}
		})
	}
}

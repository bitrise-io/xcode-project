package xcodeproj

import (
	"io/ioutil"
	"os"
	"path/filepath"
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
	// Create project with 1 asset catalog
	projectDir1, err := ioutil.TempDir("", "ios-dummy-project1")
	if err != nil {
		t.Errorf("setup: failed to create temp dir, %v", err)
	}
	defer func() {
		if err := os.RemoveAll(projectDir1); err != nil {
			t.Logf("Failed to clean up after test, error: %s", err)
		}
	}()

	appIconSetPath1 := filepath.Join(projectDir1, "ios-simple-objc", "Images.xcassets", "AppIcon.appiconset")
	if err := os.MkdirAll(appIconSetPath1, 0755); err != nil {
		t.Errorf("setup: failed top create dir %v", err)
	}

	var objects1 serialized.Object
	_, err = plist.Unmarshal([]byte(rawProj), &objects1)
	if err != nil {
		t.Fatalf("setup: failed to unmarshal project, %v", err)
	}
	// PBXProject object ID
	proj1, err := parseProj("BA3CBE6D19F7A93800CED4D5", objects1)
	if err != nil {
		t.Fatalf("setup: failed to parse project, %v", err)
	}

	// Create dummy project with 2 asset catalogs
	projectDir2, err := ioutil.TempDir("", "ios-dummy-project2")
	if err != nil {
		t.Errorf("setup: failed to create temp dir, %v", err)
	}
	defer func() {
		if err := os.RemoveAll(projectDir2); err != nil {
			t.Logf("Failed to clean up after test, error: %s", err)
		}
	}()

	appIconSetPath2 := filepath.Join(projectDir2, "Catalyst Sample", "Assets.xcassets", "AppIcon.appiconset")
	if err := os.MkdirAll(appIconSetPath2, 0755); err != nil {
		t.Errorf("setup: failed top create dir, %s", err)
	}
	if err := os.MkdirAll(filepath.Join(projectDir2, "Catalyst Sample", "Preview Content", "Preview Assets.appiconset"), 0755); err != nil {
		t.Errorf("setup: failed top create dir, %s", err)
	}

	var objects2 serialized.Object
	_, err = plist.Unmarshal([]byte(rawCatalystProj), &objects2)
	if err != nil {
		t.Fatalf("setup: failed to unmarshal project, %v", err)
	}
	proj2, err := parseProj("13917C0A243F43D00087912B", objects2)
	if err != nil {
		t.Fatalf("setup: failed to parse project, %s", err)
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
			name: "single asset catlog",
			args: args{
				project:     proj1,
				projectPath: filepath.Join(projectDir1, "ios-simple-objc.xcodeproj"),
				objects:     objects1,
			},
			want: TargetsToAppIconSets{
				proj1.Targets[0].ID: []string{appIconSetPath1},
			},
			wantErr: false,
		},
		{
			name: "2 asset catalogs",
			args: args{
				project:     proj2,
				projectPath: filepath.Join(projectDir2, "Catalyst Sample.xcodeproj"),
				objects:     objects2,
			},
			want: TargetsToAppIconSets{
				proj2.Targets[0].ID: []string{appIconSetPath2},
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

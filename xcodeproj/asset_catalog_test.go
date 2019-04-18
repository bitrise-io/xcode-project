package xcodeproj

import (
	"reflect"
	"testing"

	"github.com/bitrise-io/xcode-project/serialized"
	"howett.net/plist"
)

func TestAssetCatalogs(t *testing.T) {
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
		targets []Target
		objects serialized.Object
		want    TargetsToAssetCatalogs
		wantErr bool
	}{
		{
			name:    "good path",
			targets: proj.Targets,
			objects: objects,
			want: TargetsToAssetCatalogs{
				"BA3CBE7419F7A93800CED4D5": []string{"Images.xcassets"},
				"BA3CBE9019F7A93900CED4D5": []string{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := assetCatalogs(tt.targets, tt.objects)
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

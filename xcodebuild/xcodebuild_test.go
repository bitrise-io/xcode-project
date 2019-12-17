package xcodebuild

import (
	"reflect"
	"testing"

	"github.com/bitrise-io/xcode-project/serialized"
)

func Test_parseShowBuildSettingsOutput(t *testing.T) {
	tests := []struct {
		name    string
		out     string
		want    serialized.Object
		wantErr bool
	}{
		{
			name:    "empty output",
			out:     "",
			want:    serialized.Object{},
			wantErr: false,
		},
		{
			name: "simple output",
			out: `    ACTION = build
    AD_HOC_CODE_SIGNING_ALLOWED = NO
    ALTERNATE_GROUP = staff`,
			want:    serialized.Object{"ACTION": "build", "AD_HOC_CODE_SIGNING_ALLOWED": "NO", "ALTERNATE_GROUP": "staff"},
			wantErr: false,
		},
		{
			name: "output header",
			out: `Build settings for action build and target ios-simple-objc:
    ACTION = build
    AD_HOC_CODE_SIGNING_ALLOWED = NO`,
			want:    serialized.Object{"ACTION": "build", "AD_HOC_CODE_SIGNING_ALLOWED": "NO"},
			wantErr: false,
		},
		{
			name:    "Build setting without value",
			out:     `    ACTION = `,
			want:    serialized.Object{"ACTION": ""},
			wantErr: false,
		},
		{
			name:    "Build setting without =",
			out:     `    ACTION `,
			want:    serialized.Object{},
			wantErr: false,
		},
		{
			name:    "Build setting without key",
			out:     `    = `,
			want:    serialized.Object{},
			wantErr: false,
		},
		{
			name:    "Split the first = ",
			out:     `    ACTION = build+=test`,
			want:    serialized.Object{"ACTION": "build+=test"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseShowBuildSettingsOutput(tt.out)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseShowBuildSettingsOutput() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseShowBuildSettingsOutput() = %v, want %v", got, tt.want)
			}
		})
	}
}

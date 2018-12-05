package xcworkspace

import (
	"errors"
	"reflect"
	"testing"
)

func TestNewSchemeNotFoundError(t *testing.T) {
	tests := []struct {
		name      string
		scheme    string
		container string
		want      SchemeNotFoundError
	}{
		{
			name:      "simple test",
			scheme:    "Scheme",
			container: "Workspace",
			want:      SchemeNotFoundError{scheme: "Scheme", container: "Workspace"},
		},
		{
			name:      "empty test",
			scheme:    "",
			container: "",
			want:      SchemeNotFoundError{scheme: "", container: ""},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewSchemeNotFoundError(tt.scheme, tt.container); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSchemeNotFoundError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSchemeNotFoundError_Error(t *testing.T) {
	tests := []struct {
		name      string
		scheme    string
		container string
		want      string
	}{
		{
			name:      "simple test",
			scheme:    "Scheme",
			container: "Workspace",
			want:      "scheme Scheme not found in Workspace",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := SchemeNotFoundError{
				scheme:    tt.scheme,
				container: tt.container,
			}
			if got := e.Error(); got != tt.want {
				t.Errorf("SchemeNotFoundError.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsSchemeNotFoundError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "SchemeNotFoundError",
			err:  NewSchemeNotFoundError("Scheme", "Workspace"),
			want: true,
		},
		{
			name: "not SchemeNotFoundError",
			err:  errors.New("other error"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsSchemeNotFoundError(tt.err); got != tt.want {
				t.Errorf("IsSchemeNotFoundError() = %v, want %v", got, tt.want)
			}
		})
	}
}

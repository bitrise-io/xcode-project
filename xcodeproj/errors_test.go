package xcodeproj

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConfigurationNotFoundError(t *testing.T) {
	err := NewConfigurationNotFoundError("AppStore")
	require.EqualError(t, err, "configuration not found with name: AppStore")
}

package xcodeproj

import (
	"testing"

	"github.com/bitrise-tools/xcode-project/serialized"
	"github.com/stretchr/testify/require"
)

func TestResolve(t *testing.T) {
	{
		bundleID := BundleID("Bitrise.XcodeProjectSample.watch")
		resolved, err := bundleID.Resolve(nil)
		require.NoError(t, err)
		require.Equal(t, "Bitrise.XcodeProjectSample.watch", resolved)
	}

	{
		bundleID := BundleID("Bitrise.$(PRODUCT_NAME:rfc1034identifier).watch")
		buildSettings := serialized.Object(map[string]interface{}{"PRODUCT_NAME": "XcodeProjectSample"})
		resolved, err := bundleID.Resolve(buildSettings)
		require.NoError(t, err)
		require.Equal(t, "Bitrise.XcodeProjectSample.watch", resolved)
	}

	{
		bundleID := BundleID("Bitrise.$PRODUCT_NAME.watch")
		buildSettings := serialized.Object(map[string]interface{}{"PRODUCT_NAME": "XcodeProjectSample"})
		resolved, err := bundleID.Resolve(buildSettings)
		require.NoError(t, err)
		require.Equal(t, "Bitrise.XcodeProjectSample.watch", resolved)
	}

	{
		bundleID := BundleID("Bitrise.$PRODUCT_NAME.watch")
		buildSettings := serialized.Object(map[string]interface{}{"PRODUCT_NAME": "$TARGET_NAME", "TARGET_NAME": "XcodeProjectSample"})
		resolved, err := bundleID.Resolve(buildSettings)
		require.NoError(t, err)
		require.Equal(t, "Bitrise.XcodeProjectSample.watch", resolved)
	}

	{
		bundleID := BundleID("Bitrise.$PRODUCT_NAME.watch")
		buildSettings := serialized.Object(map[string]interface{}{"PRODUCT_NAME": "$TARGET_NAME", "TARGET_NAME": "$PRODUCT_NAME"})
		resolved, err := bundleID.Resolve(buildSettings)
		require.EqualError(t, err, "failed to resolve bundle id (Bitrise.$PRODUCT_NAME.watch): reference cycle found")
		require.Equal(t, "", resolved)
	}
}

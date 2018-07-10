package serialized

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValue(t *testing.T) {
	o := Object(map[string]interface{}{"key": "value"})

	{
		v, err := o.Value("key")
		require.NoError(t, err)
		require.Equal(t, "value", v)
	}

	{
		v, err := o.Value("not_existing_key")
		require.EqualError(t, err, "key (not_existing_key) not found in:\nmap[key:value]")
		require.Equal(t, nil, v)
	}
}

func TestString(t *testing.T) {
	{
		o := Object(map[string]interface{}{"key": "value"})
		v, err := o.String("key")
		require.NoError(t, err)
		require.Equal(t, "value", v)
	}

	o := Object(map[string]interface{}{"key": 0})
	{
		v, err := o.String("key")
		require.EqualError(t, err, "value (0) for key (key) is not a string")
		require.Equal(t, "", v)
	}

	{
		v, err := o.Object("key")
		require.EqualError(t, err, "value (0) for key (key) is not a map[string]interface {}")
		require.Equal(t, Object(nil), v)
	}
}

func TestObject(t *testing.T) {
	o := Object(map[string]interface{}{"key": map[string]interface{}{"object_key": "object_value"}})
	v, err := o.Object("key")
	require.NoError(t, err)
	require.Equal(t, Object(map[string]interface{}{"object_key": "object_value"}), v)
}

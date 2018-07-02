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
		require.EqualError(t, err, "key (not_existing_key) not found")
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

	{
		o := Object(map[string]interface{}{"key": 0})
		v, err := o.String("key")
		require.EqualError(t, err, "value (0) for key (key) is not a string")
		require.Equal(t, "", v)
	}
}

func TestStringSlice(t *testing.T) {
	{
		raw := map[string][]string{"key": []string{"value1", "value2"}}
		o := raw.(Object)
		v, err := o.StringSlice("key")
		require.NoError(t, err)
		require.Equal(t, []string{"value1", "value2"}, v)
	}

	{
		o := Object(map[string]interface{}{"key": 0})
		v, err := o.String("key")
		require.EqualError(t, err, "value (0) for key (key) is not a string")
		require.Equal(t, "", v)
	}
}

/*
// StringSlice ...
func (o Object) StringSlice(key string) ([]string, error) {
	value, err := o.Value(key)
	if err != nil {
		return nil, err
	}

	casted, ok := value.([]interface{})
	if !ok {
		return nil, fmt.Errorf("value (%v) for key (%s) is not an array", value, key)
	}

	slice := []string{}
	for _, v := range casted {
		item, ok := v.(string)
		if !ok {
			return nil, fmt.Errorf("value (%v) for key (%s) is not a string array", casted, key)
		}

		slice = append(slice, item)
	}

	return slice, nil
}

// Object ...
func (o Object) Object(key string) (Object, error) {
	value, err := o.Value(key)
	if err != nil {
		return nil, err
	}

	casted, ok := value.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("value (%v) for key (%s) is not a map string interface", value, key)
	}

	return casted, nil
}

// ObjectSlice ...
func (o Object) ObjectSlice(key string) ([]Object, error) {
	value, err := o.Value(key)
	if err != nil {
		return nil, err
	}

	casted, ok := value.([]interface{})
	if !ok {
		return nil, fmt.Errorf("value (%v) for key (%s) is not an array", value, key)
	}

	slice := []Object{}
	for _, v := range casted {
		item, ok := v.(Object)
		if !ok {
			return nil, fmt.Errorf("value (%v) for key (%s) is not a string array", casted, key)
		}

		slice = append(slice, item)
	}

	return slice, nil
}
*/

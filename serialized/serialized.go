package serialized

import (
	"fmt"
)

// Object ...
type Object map[string]interface{}

// Value ...
func (o Object) Value(key string) (interface{}, error) {
	value, ok := o[key]
	if !ok {
		return nil, fmt.Errorf("key (%s) not found", key)
	}
	return value, nil
}

// String ...
func (o Object) String(key string) (string, error) {
	value, err := o.Value(key)
	if err != nil {
		return "", err
	}

	casted, ok := value.(string)
	if !ok {
		return "", fmt.Errorf("value (%v) for key (%s) is not a string", value, key)
	}

	return casted, nil
}

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

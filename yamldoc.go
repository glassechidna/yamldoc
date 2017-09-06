package yamldoc

import (
	"gopkg.in/yaml.v2"
	"fmt"
)

type YamlDocument struct {
	root yaml.MapSlice
}

func NewYamlDocument(data []byte) (*YamlDocument, error) {
	slice := yaml.MapSlice{}

	err := yaml.Unmarshal(data, &slice)
	if err != nil {
		return nil, err
	}

	return &YamlDocument{root: slice}, nil
}

func (y *YamlDocument) Serialize() ([]byte, error) {
	return yaml.Marshal(y.root)
}

func (y *YamlDocument) Get(path ...interface{}) (interface{}, error) {
	var current interface{} = y.root

	for _, keyIface := range path {
		switch key := keyIface.(type) {
		case string:
			item, _ := itemForKey(current.(yaml.MapSlice), key)
			current = item.Value
		case int:
			arr := current.([]interface{})
			current = arr[key]
		default:
			return nil, error(fmt.Sprintf("unexpected key type: %+v\n", key))
		}
	}

	return current, nil
}

func (y *YamlDocument) Set(value interface{}, path ...interface{}) error {
	last, path := path[len(path)-1], path[:len(path)-1]
	current, err := y.Get(path...)
	if err != nil { return err }

	switch key := last.(type) {
	case string:
		slice := current.(yaml.MapSlice)
		item, idx := itemForKey(slice, key)

		if item != nil {
			item.Value = value
			slice[idx] = *item
		} else {
			item = &yaml.MapItem{Key: key, Value: value}
			slice = append(slice, *item)
			y.Set(slice, path...)
		}

	case int:
		arr := current.([]interface{})
		arr[key] = value
	default:
		return error(fmt.Sprintf("unexpected key type: %+v\n", key))
	}

	return nil
}

func itemForKey(s yaml.MapSlice, key interface{}) (*yaml.MapItem, int) {
	for idx, pair := range s {
		if pair.Key == key {
			return &pair, idx
		}
	}

	return nil, 0
}

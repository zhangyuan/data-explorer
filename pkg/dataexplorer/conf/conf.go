package conf

import (
	"os"

	"gopkg.in/yaml.v3"
)

func LoadFromYAML[T any](path string) (*T, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var t T
	if err := yaml.Unmarshal(bytes, &t); err != nil {
		return nil, err
	}
	return &t, nil
}

package types

import (
	jsoniter "github.com/json-iterator/go"
)

type Optional[T any] struct {
	Value  *T
	exists bool
}

func (optionalField *Optional[T]) UnmarshalJSON(data []byte) error {
	var value *T
	if err := jsoniter.Unmarshal(data, &value); err != nil {
		return err
	}
	optionalField.exists = true
	optionalField.Value = value

	return nil
}

func (optionalField *Optional[T]) HasValue() bool {
	return optionalField.exists
}

package utils

import "encoding/json"

// Perform JSON marshalling and unmarshalling
func ConvertStruct(v any, target any) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}
	err = json.Unmarshal(b, target)
	return err
}

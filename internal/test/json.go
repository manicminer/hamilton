package test

import (
	"encoding/json"
	"fmt"
)

func AssertJsonMarshalEquals(value interface{}, expected string) error {
	bytes, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling failed with error %s", err)
	}

	actual := string(bytes)
	if actual != expected {
		return fmt.Errorf("expected marshalled json to equal %s but was %s", expected, actual)
	}

	return nil
}

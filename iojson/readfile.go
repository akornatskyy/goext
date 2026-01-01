// Package iojson provides JSON file I/O utilities.
package iojson

import (
	"encoding/json"
	"os"
)

// ReadFile reads the file named by filename as JSON-encoded
// value and stores it in the value pointed to by v.
func ReadFile(filename string, v interface{}) error {
	value, err := os.ReadFile(filename) // #nosec G304
	if err != nil {
		return err
	}
	return json.Unmarshal(value, v)
}

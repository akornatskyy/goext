package iojson

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReadFile(t *testing.T) {
	// Create a temporary directory for test files
	tmpDir, err := os.MkdirTemp("", "iojson-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	t.Run("ReadValidJSON", func(t *testing.T) {
		// Create a test JSON file
		testFile := filepath.Join(tmpDir, "valid.json")
		testData := `{"name":"John","age":30,"active":true}`
		if err := os.WriteFile(testFile, []byte(testData), 0600); err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}

		// Test reading into a struct
		type Person struct {
			Name   string `json:"name"`
			Age    int    `json:"age"`
			Active bool   `json:"active"`
		}

		var person Person
		err := ReadFile(testFile, &person)
		if err != nil {
			t.Errorf("ReadFile() error = %v, want nil", err)
		}

		if person.Name != "John" {
			t.Errorf("person.Name = %q, want %q", person.Name, "John")
		}
		if person.Age != 30 {
			t.Errorf("person.Age = %d, want %d", person.Age, 30)
		}
		if !person.Active {
			t.Errorf("person.Active = %v, want true", person.Active)
		}
	})
}

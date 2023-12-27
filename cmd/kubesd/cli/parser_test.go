package cli

import (
	"os"
	"testing"
)

func TestReadObject(t *testing.T) {
	data, err := os.OpenFile("../../../testdata/generic/opaque-data.yaml", os.O_RDONLY, 0755)

	if err != nil {
		t.Fatal(err)
	}

	oldStdIn := os.Stdin
	defer func() { os.Stdin = oldStdIn }()

	os.Stdin = data
	s, _ := readObject()

	if s.String() == "" {
		t.Error("No data from STDIN")
	}

}

func TestUnmarshal(t *testing.T) {
	var s SecretYaml
	data, err := os.OpenFile("../../../testdata/generic/opaque-data.yaml", os.O_RDONLY, 0755)

	if err != nil {
		t.Fatal(err)
	}

	oldStdIn := os.Stdin
	defer func() { os.Stdin = oldStdIn }()

	os.Stdin = data
	objectType, _ := s.unmarshal()

	if objectType != "opaque" {
		t.Error("Failed to Unmarshal the object")
	}
}

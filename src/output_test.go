package src

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func TestMarshal(t *testing.T){
	var s SecretYaml = map[string]interface{}{
		"data": map[string]interface{}{
			"PASSWORD": "YWRtaW4=",
			"USERNAME": "YWRtaW4=",
		},
	}

	tmpfile, err := ioutil.TempFile("", "example")
	if err != nil {
		t.Fatal(err)
	}

	defer os.Remove(tmpfile.Name())

    if _, err := tmpfile.Write([]byte(fmt.Sprint(s))); err != nil {
        t.Fatal(err)
    }

    if _, err := tmpfile.Seek(0, 0); err != nil {
        t.Fatal(err)
    }

	oldStdIn := os.Stdin
	defer func() {os.Stdin = oldStdIn}()

	os.Stdin = tmpfile
	object, _ := s.marshal()

	if object == "" {
		t.Error("Decode failed for secret object")
	}
}
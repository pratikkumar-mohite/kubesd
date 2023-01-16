package src

import (
	base64 "encoding/base64"
	"os"
	"strings"
	"testing"
	"io/ioutil"
	"fmt"
)

func TestDecodeBase64(t *testing.T){
	var encodedData = base64.StdEncoding.EncodeToString([]byte("admin"))
	decodedValue, _ := decodeBase64(encodedData)
	if decodedValue != "admin" {
		t.Error("Base64 decode test failed")
	}
}

func TestDecode(t *testing.T){
	var d string = `apiVersion: v1
kind: Secret
metadata:
  name: opaque-secret
  namespace: default
type: Opaque
data:
 USERNAME: cm9vdAo=
 PASSWORD: cm9vdAo=`

	tmpfile, err := ioutil.TempFile("", "example")
	if err != nil {
		t.Fatal(err)
	}

	defer os.Remove(tmpfile.Name())

    if _, err := tmpfile.Write([]byte(fmt.Sprint(d))); err != nil {
        t.Fatal(err)
    }

    if _, err := tmpfile.Seek(0, 0); err != nil {
        t.Fatal(err)
    }

	data, _ := os.OpenFile(tmpfile.Name(), os.O_RDONLY, 0755)

	oldStdIn := os.Stdin
	defer func() {os.Stdin = oldStdIn}()

	os.Stdin = data
	objectType, err := Decode()
	if err != nil {
		t.Fatal(err)
	}

	if strings.Count(objectType,"root") != 2 {
		t.Error("Decode failed for secret object")
	}
}

func TestDecodeMultiple(t *testing.T){
	var testDataPath = "../testdata/"
	items, _ := ioutil.ReadDir(testDataPath)
    for _, item := range items {
        if item.IsDir() {
            subitems, _ := ioutil.ReadDir(testDataPath+item.Name())
            for _, subitem := range subitems {
                if !subitem.IsDir() {
					path := testDataPath+item.Name() + "/" + subitem.Name()
					data, err := os.OpenFile(path, os.O_RDONLY, 0755)

					if err != nil {
						t.Fatal(err)
					}
				
					oldStdIn := os.Stdin
					defer func() {os.Stdin = oldStdIn}()
				
					os.Stdin = data
					var secretObject, _ = Decode()
				
					if strings.Contains(secretObject,"!!binary") {
						t.Error("Decode failed for secret object : "+path)
					}
				}
            }
        }
    }
}
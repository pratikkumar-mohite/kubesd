package src

import (
	"fmt"
	base64 "encoding/base64"
)

// Decode the data to secret object
func decodeBase64(value string) string {
	keyval, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		fmt.Printf("Failed to decode object %v",err)
	}
	return string(keyval)
}

func decodeOpaque(y *secretData) {
	for key, value := range y.Data {
		y.Data[key] = decodeBase64(value)
	}
}

func Decode() {
	yamlContent, objectType := readObject()
	if objectType == "opaque" {
		decodeOpaque(yamlContent)
	}
	fmt.Println(yamlContent.Data)
}
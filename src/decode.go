package src

import (
	"fmt"
	base64 "encoding/base64"
)

// Decode the data to secret object
func decodeBase64(value string) string {
	decodedData, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		fmt.Printf("Failed to decode object %v",err)
	}
	return string(decodedData)
}

// Opaque : arbitrary user-defined data
func decodeOpaque() {
	for key, value := range Data {
		Data[key] = decodeBase64(value)
	}
}

func Decode() {
	secretCompleteObject, objectType := readObject()
	if objectType == "opaque" {
		decodeOpaque()
	}
	printSecretObject(secretCompleteObject)
}
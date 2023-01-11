package src

import (
	base64 "encoding/base64"
	"fmt"
)

var supportedObjectTypes = []string{"opaque","kubernetes.io/dockerconfigjson","kubernetes.io/dockercfg","kubernetes.io/tls"}

// Decode the data to secret object
func decodeBase64(value string) string {
	decodedData, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		fmt.Printf("Failed to decode object %v", err)
	}
	return string(decodedData)
}

// Opaque : arbitrary user-defined data
func decodeOpaque() {
	for key, value := range Data {
		Data[key] = decodeBase64(value)
	}
}

func contains(key string, list []string) bool {
	for _, value := range list {
		if value == key {
			return true
		}
	}
	return false
}

func Decode() {
	secretCompleteObject, objectType := readObject()
	if !isStringData {
		if contains(objectType, supportedObjectTypes) {
			decodeOpaque()
		} else {
			fmt.Printf("Invalid secret object type : %v", objectType)
			return
		}
	}
	printSecretObject(secretCompleteObject)
}

package src

import (
	base64 "encoding/base64"
	"fmt"
)

// Secret types: https://kubernetes.io/docs/concepts/configuration/secret/#secret-types
// Opaque	arbitrary user-defined data
// kubernetes.io/service-account-token	ServiceAccount token
// kubernetes.io/dockercfg	serialized ~/.dockercfg file
// kubernetes.io/dockerconfigjson	serialized ~/.docker/config.json file
// kubernetes.io/basic-auth	credentials for basic authentication
// kubernetes.io/ssh-auth	credentials for SSH authentication
// kubernetes.io/tls	data for a TLS client or server
// bootstrap.kubernetes.io/token	bootstrap token data

type SecretYaml map[string]interface{}

var supportedObjectTypes = []string{"opaque","kubernetes.io/service-account-token","kubernetes.io/dockercfg","kubernetes.io/dockerconfigjson","kubernetes.io/basic-auth","kubernetes.io/ssh-auth","kubernetes.io/tls","bootstrap.kubernetes.io/token"}

// Decode the data to secret object
func decodeBase64(value string) string {
	decodedData, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		fmt.Printf("Failed to decode object %v", err)
	}
	return string(decodedData)
}

func decodeData() {
	for key, value := range Data {
		Data[key] = decodeBase64(value)
	}
}

func doesListContains(key string, list []string) bool {
	for _, value := range list {
		if value == key {
			return true
		}
	}
	return false
}

func Decode() {
	var s SecretYaml
	var objectType = s.unmarshal()

	if !isStringData {
		if doesListContains(objectType, supportedObjectTypes) {
			decodeData()
		} else {
			fmt.Printf("Invalid secret object type : %v", objectType)
			return
		}
	}
	fmt.Println(s.marshal())
}

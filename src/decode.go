package src

import (
	base64 "encoding/base64"
	"errors"
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
func decodeBase64(value string) (string, error) {
	decodedData, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		return "", err
	}
	return string(decodedData), nil
}

func decodeData() (error){
	var err error
	for key, value := range Data {
		Data[key], err = decodeBase64(value)
		if err != nil{
			return err
		}
	}
	return nil
}

func doesListContains(key string, list []string) (bool) {
	for _, value := range list {
		if value == key {
			return true
		}
	}
	return false
}

func Decode() (string, error) {
	var s SecretYaml
	var objectType, err = s.unmarshal()

	if err != nil {
		return "", err
	}

	if !isStringData {
		if doesListContains(objectType, supportedObjectTypes) {
			decodeData()
		} else {
			return "", errors.New("Invalid secret object type : "+objectType)
		}
	}

	marshalData, err := s.marshal()
	if err != nil {
		return "", err
	}
	return marshalData, nil
}

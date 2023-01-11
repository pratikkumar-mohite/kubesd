package src

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	json "encoding/json"
	yaml "gopkg.in/yaml.v2"
)

type SecretYaml map[string]interface{}

var Kind, Type, OutputType string
var Data map[string]string

var isStringData = false 

// Read secret object from stdin
func readObject() (sComplete SecretYaml, objectType string) {
	var secretCompleteObject SecretYaml
	var secretString strings.Builder
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Err() != nil {
		fmt.Printf("Failed to read from STDIN %v\n", scanner.Err())
		return
	}
	for scanner.Scan() {
		secretString.WriteString(scanner.Text() + "\n")
	}

	// check if data is json
	OutputType = isJson(secretString.String())

	if OutputType == "json" {
		err := json.Unmarshal([]byte(secretString.String()), &secretCompleteObject)
		if err != nil {
			fmt.Printf("Failed to decode object %v\n", err)
		}

	} else {
		err := yaml.Unmarshal([]byte(secretString.String()), &secretCompleteObject)
		if err != nil {
			fmt.Printf("Failed to decode object %v\n", err)
		}
	}

	// set secret objects
	Kind = secretCompleteObject["kind"].(string)
	Type = secretCompleteObject["type"].(string)
	var data = secretCompleteObject["data"]
	var stringData = secretCompleteObject["stringData"]

	if data == nil && stringData == nil{
		fmt.Println("No valid data field found")
		return
	} else if data == nil {
		isStringData = true
		data = stringData
	}

	Data = make(map[string]string)

	if OutputType == "json" {
		convertJsonInterfaceObject(data.(map[string]interface{}))
	} else {
		convertYamlInterfaceObject(data.(map[interface{}]interface{}))
	}

	verifyObject(strings.ToLower(Kind), "secret")

	return secretCompleteObject, strings.ToLower(Type)
}

func convertYamlInterfaceObject(data map[interface{}]interface{}) {
	for key, value := range data {
		strKey := fmt.Sprintf("%v", key)
		strValue := fmt.Sprintf("%v", value)
		Data[strKey] = strValue
	}
}

func convertJsonInterfaceObject(data map[string]interface{}) {
	for key, value := range data {
		strValue := fmt.Sprintf("%v", value)
		Data[key] = strValue
	}
}

// Verify object
func verifyObject(actual string, expected string) {
	if actual != expected {
		fmt.Println("The given object is not a Secret object")
		return
	}
}

func isJson(s string) string {
	var js interface{}
	if json.Unmarshal([]byte(s), &js) == nil {
		return "json"
	}
	return "nil"
}

// Identify the secret object from following types:
// Opaque	arbitrary user-defined data
// kubernetes.io/service-account-token	ServiceAccount token
// kubernetes.io/dockercfg	serialized ~/.dockercfg file
// kubernetes.io/dockerconfigjson	serialized ~/.docker/config.json file
// kubernetes.io/basic-auth	credentials for basic authentication
// kubernetes.io/ssh-auth	credentials for SSH authentication
// kubernetes.io/tls	data for a TLS client or server
// bootstrap.kubernetes.io/token	bootstrap token data
// func identifySecretObjectType(objectType string) string {
//
// }

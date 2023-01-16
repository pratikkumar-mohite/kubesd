package src

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"errors"

	json "encoding/json"
	yaml "gopkg.in/yaml.v2"
)

var Kind, Type, OutputType string
var Data map[string]string

var isStringData = false 

// Read secret object from stdin
func readObject() (strings.Builder, error) {
	var secretString strings.Builder
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Err() != nil {
		return strings.Builder{}, scanner.Err()
	}
	for scanner.Scan() {
		secretString.WriteString(scanner.Text() + "\n")
	}
	return secretString, nil
}

func (sComplete *SecretYaml)unmarshal() (string, error){
	var secretString, err = readObject()

	if err != nil {
		return "", err
	}
	
	// check if data is json
	OutputType = isJson(secretString.String())

	if OutputType == "json" {
		err := json.Unmarshal([]byte(secretString.String()), &sComplete)
		if err != nil {
			return "", err
		}

	} else {
		err := yaml.Unmarshal([]byte(secretString.String()), &sComplete)
		if err != nil {
			return "", err
		}
	}

	// set secret objects - Kind, Type, Data
	Kind = (*sComplete)["kind"].(string)
	err = verifyObjectKind(strings.ToLower(Kind), "secret")
	if err != nil {
		return "", err
	}

	Type = (*sComplete)["type"].(string)
	
	var data = (*sComplete)["data"]
	var stringData = (*sComplete)["stringData"]

	if data == nil && stringData == nil{
		return "", errors.New("the data/stringData field not found")
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

	return strings.ToLower(Type), nil
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

func verifyObjectKind[T comparable](actual T, expected T) (error){
	if actual != expected {
		return errors.New("the given object is not a secret object")
	}
	return nil
}

func isJson(s string) (string) {
	var js interface{}
	if json.Unmarshal([]byte(s), &js) == nil {
		return "json"
	}
	return "nil"
}

package src

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	json "encoding/json"
	yaml "gopkg.in/yaml.v2"
)

var Kind, Type, OutputType string
var Data map[string]string

var isStringData = false 

// Read secret object from stdin
func readObject() (s strings.Builder) {
	var secretString strings.Builder
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Err() != nil {
		fmt.Printf("Failed to read from STDIN %v\n", scanner.Err())
		return
	}
	for scanner.Scan() {
		secretString.WriteString(scanner.Text() + "\n")
	}
	return secretString
}

func (sComplete *SecretYaml)unmarshal() (objectType string){
	var secretString = readObject()
	
	// check if data is json
	OutputType = isJson(secretString.String())

	if OutputType == "json" {
		err := json.Unmarshal([]byte(secretString.String()), &sComplete)
		if err != nil {
			fmt.Printf("Failed to decode object %v\n", err)
		}

	} else {
		err := yaml.Unmarshal([]byte(secretString.String()), &sComplete)
		if err != nil {
			fmt.Printf("Failed to decode object %v\n", err)
		}
	}

	// set secret objects - Kind, Type, Data
	Kind = (*sComplete)["kind"].(string)
	verifyObjectKind(strings.ToLower(Kind), "secret")

	Type = (*sComplete)["type"].(string)
	
	var data = (*sComplete)["data"]
	var stringData = (*sComplete)["stringData"]

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

	return strings.ToLower(Type)
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

func verifyObjectKind[T comparable](actual T, expected T) {
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

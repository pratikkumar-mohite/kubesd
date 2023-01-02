package src

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	yaml "gopkg.in/yaml.v2"
	json "encoding/json"
)

type secretData struct {
	Kind string `yaml:"kind" json:"kind"`
	Type string `yaml:"type" json:"type"`
	Data map[string]string `yaml:"data" json:"data"`
}

type SecretYaml map[string]interface{}

var OutputType string

// Read secret object from stdin
func readObject() (yData *secretData, yComplete SecretYaml, objectType string){
	var secretDataObject secretData
	var secretCompleteObject SecretYaml
	var secretString strings.Builder
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Err() != nil {
		fmt.Printf("Failed to read from STDIN %v\n",scanner.Err())
		return
	}
	for scanner.Scan() {
		secretString.WriteString(scanner.Text()+"\n")
	}

	// check if data is json
	OutputType = isJson(secretString.String())

	if OutputType == "json" {
		// only data part
		err := json.Unmarshal([]byte(secretString.String()),&secretDataObject)
		if err != nil {
			fmt.Printf("Failed to decode object %v\n",err)
		}
		// complete json object
		err = json.Unmarshal([]byte(secretString.String()),&secretCompleteObject)
		if err != nil {
			fmt.Printf("Failed to decode object %v\n",err)
		}

	} else {
		// only data part
		err := yaml.Unmarshal([]byte(secretString.String()),&secretDataObject)
		if err != nil {
			fmt.Printf("Failed to decode object %v\n",err)
		}
		// complete yaml object
		err = yaml.Unmarshal([]byte(secretString.String()),&secretCompleteObject)
		if err != nil {
			fmt.Printf("Failed to decode object %v\n",err)
		}
	}

	if !isSecretObject(&secretDataObject) {
		fmt.Println("The given object is not a Secret object")
		return
	}
	return &secretDataObject, secretCompleteObject, strings.ToLower(secretDataObject.Type)
}

// Validate the secret object kind
func isSecretObject(y *secretData) bool {
	return strings.ToLower(y.Kind) == "secret"
}

func isJson(s string) string{
	var js interface{}
	if json.Unmarshal([]byte(s),&js) == nil {
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
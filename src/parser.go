package src

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

type secretData struct {
	Kind string `yaml:"kind"`
	Type string `yaml:"type"`
	Data map[string]string `yaml:"data"`
}

type Secret map[string]interface{}

// Read secret object from stdin
func readObject() (yData *secretData, yComplete Secret, objectType string){
	var yamlData secretData
	var yamlComplete Secret
	var secretString strings.Builder
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Err() != nil {
		fmt.Printf("Failed to read from STDIN %v\n",scanner.Err())
		return
	}
	for scanner.Scan() {
		secretString.WriteString(scanner.Text()+"\n")
	}
	// only data part
	err := yaml.Unmarshal([]byte(secretString.String()),&yamlData)
	if err != nil {
		fmt.Printf("Failed to decode yaml %v\n",err)
	}
	// complete yaml object
	err = yaml.Unmarshal([]byte(secretString.String()),&yamlComplete)
	if err != nil {
		fmt.Printf("Failed to decode yaml %v\n",err)
	}
	if !isSecretObject(&yamlData) {
		fmt.Println("The given object is not a Secret object")
		return
	}
	return &yamlData, yamlComplete, strings.ToLower(yamlData.Type)
}

// Validate the secret object kind
func isSecretObject(y *secretData) bool {
	return strings.ToLower(y.Kind) == "secret"
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
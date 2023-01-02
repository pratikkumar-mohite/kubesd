package src

import (
	"fmt"

	yaml "gopkg.in/yaml.v2"
	json "encoding/json"
)

const dataObject = "data"

func printSecretObject(yData *secretData, yComplete SecretYaml){
	var object []uint8
	var err error
	for key := range yComplete {
		if key == dataObject {
			yComplete[key] = yData.Data
		}
	}
	if OutputType == "json" {
		object, err = json.MarshalIndent(yComplete,"","\t")
	} else {
		object, err = yaml.Marshal(yComplete)
	}
	if err != nil {
		fmt.Printf("Failed to encode the secret object while printing %v\n",err)
	}
	fmt.Println(string(object))
}
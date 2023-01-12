package src

import (
	"fmt"

	json "encoding/json"
	yaml "gopkg.in/yaml.v2"
)

const dataObject = "data"
const stringDataObject = "stringData"

func (sComplete *SecretYaml)marshal() (output string){
	var object []uint8
	var err error

	for key := range *sComplete {
		if key == dataObject || key == stringDataObject {
			(*sComplete)[key] = Data
		}
	}

	if OutputType == "json" {
		object, err = json.MarshalIndent(sComplete, "", "\t")
	} else {
		object, err = yaml.Marshal(sComplete)
	}
	if err != nil {
		fmt.Printf("Failed to encode the secret object while printing %v\n", err)
	}
	return string(object)
}

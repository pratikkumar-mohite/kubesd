package src

import (
	"fmt"

	yaml "gopkg.in/yaml.v2"
	json "encoding/json"
)

const dataObject = "data"

func printSecretObject(s SecretYaml){
	var object []uint8
	var err error
	for key := range s {
		if key == dataObject {
			s[key] = Data
		}
	}
	if OutputType == "json" {
		object, err = json.MarshalIndent(s,"","\t")
	} else {
		object, err = yaml.Marshal(s)
	}
	if err != nil {
		fmt.Printf("Failed to encode the secret object while printing %v\n",err)
	}
	fmt.Println(string(object))
}
package src

import (
	"fmt"

	yaml "gopkg.in/yaml.v2"
)

const dataObject = "data"

func printSecretObject(yData *secretData, yComplete Secret){
	for key := range yComplete {
		if key == dataObject {
			yComplete[key] = yData.Data
		}
	}
	object, err := yaml.Marshal(yComplete)
	if err != nil {
		fmt.Printf("Failed to encode the secret object while printing %v\n",err)
	}
	fmt.Println(string(object))
}
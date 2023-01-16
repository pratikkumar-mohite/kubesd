package main

import (
	"fmt"
	"github.com/pratikkumar-mohite/kubesd/src"
)

func main() {
	var decodedData, err = src.Decode()
	if err != nil {
		fmt.Println("Secret object decode failed")
		return
	}
	fmt.Println(decodedData)
}

package main

import (
	"github.com/pratikkumar-mohite/kubesd/cmd/kubesd/cli"
	"github.com/pratikkumar-mohite/kubesd/pkg/logger"
)

func main() {
	log := logger.NewLogger()
	var decodedData, err = cli.Decode()
	if err != nil {
		log.Error(err.Error())
		return
	}
	log.Info(decodedData)
}

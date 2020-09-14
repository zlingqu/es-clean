package main

import (
	"log"

	"github.com/zlingqu/es-clean/cmd"
)

func main() {
	newClean := cmd.NewEsCleanCommand()
	if err := newClean.Execute(); err != nil {
		log.Fatal(err)
	}

}

// https://gitlab.dm-ai.cn/application-engineering/devops/service-clean-es-data/blob/master/modules/es/es.go

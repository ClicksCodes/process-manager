package main

import (
	"clicksminuteper.net/process-manager/containerManager"
	"log"
)

func main() {
	err := containerManager.BuildContainer("mini-test", "@latest")
	if err != nil {
		log.Println("Container failed to build, " + err.Error())
		return
	}
	log.Println("Container built, running...")

	err = containerManager.RunContainer("mini-test", "@latest")
	if err != nil {
		log.Println("Container failed to run, " + err.Error())
		return
	}
	log.Println("Ran container successfully")
}

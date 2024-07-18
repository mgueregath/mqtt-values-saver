package main

import (
	"ValuesImporter/facade/container"
	"ValuesImporter/facade/environment"
	"gopkg.in/yaml.v3"
	"os"
)

func main() {
	container.CreateApp(loadEnvironment())
}

func loadEnvironment() *environment.Environment {
	ENV := environment.Environment{}

	yamlFile, err := os.ReadFile("config.yml")
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(yamlFile, &ENV)
	if err != nil {
		panic(err)
	}

	return &ENV
}

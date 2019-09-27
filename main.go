package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v3"
)

type conf struct {
	Hits int64 `yaml:"hits"`
	Time int64 `yaml:"time"`
}

func (c *conf) getConf(pathToFile string) *conf {

	yamlFile, err := ioutil.ReadFile(pathToFile)
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	return c
}

func main() {
	var c conf
	c.getConf("conf.yaml")

	fmt.Println(c)
}

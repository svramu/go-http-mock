package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

// Rule ...
type Rule struct {
	Request  string
	Callback string
	Delay    int // In Millis
}

// Conf ...
type Conf struct {
	Env   map[string]string
	Rules []Rule
}

func (c *Conf) parse(pathToFile string) *Conf {
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
	var c Conf
	c.parse("conf.yaml")

	fmt.Println(c)

	for i, r := range c.Rules {
		fmt.Println(i, r.Callback)
		t := template.Must(template.New("callback").Parse(r.Callback))
		t.Execute(os.Stdout, c.Env)
		fmt.Println()
	}
}

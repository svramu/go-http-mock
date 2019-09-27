package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

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

func (c *Conf) match(path string) (Rule, error) {
	for _, r := range c.Rules {
		if strings.HasPrefix(r.Request, path) {
			return r, nil
		}
	}
	return Rule{}, fmt.Errorf("%q", path)
}

var c Conf

func main() {
	c.parse("conf.yaml")

	fmt.Println("---- ---- ---- ----")
	fmt.Println(time.Now().Format("2006 Aug 3"))
	fmt.Println(c)

	for i, r := range c.Rules {
		fmt.Println("...")
		fmt.Println(i)
		fmt.Println("  ", r.Request)
		fmt.Println("  ", r.Callback)
		t := template.Must(template.New("callback").Parse(r.Callback))
		fmt.Print("  ")
		t.Execute(os.Stdout, c.Env)
		fmt.Println()
	}

	fmt.Println("...")
	http.HandleFunc("/", handleAll)
	http.ListenAndServe(":6174", nil)
}

func handleAll(w http.ResponseWriter, req *http.Request) {
	out := req.URL.Path + "?" + req.URL.RawQuery
	fmt.Println(out, req.URL.Path)

	r, err := c.match(req.URL.Path)
	if err != nil {
		fmt.Println("error:", err.Error())
	} else {
		fmt.Println("ok:", r)
		t := template.Must(template.New("callback").Parse(r.Callback))
		t.Execute(os.Stdout, c.Env)
		fmt.Println()
	}

	// resp, err := http.Get("http://example.com/")
	// if err != nil {
	// 	// handle error
	// }
	// defer resp.Body.Close()

	// body, err := ioutil.ReadAll(resp.Body)
	// fmt.Println(string(body))
}

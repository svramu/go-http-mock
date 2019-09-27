package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
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
	http.HandleFunc("/", handleAll)
	http.ListenAndServe(":6174", nil)
}

func handleAll(w http.ResponseWriter, req *http.Request) {
	out := req.URL.Path + "?" + req.URL.RawQuery
	fmt.Println(out)

	r, err := c.match(req.URL.Path)
	if err != nil {
		fmt.Println("error 0:", err.Error())
	} else {
		callHTTP(transform(r.Callback, c.Env))
	}
}

func transform(ts string, data interface{}) string {
	buf := &bytes.Buffer{}
	t := template.Must(template.New("callback").Parse(ts))
	t.Execute(buf, data)
	fmt.Println("transform:", ts, buf.String())
	return buf.String()
}

func callHTTP(url string) {
	fmt.Println("http:", url)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("error 1:", err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("error 2:", err)
		return
	}

	fmt.Println(len(string(body)))
}

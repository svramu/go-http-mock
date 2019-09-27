package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"regexp"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// Rule ...
type Rule struct {
	Request         string
	RequestTemplate string `yaml:"request-template"`
	Callback        string
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
	fmt.Println(".")
	http.HandleFunc("/", handleAll)
	http.ListenAndServe(":6174", nil)
}

func handleAll(w http.ResponseWriter, req *http.Request) {
	out := req.URL.Path + "?" + req.URL.RawQuery
	fmt.Println("req\t>", out)

	r, err := c.match(req.URL.Path)
	if err != nil {
		//fmt.Println("error 0:", err.Error())
	} else {
		data := captureRegex(r.RequestTemplate, out)
		//fmt.Println(data)

		dur := time.Duration(500 + rand.Intn(500))
		time.Sleep(dur * time.Millisecond)
		callHTTP(transform(r.Callback, data))
	}
}

func captureRegex(tpl, txt string) map[string]string {
	re := regexp.MustCompile(tpl)
	values := re.FindStringSubmatch(txt)
	keys := re.SubexpNames()
	outMap := make(map[string]string)
	for i := 1; i < len(keys); i++ {
		outMap[keys[i]] = values[i]
	}
	return outMap
}

func transform(ts string, data interface{}) string {
	buf := &bytes.Buffer{}
	t := template.Must(template.New("callback").Parse(ts))
	t.Execute(buf, data)
	//fmt.Println("transform:", ts, buf.String())
	return buf.String()
}

func callHTTP(url string) {
	fmt.Println("call\t>", url)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("error 1:", err)
		return
	}
	defer resp.Body.Close()
	fmt.Println(".")
}

func getBody(body io.ReadCloser) string {
	b, err := ioutil.ReadAll(body)
	if err != nil {
		panic(err)
	}
	return string(b)
}

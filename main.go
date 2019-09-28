package main

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"regexp"
	"time"

	"gopkg.in/yaml.v3"
)

// Call ...
type Call struct {
	URL    string
	Method string
}

func (c *Call) method() string {
	if c.Method == "" {
		return "GET"
	}
	return c.Method
}

// Rule ...
type Rule struct {
	Request  Call
	Callback Call
}

// Conf ...
type Conf struct {
	//Env   map[string]string // TBD - not used now
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

func (c *Conf) match(method, url string) (Call, error) {
	outCall := Call{}
	for _, r := range c.Rules {
		//fmt.Println("**", r.Request)
		data := captureRegex(r.Request.URL, url)
		if data == nil {
			continue
		}
		if method != r.Request.method() {
			continue
		}
		outCall.URL = transform(r.Callback.URL, data)
		outCall.Method = r.Callback.Method
		return outCall, nil
	}
	//fmt.Println("***", method, url)
	return outCall, errors.New("no match")
}

func captureRegex(tpl, txt string) map[string]string {
	re := regexp.MustCompile(tpl)
	values := re.FindStringSubmatch(txt)
	if values == nil {
		return nil
	}
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

var c Conf

func main() {
	c.parse("conf.yaml")
	fmt.Println("---- ---- ---- ----")
	fmt.Println(time.Now())
	fmt.Println(".")
	http.HandleFunc("/", handleAll)
	http.ListenAndServe(":6174", nil)
}

func handleAll(w http.ResponseWriter, req *http.Request) {
	fullURL := req.URL.Path + "?" + req.URL.RawQuery
	method := req.Method
	fmt.Println(method, fullURL)

	call, err := c.match(method, fullURL)
	if err == nil {
		dur := time.Duration(500 + rand.Intn(500))
		time.Sleep(dur * time.Millisecond)
		callHTTP(call)
	}
}

func callHTTP(call Call) {
	fmt.Println(":", call.method(), call.URL)
	request, err := http.NewRequest(call.Method, call.URL, nil)
	if err != nil {
		log.Fatalln(err)
	}
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		fmt.Println(":", "Error: Unable to connect")
	} else {
		defer resp.Body.Close()
	}
}

func getBody(body io.ReadCloser) string {
	b, err := ioutil.ReadAll(body)
	if err != nil {
		panic(err)
	}
	return string(b)
}

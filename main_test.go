package main

import (
	"testing"
)

func TestCaptureRegex_Good1(t *testing.T) {
	got := captureRegex("/say/color.+blue=(?P<blue>[a-zA-Z]+).*", "http://localhost:6174/say/color?blue=orange")
	if got["blue"] != "orange" {
		t.Errorf("captureRegex not capturing: %v", got)
	}
}

func TestCaptureRegex_Good1_Var(t *testing.T) {
	got := captureRegex("/say/color.+blue=(?P<blue>[a-zA-Z]+).*", "/say/color?blue=orange")
	if got["blue"] != "orange" {
		t.Errorf("captureRegex not capturing: %v", got)
	}
}

func TestCaptureRegex_Negative(t *testing.T) {
	got := captureRegex("king", "bingo")
	if len(got) != 0 {
		t.Errorf("captureRegex not capturing: %v", got)
	}
}

func TestTransform_Good1(t *testing.T) {
	got := transform("http://localhost:6174/{{.blue}}", map[string]string{"blue": "green"})
	if got != "http://localhost:6174/green" {
		t.Errorf("transform not working: %v", got)
	}
}

func TestTransform_Negative1(t *testing.T) {
	got := transform("hola{{.red}}", map[string]string{"blue": "green"})
	if got != "hola" {
		t.Errorf("transform not working: %v", got)
	}
}

func TestTransform_Negative2(t *testing.T) {
	got := transform("hola{{.red}}", map[string]string{})
	if got != "hola" {
		t.Errorf("transform not working: %v", got)
	}
}

func testConf() Conf {
	var c Conf
	c.Rules = append(c.Rules, Rule{
		Request: Call{
			URL:    ".*/say/color.+blue=(?P<blue>[a-zA-Z]+).*",
			Method: "POST",
		},
		Callback: Call{
			URL:    "http://localhost:6174/{{.blue}}",
			Method: "GET",
		},
	})
	c.Rules = append(c.Rules, Rule{
		Request: Call{
			URL: "/say/23",
		},
		Callback: Call{
			URL: "http://localhost:3000/answer",
		},
	})
	return c
}

func TestConfMatch_Good1(t *testing.T) {
	c := testConf()
	got, err := c.match("POST", "http://localhost:6174/say/color?blue=orange")
	if err != nil {
		t.Errorf("ConfMatch unknown error: %v", err)
		return
	}
	if got.URL != "http://localhost:6174/orange" {
		t.Errorf("ConfMatch unknown match: %v", got)
		return
	}
	if got.method() != "GET" {
		t.Errorf("ConfMatch unknown Method: %v", got)
		return
	}
}

func TestConfMatch_Good2(t *testing.T) {
	c := testConf()
	got, err := c.match("GET", "http://localhost:6174/say/23")
	if err != nil {
		t.Errorf("ConfMatch unknown error: %v", err)
		return
	}
	if got.URL != "http://localhost:3000/answer" {
		t.Errorf("ConfMatch unknown match: %v", got)
		return
	}
	if got.method() != "GET" {
		t.Errorf("ConfMatch unknown Method: %v", got)
		return
	}
}

func TestConfMatch_Unknown(t *testing.T) {
	c := testConf()
	got, err := c.match("GET", "http://localhost:6174/say/23")
	if err != nil {
		t.Errorf("ConfMatch unknown error: %v", err)
		return
	}
	if got.URL != "http://localhost:3000/answer" {
		t.Errorf("ConfMatch unknown match: %v", got)
		return
	}
	if got.method() != "GET" {
		t.Errorf("ConfMatch unknown Method: %v", got)
		return
	}
}

func TestConfMatch_Fail_Method(t *testing.T) {
	c := testConf()
	got, err := c.match("POST", "http://localhost:6174/say/23")
	if err == nil {
		t.Errorf("ConfMatch should be an error: %v", got)
		return
	}
}

func TestConfMatch_Fail_Match(t *testing.T) {
	c := testConf()
	got, err := c.match("GET", "http://localhost:6174/say/x23")
	if err == nil {
		t.Errorf("ConfMatch should be an error: %v", got)
		return
	}
}

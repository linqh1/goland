package main

import (
	"fmt"
	"github.com/go-yaml/yaml"
	"io/ioutil"
)

type Config struct {
	Timeout int `yaml:timeout`
	URL     struct {
		Source      string
		Destination string
	}
}

func main() {
	config := Config{}
	bytes, e := ioutil.ReadFile("C:\\Users\\linqh1\\Desktop\\test.txt")
	if e != nil {
		panic("Read File Error!" + e.Error())
	}
	e = yaml.Unmarshal(bytes, &config)
	if e != nil {
		panic("Unmarshal File Error!" + e.Error())
	}
	fmt.Println(config.Timeout, config.URL.Destination, config.URL.Source)
}

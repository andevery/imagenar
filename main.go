package main

import (
	"flag"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"time"
)

var (
	envArg = flag.String("e", "development", "Environment")
)

func main() {
	flag.Parse()

	err := setEnv(*envArg)
	if err != nil {
		log.Fatal(err)
	}

	data, err := ioutil.ReadFile("db/dbconf.yml")
	if err != nil {
		log.Fatal(err)
	}

	var conf DBConf
	err = yaml.Unmarshal(data, &conf)
	if err != nil {
		log.Fatal(err)
	}

	dispatcher, err := NewDispatcher(conf)
	if err != nil {
		log.Fatal(err)
	}

	dispatcher.Start()
	time.Sleep(time.Minute)
	dispatcher.Stop()
}

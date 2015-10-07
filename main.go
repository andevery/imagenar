package main

import (
	"flag"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
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

	var conf dbConf
	err = yaml.Unmarshal(data, &conf)
	if err != nil {
		log.Fatal(err)
	}

	err = dbConnect(conf)
	if err != nil {
		log.Fatal(err)
	}
}

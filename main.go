package main

import (
	"flag"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"syscall"
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
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	log.Println(<-ch)

	// Stop the service gracefully.
	dispatcher.Stop()
}

package main

import (
	"flag"
	"log"

	"autocare.org/sandpiper/pkg/api"
	"autocare.org/sandpiper/pkg/config"
)

func main() {

	cfgPath := flag.String("p", "./sandpiper.config.yaml", "Path to config file")
	flag.Parse()

	cfg, err := config.Load(*cfgPath)
	if err != nil {
		log.Fatal("ERROR: ",err)
	}

	err = api.Start(cfg)
	if err != nil {
		panic(err.Error())
	}
}


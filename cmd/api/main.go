// Copyright Auto Care Association. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

package main

import (
	"flag"
	"fmt"
	"log"

	"autocare.org/sandpiper/pkg/api"
	"autocare.org/sandpiper/pkg/config"
)

var version = "undefined"

func main() {
	fmt.Printf("Sandpiper (%s)\n", version)

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


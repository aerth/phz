package main

import (
	"flag"
	"log"
	"os"
	"strings"
	"x/phzd/phz"

	"github.com/BurntSushi/toml"
)

func main() {
	var (
		configpath = flag.String("conf", "", "config path")
		debug      = flag.Bool("v", false, "debug mode")
		err        error
		envmap     = map[string]string{}
	)

	flag.Parse()

	// default config, change this
	conf := &phz.Config{
		Data: map[string]interface{}{
			"Env": envmap,
		},
		Debug: *debug,
	}

	// read config into default
	if *configpath != "" {
		if _, err = toml.DecodeFile(*configpath, conf); err != nil {
			log.Fatalln(err)
		}
	}

	// read env into config
	for _, v := range os.Environ() {
		split := strings.Split(v, "=")
		key, val := split[0], split[1]
		envmap[key] = val
	}

	// execute the phz templates
	for _, filename := range flag.Args() {
		if err := conf.ExecFile(os.Stdout, filename); err != nil {
			log.Fatalln(err)
		}

	}
}

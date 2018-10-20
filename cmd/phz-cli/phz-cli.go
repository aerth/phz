// Copyright (c) 2018 aerth. All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are
// met:
//
//    * Redistributions of source code must retain the above copyright
// notice, this list of conditions and the following disclaimer.
//    * Redistributions in binary form must reproduce the above
// copyright notice, this list of conditions and the following disclaimer
// in the documentation and/or other materials provided with the
// distribution.
//    * Neither the name of aerth nor the names of this project's
// contributors may be used to endorse or promote products derived from
// this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
// "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
// LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
// A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
// OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
// SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
// LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
// DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
// THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

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

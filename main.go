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

package main // import "x/phzd"

import (
	"flag"
	"log"
	serverlib "x/phzd/phz"

	"github.com/BurntSushi/toml"
)

func main() {

	var (
		confpath  = flag.String("conf", "config.toml", "path to TOML config")
		addrflag  = flag.String("addr", "", "address to override config (format: 127.0.0.1:8080)")
		debugflag = flag.Bool("v", false, "verbose / debug logs")
	)
	log.SetFlags(0)
	flag.Parse()
	config := serverlib.NewDefaultConfig()
	if _, err := toml.DecodeFile(*confpath, config); err != nil {
		log.Fatalln(err)
	}
	if *addrflag != "" {
		config.Addr = *addrflag
	}
	if *debugflag {
		config.Debug = true
	}
	if config.Debug {
		log.SetFlags(log.Ltime | log.Lshortfile)
	}
	srv := serverlib.NewServer(*config)
	log.Println("Serving http://" + config.Addr)
	log.Fatalln(srv.ListenAndServe())

}

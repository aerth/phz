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

package phz

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"time"
)

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method, r.Host, r.URL.Path, r.RemoteAddr, r.UserAgent())
	switch r.Method {
	case http.MethodGet: // ok!
	default:
		http.Error(w, "bad method", http.StatusMethodNotAllowed)
		return
	}

	for _, v := range RestrictedPathKeywords {
		if strings.Contains(r.URL.Path, v) {
			http.Error(w, "bad url", http.StatusForbidden)
			return
		}
	}

	if strings.HasSuffix(r.URL.Path, "/") {
		r.URL.Path += "index.phz"
	}
	if strings.HasSuffix(r.URL.Path, ".html") {
		r.URL.Path = strings.TrimSuffix(r.URL.Path, ".html")
		r.URL.Path += ".phz"
	}
	path := strings.Split(strings.TrimPrefix(r.URL.Path, "/"), "/")
	log.Println("Checking path[0]:", path[0])
	switch path[0] {
	default:
		split := strings.Split(r.URL.Path, "/")
		if len(split) > 1 {
			if strings.HasSuffix(split[len(split)-1], ".phz") {
				s.handleGETphz(w, r, strings.TrimPrefix(r.URL.Path, "/")) // TODO: dry
				return
			} else {
				http.ServeFile(w, r, filepath.Join(s.config.TemplatePath, r.URL.Path))
				return
			}
		}
		log.Println("not found:", r.URL.Path, split)
		http.NotFound(w, r)
		return
	case "a":
		fmt.Println("AAA")
	case "stats":
		fmt.Println(time.Now().UTC())
	}
}

func (s *Server) handleGETphz(w http.ResponseWriter, r *http.Request, pathnoslash string) {
	log.Println("PHZ:", pathnoslash)
	if err := s.ServeTemplate(w, r, pathnoslash); err != nil {
		log.Println("servtemplate:", err)
	}
}

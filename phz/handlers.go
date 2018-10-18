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

func containsBadWords(s ...string) bool {
	// O(9000)
	for _, v := range RestrictedPathKeywords {
		for i := range s {
			if strings.Contains(s[i], v) {
				return true
			}
		}
	}
	return false
}
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	reqdata := map[string]interface{}{}
	t1 := time.Now()
	reqdata["t1"] = t1
	if s.config.Debug {
		defer func() {
			//fmt.Fprintf(w, "\nRequest process took %s. Powered by <a href='https://phz'>phz</a>!\n", time.Since(t1))
			fmt.Fprintf(w, "\nRequest processed: %s\n", time.Since(t1))
		}()
	}
	log.Println(r.Method, r.Host, r.URL.Path, r.RemoteAddr, r.UserAgent())
	switch r.Method {
	case http.MethodGet: // ok!
	case http.MethodPost:
		r.ParseMultipartForm(1024)
		for k, v := range r.PostForm {
			reqdata["post_"+k] = v
		}
	default:
		s.Error(w, r, http.StatusMethodNotAllowed)
		return
	}
	for k, v := range r.URL.Query() {
		reqdata["get_"+k] = v
	}
	if containsBadWords(r.URL.Path) {
		s.Error(w, r, http.StatusForbidden)
		return
	}

	if strings.HasSuffix(r.URL.Path, "/") {
		r.URL.Path += "index.phz"
	}

	for _, old := range []string{".html", ".md"} {
		if strings.HasSuffix(r.URL.Path, old) {
			r.URL.Path = strings.TrimSuffix(r.URL.Path, old)
			r.URL.Path += ".phz"
		}
	}
	pathnoprefix := strings.TrimPrefix(r.URL.Path, "/")
	pathpart := strings.Split(pathnoprefix, "/")

	if builtinWebHandler(w, r, pathpart) {
		return
	}
	filename := filepath.Base(pathnoprefix)
	if strings.HasSuffix(filename, ".phz") {
		log.Printf("Serving dynamic phz file: %q from %q", filename, pathnoprefix)
		if err := s.phzhandler(w, r, pathnoprefix, reqdata); err != nil {
			log.Printf("Error serving phz file: %s %v", filename, err)
		}
		return
	}

	// handle everything else (and 404s)
	staticfilepath := filepath.Join(s.config.TemplatePath, pathnoprefix)
	log.Printf("Serving static file: %q from %q", staticfilepath, pathnoprefix)
	http.ServeFile(w, r, staticfilepath)
}

func builtinWebHandler(w http.ResponseWriter, r *http.Request, pathpart []string) (handledProperly bool) {
	switch pathpart[0] {
	default:
		return false
	case "a":
		fmt.Fprintln(w, "AAA")
		fmt.Fprintln(w)
		return true
	case "stats":
		fmt.Fprintln(w, time.Now().UTC())
		fmt.Fprintln(w)
		return true
	}
}

func (s *Server) handleGETphz(w http.ResponseWriter, r *http.Request, pathnoslash string, try int) error {
	if try > 1 {
		return fmt.Errorf("too many tries: %s", pathnoslash)
	}
	err := s.ServeTemplate(w, r, pathnoslash)
	if err == nil {
		return nil
	}
	if strings.Contains(err.Error(), "no such template") {
		log.Println("servtemplate:", err)
		deps := strings.Split(err.Error(), "no such template ")
		if len(deps) == 2 {
			// dependent / include
			dtmpl := strings.TrimPrefix(strings.TrimSuffix(deps[1], `"`), `"`)
			log.Println("template depends on", dtmpl, "-- reloading it")
			if err := s.reloadtemplate(dtmpl); err != nil {
				log.Println("err reloading", err)
			}
			log.Println("relaoding:", pathnoslash)
			s.reloadtemplate(pathnoslash)
			return s.handleGETphz(w, r, pathnoslash, try+1)

		}
	} else {
		log.Println("Error handling template:", err)
	}

	// print debug all template names

	if s.config.Debug {
		log.Println("all templates:", s.template.Templates())
		/*		if s.templates[pathnoslash] != nil {
					ts := s.templates[pathnoslash].Templates()
					for i := range ts {
						log.Println(ts[i].Name())
					}
				}
		*/
	}

	return nil
}

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
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var RestrictedPathKeywords = []string{
	"..", // up
	"/.", // hidden files
}

func ContainsBadWords(s ...string) bool {
	for i := range s {
		for j := 0; j < len(RestrictedPathKeywords); j++ {
			if strings.Contains(s[i], RestrictedPathKeywords[j]) {
				return true
			}
		}
	}
	return false
}

func contains(needle string, haystack []string) bool {
	for i := range haystack {
		if needle == haystack[i] {
			return true
		}
	}
	return false
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if isBuiltInHandled(w, r) {
		return
	}
	if ContainsBadWords(r.URL.Path) {
		s.Error(w, r, 403)
		log.Println("WARN", "intruder alert")
		return
	}
	if strings.HasSuffix(r.URL.Path, "/") {
		r.URL.Path += "index.phz"
	}
	filename := filepath.Join(s.config.TemplatePath, r.URL.Path)

	if !contains(r.Method, []string{http.MethodGet, http.MethodPost, http.MethodHead}) {
		s.Error(w, r, http.StatusMethodNotAllowed)
		return
	}

	if _, err := os.Stat(filename); err != nil {
		http.NotFound(w, r)
		log.Println(r.URL.Path, err)
		return
	}

	if strings.HasSuffix(r.URL.Path, ".phz") {
		if err := s.ServePHZ(w, r); err != nil {
			log.Println(err)
			s.Error(w, r, 503)
		}
		return
	}
	http.ServeFile(w, r, filename)
}

func isBuiltInHandled(w http.ResponseWriter, r *http.Request) bool {
	switch r.URL.Path {
	case "/stats":
		fmt.Fprintln(w, time.Now().UTC())
		return true
	default:
	}

	switch strings.Split(strings.TrimPrefix(r.URL.Path, "/"), "/")[0] {
	case "a":
		fmt.Println("AAA") // admin?
		return true
	}

	return false
}

func getformdata(r *http.Request) map[string]interface{} {
	log.Println(r.Method, r.Host, r.URL.Path, r.RemoteAddr, r.UserAgent())
	reqdata := map[string]interface{}{}
	switch r.Method {
	case http.MethodGet: // ok!
	case http.MethodPost:
		r.ParseMultipartForm(1024)
		for k, v := range r.PostForm {
			reqdata["post_"+k] = v
		}
	default:
		return nil
	}
	for k, v := range r.URL.Query() {
		reqdata["get_"+k] = v
	}

	return reqdata
}

func (s *Server) ServePHZ(w http.ResponseWriter, r *http.Request) error {
	mainfile := filepath.Join(s.config.TemplatePath, r.URL.Path)
	if s.config.Debug {
		log.Println("phz running", mainfile)
	}
	files := []string{mainfile}
	t, err := template.ParseFiles(files...)
	if err != nil {
		log.Println("error parsing template:", err)
		return err
	}

	t = t.Option("missingkey=zero").Funcs(s.globalfuncs)
	formdata := getformdata(r)
	if formdata == nil {
		s.Error(w, r, http.StatusMethodNotAllowed)
		return nil
	}
	inputdata := map[string]interface{}{
		"Now":  time.Now(),
		"Foo":  "bar",
		"Req":  *r,
		"Form": formdata,
	}
	buf := new(bytes.Buffer)
	err = t.Execute(buf, inputdata)
	if err == nil {
		w.Write(ParseMarkdown(buf.Bytes()))
		return nil
	}

	buf.Reset()

	log.Println("Reload experiment 1 START")
	pathnoslash := strings.TrimPrefix(r.URL.Path, "/")
	// there was an error in execution
	if err := s.refreshinclude(pathnoslash, err); err != nil {
		log.Println("error refreshing template:", err)
		return nil
	}

	err = t.Execute(buf, inputdata)
	if err == nil {
		log.Println("Reload experiment 1 PASS")
		w.Write(ParseMarkdown(buf.Bytes()))
		return nil
	}

	log.Println("Reload experiment 1 FAIL", err)
	return err
}

func (s *Server) reloadtemplate(name string) error {
	_, err := s.template.ParseFiles(filepath.Join(s.config.TemplatePath, name))
	return err
}

func (s *Server) refreshinclude(pathnoslash string, err error) error {
	if !strings.Contains(err.Error(), "no such template") {
		return err
	}
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
	}
	return nil
}

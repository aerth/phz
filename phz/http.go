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

// ServeHTTP main entrypoint
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
		t1 := time.Now()
		if err := s.ServePHZ(w, r); err != nil {
			log.Println(err)
			s.Error(w, r, 503)
		}
		fmt.Fprintf(w, "phz processed: %s\n", time.Since(t1))
		return
	}
	http.ServeFile(w, r, filename)
}

// a few built-ins, can add more
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
	case http.MethodGet: // get and post are ok
	case http.MethodPost:
		r.ParseMultipartForm(1024)
		for k, v := range r.PostForm {
			reqdata["post_"+k] = v[0] // TODO: more than the first value?
		}
	default:
		return nil // if nil, throw err
	}
	for k, v := range r.URL.Query() {
		reqdata["get_"+k] = v[0] // TODO: more than the first value? or get last value?
	}
	return reqdata
}

func (s *Server) ServePHZ(w http.ResponseWriter, r *http.Request) error {
	var (
		mainfile = filepath.Join(s.config.TemplatePath, r.URL.Path)
		files    = []string{mainfile}
		t        = template.New(".root").Option("missingkey=zero").Funcs(s.globalfuncs)
		buf      bytes.Buffer
	)
	/* if s.config.Debug {
		log.Println("Funcs:", len(s.globalfuncs))
		for i, v := range s.globalfuncs {
			log.Printf("func: %s: %s", i, funcName(v))
		}
	} */

	// get file, (prepare for "files soon")
	t.ParseFiles(files...)

	// get first level phz files ( even unservable dot files! )
	t.ParseGlob(filepath.Join(s.config.TemplatePath, "*.phz"))
	formdata := getformdata(r)
	if formdata == nil {
		s.Error(w, r, http.StatusMethodNotAllowed)
		return nil
	}
	inputdata := map[string]interface{}{
		"Now":  time.Now(),
		"Req":  *r,
		"Form": formdata,
	}
	// header strings (should use first or last? or join? how to index
	// easier?)
	hmap := map[string]string{}
	for i, v := range r.Header {
		hmap[i] = v[0]
	}
	inputdata["Header"] = hmap

	// begin execute
	pathnoslash := strings.TrimPrefix(r.URL.Path, "/")
	templatename := filepath.Base(pathnoslash)
	t2, err := t.Clone()
	if err != nil {
		log.Println("err cloning", err)
		return err
	}
	err = t2.ExecuteTemplate(&buf, templatename, inputdata)
	if err == nil {
		b := buf.Bytes()
		if len(b) >= 14 && bytes.Compare(b[:14], []byte("<!DOCTYPE html")) == 0 {
			w.Write(buf.Bytes())
			return nil
		}
		w.Write(ParseMarkdown(buf.Bytes()))
		return nil
	}

	buf.Reset()

	log.Println("Reload experiment 1, START err:", err)
	// there was an error in execution
	if err := s.refreshinclude(t, pathnoslash, err); err != nil {
		log.Println("error refreshing template:", err)
		return nil
	}

	err = t.ExecuteTemplate(&buf, templatename, inputdata)
	if err == nil {
		log.Println("Reload experiment 1 PASS")
		w.Write(ParseMarkdown(buf.Bytes()))
		return nil
	}

	log.Println("Reload experiment 1 FAIL", err)
	return err
}

func (s *Server) reloadtemplate(t *template.Template, name string) error {
	t, err := t.Clone()
	if err != nil {
		return err
	}
	_, err = t.ParseFiles(filepath.Join(s.config.TemplatePath, name))
	return err
}

func (s *Server) refreshinclude(t *template.Template, pathnoslash string, err error) error {
	if !strings.Contains(err.Error(), "no such template") {
		return err
	}
	deps := strings.Split(err.Error(), "no such template ")
	if len(deps) == 2 {
		// dependent / include
		dtmpl := strings.TrimPrefix(strings.TrimSuffix(deps[1], `"`), `"`)
		log.Println("template depends on", dtmpl, "-- reloading it")
		if err := s.reloadtemplate(t, dtmpl); err != nil {
			log.Println("err reloading", err)
		}
		log.Println("relaoding:", pathnoslash)
		s.reloadtemplate(t, pathnoslash)
	}
	return nil
}

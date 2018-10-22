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
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Exec a template
// web server should use this, it doesnt.
func (c Config) Exec(w io.Writer, path string, b []byte) error {
	var err error
	t := template.New(path).Option("missingkey=zero").Funcs(DefaultFuncMap)
	// includes
	if c.TemplatePath != "" {
		t, err = t.ParseGlob(filepath.Join(c.TemplatePath, "*.phz"))
		if err != nil {
			return fmt.Errorf("could not parse includable templates, clean up TemplatePath: %v", err)
		}
	}
	t, err = t.Parse(string(b))
	if err != nil {
		return err
	}
	if c.Data == nil {
		c.Data = make(map[string]interface{})
	}
	c.Data["Now"] = time.Now().UTC()
	c.Data["Path"] = path
	envmap := map[string]string{}
	for _, v := range os.Environ() {
		split := strings.Split(v, "=")
		envmap[split[0]] = split[1]
	}
	c.Data["Env"] = envmap

	// dummy web stuff
	c.Data["Req"] = http.Request{}
	c.Data["Header"] = map[string]string{}
	c.Data["Form"] = map[string]string{}

	args := os.Args[0:]
	if len(args) > 1 {
		args = args[1:]
	}
	c.Data["Args"] = args
	return t.Execute(w, c.Data)
}

func (c Config) ExecFile(w io.Writer, path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	b, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	index := 0
	if b[0] == '#' && b[1] == '!' {
		for i := range b {
			if b[i] == '\n' {
				index = i
				break
			}
		}
	}
	return c.Exec(w, path, b[index:])
}

// things are not printed to web site, only to server log
func (s *Server) Error(w http.ResponseWriter, r *http.Request, code int, things ...interface{}) {
	log.Printf("%s: %v %s", r.URL.Path, code, fmt.Sprint(things...))
	http.Error(w, http.StatusText(code), code)
}

func (s *Server) phzhandler(w http.ResponseWriter, r *http.Request, path string, formdata map[string]interface{}) error {
	buf := new(bytes.Buffer)
	t, err := s.template.Clone()
	if err != nil {
		http.Error(w, "dang", 503)
		log.Println(t.DefinedTemplates())
		return err
	}
	data := map[string]interface{}{}
	for i, v := range s.data {
		data[i] = v
	}
	data["Path"] = path
	data["GenTime"] = time.Since(formdata["t1"].(time.Time))
	data["Req"] = r
	data["Now"] = time.Now().UTC()
	data["Form"] = formdata
	data["t1"] = formdata["t1"]
	if err := t.ExecuteTemplate(buf, path, data); err != nil {
		http.Error(w, "dang", 503)
		log.Println(t.DefinedTemplates())
		return fmt.Errorf("Error executing template: %v", err)
	}
	w.Write(ParseMarkdown(buf.Bytes()))
	return nil

}

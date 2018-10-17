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
	"html/template"
	"net/http"
	"sync"
)

var RestrictedPathKeywords = []string{
	"..", // up
	"/.", // hidden files
}

type Config struct {
	Addr         string
	TemplatePath string
	Debug        bool
}

type Server struct {
	config Config
	mu     *sync.Mutex // guards global data map
	data   map[string]interface{}

	templates    map[string]*template.Template
	templatelock *sync.Mutex // guards template map
}

func NewDefaultConfig() *Config {
	return &Config{
		Addr:         "127.0.0.1:8000",
		TemplatePath: "./",
	}
}

func NewServer(c Config) *Server {
	return &Server{
		config:       c,
		data:         map[string]interface{}{},
		templates:    map[string]*template.Template{},
		mu:           new(sync.Mutex),
		templatelock: new(sync.Mutex),
	}
}

func (s *Server) ListenAndServe() error {
	return http.ListenAndServe(s.config.Addr, s)
}

func (s *Server) DataSet(str string, v interface{}) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[str] = v
}

func (s *Server) DataGet(str string) interface{} {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.data[str]
}

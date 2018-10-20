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
	"io"
	"io/ioutil"
	logg "log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func init() {
	//	log.SetFlags(logg.Lshortfile)
	//	log.SetOutput(ioutil.Discard)
	_ = ioutil.Discard
	log.SetFlags(logg.Lshortfile)
}

func TestOne(t *testing.T) {
	type testcase struct {
		path   string
		status int // status code response
	}
	for _, tc := range []testcase{
		{"/a/lol", 200}, //builtin prefix
		{"/a/lol2", 200},
		{"/a/lol3", 200},
		{"/", 200},    // good, not cached
		{"/", 200},    // good, cached
		{"/", 200},    // good, cached
		{"/../", 403}, // bad paths
		{"//.", 403},
		{"/.", 403},
		{"/.htaccess", 403},
		{"/.config", 403},
		{"/stats", 200},  // builtin
		{"/status", 404}, // not a builtin, not a real file
		{"/this/path/doesnt/exist.phz", 404},
		{"/this/path/doesnt/exist.png", 404},
	} {
		//fmt.Println("Testing:", tc.path)
		fmt.Println()
		testone(t, tc.path, tc.status)
	}
}

func testone(t *testing.T, path string, wantstatus int) {
	s := NewServer(Config{Debug: true, TemplatePath: "../testdata"})
	testserver := httptest.NewServer(s)
	url := testserver.URL
	req, _ := http.NewRequest(http.MethodGet, url+path, nil)
	req.Header.Add("User-Agent", "phz test client")
	resp, err := new(http.Client).Do(req)
	if err == nil && resp.StatusCode != wantstatus {
		err = fmt.Errorf("STATUS RETURNED %v for path=%q", resp.StatusCode, path)
	}
	if err != nil {
		t.Error(err)
		return
	}
	buf := new(bytes.Buffer)
	io.Copy(buf, resp.Body)
	resp.Body.Close()
	if true || resp.StatusCode == 200 && !strings.HasPrefix(path, "/a") {
		t.Log(buf.String())
	}
}

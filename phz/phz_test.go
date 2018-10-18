package phz

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func init() {
	log.SetFlags(log.Lshortfile)
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

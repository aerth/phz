package phz

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestOne(t *testing.T) {
	type testcase struct {
		path   string
		status int // status code response
	}
	for _, tc := range []testcase{
		{"/", 200},
		{"/../", 403},
		{"//.", 403},
		{"/.htaccess", 403},
		{"/.config", 403},
		{"/stats", 200},
		{"/status", 404},
		{"/.", 403},
		{"/a/lol", 200},
		{"/a/lol2", 200},
		{"/a/lol3", 200},
	} {
		fmt.Println("Testing:", tc.path)
		testone(t, tc.path, tc.status)
	}
}

func testone(t *testing.T, path string, wantstatus int) {
	s := &Server{}
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
	}
	buf := new(bytes.Buffer)
	io.Copy(buf, resp.Body)
	resp.Body.Close()
	fmt.Printf("% X\n", buf.Bytes())

}

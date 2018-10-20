package phz

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSimple(t *testing.T) {
	conf := Config{TemplatePath: "../testdata", Addr: ":0", Debug: true}
	server := NewServer(conf)
	testserver := httptest.NewServer(server)
	url := testserver.URL
	buf := new(bytes.Buffer)

	path := "/"
	// create request
	req, _ := http.NewRequest(http.MethodGet, url+path, nil)
	req.Header.Add("User-Agent", "phz test client")
	resp, err := new(http.Client).Do(req)
	if err != nil {
		t.Error(err)
		return
	}
	if resp.StatusCode != 200 {
		t.FailNow()
	}
	io.Copy(buf, resp.Body)
	resp.Body.Close()
	t.Log(buf.String())
	buf.Reset()
}

package phz

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

// Exec a template
// web server should use this, it doesnt.
func (c Config) Exec(w io.Writer, path string, b []byte) error {
	t1 := time.Now()
	t, err := template.New(path).Funcs(DefaultFuncMap).Parse(string(b))
	if err != nil {
		return err
	}
	if c.Data == nil {
		c.Data = make(map[string]interface{})
	}
	c.Data["Now"] = time.Now().UTC()
	c.Data["Path"] = path
	c.Data["GenTime"] = time.Since(t1)
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

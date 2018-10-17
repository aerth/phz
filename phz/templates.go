package phz

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"time"
)

func (s *Server) ExecuteTemplate(templatename string, input interface{}, output io.Writer) error {
	var err error
	if s.templates[templatename] == nil {
		log.Println("INitializing template:", templatename)
		if err := s.reloadtemplate(templatename); err != nil {
			return err
		}
	}
	markdowner := new(bytes.Buffer)
	if err := s.templates[templatename].ExecuteTemplate(markdowner, templatename, input); err != nil {
		return err
	}
	_, err = output.Write(ParseMarkdown(markdowner.Bytes()))
	return err

}

func (s *Server) reloadtemplate(templatename string) error {
	log.Println("reloading template:", templatename)
	t, err := s.template.Clone()

	if err != nil {
		return err
	}

	t, err = t.New(templatename).Parse(s.gettemplatestring(templatename))
	if err != nil {
		return err
	}
	s.templates[templatename] = t
	return nil
}
func (s *Server) gettemplatestring(name string) string {
	if name == "" {
		return "errar 1"
	}

	for _, v := range RestrictedPathKeywords {
		if strings.Contains(name, v) {
			return "errar 2"
		}
	}
	if name == "/" {
		name = "index"
	}
	b, err := ioutil.ReadFile(filepath.Join(s.config.TemplatePath, name))
	if err != nil {
		log.Println("Templates:", err)
		return "errar 4"
	}
	return string(b)
}

func (s *Server) ServeTemplate(w http.ResponseWriter, r *http.Request, path string) error {
	data := map[string]interface{}{
		"now": time.Now().UTC(),
		"req": *r,
	}
	return s.ExecuteTemplate(path, data, w)
}

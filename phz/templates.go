package phz

import (
	"bytes"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"time"
)

func (s *Server) ParseTemplate(templatename string, input interface{}, output io.Writer) error {
	s.templatelock.Lock()
	defer s.templatelock.Unlock()
	if s.templates[templatename] == nil {
		tmpl, err := template.New(templatename).Parse(s.gettemplatestring(templatename))
		if err != nil {
			return err
		}
		s.templates[templatename] = tmpl
	}

	markdowner := new(bytes.Buffer)
	if err := s.templates[templatename].Execute(markdowner, input); err != nil {
		return err
	}
	_, err := output.Write(ParseMarkdown(markdowner.Bytes()))
	return err

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
	b, err := ioutil.ReadFile(filepath.Join(s.config.TemplatePath, name+".phz"))
	if err != nil {
		log.Println("Templates:", err)
		return "errar 3"
	}
	return string(b)
}

func (s *Server) ServeTemplate(w http.ResponseWriter, r *http.Request, path string) error {
	data := map[string]interface{}{
		"now": time.Now().UTC(),
		"req": *r,
	}
	return s.ParseTemplate(path, data, w)
}

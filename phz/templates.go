package phz

import (
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"strings"
)

func (s *Server) ParseTemplate(templatename string, input []byte, output io.Writer) error {
	s.templatelock.Lock()
	defer s.templatelock.Unlock()
	if s.templates[templatename] == nil {
		tmpl, err := template.New(templatename).Parse(s.gettemplatestring(templatename))
		if err != nil {
			return err
		}
		s.templates[templatename] = tmpl
	}
	return s.templates[templatename].Execute(output, input)
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
	return s.ParseTemplate(path, nil, w)
}

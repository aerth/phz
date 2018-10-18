package phz

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var CacheTime = time.Minute

func (s *Server) executeTemplate(templatename string, input map[string]interface{}, output io.Writer) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	t, err := s.template.Clone()
	if err != nil {
		return err
	}
	/* if s.templates[templatename] == nil || time.Since(s.cache[templatename]) > CacheTime {
		log.Println("(re)initializing template:", templatename)
		if err = s.reloadtemplate(templatename); err != nil {
			return err
		}
	}
	*/
	// turn the template into html now
	buf := new(bytes.Buffer)
	now := time.Now().UTC()
	input["Now"] = now
	//s.cache[templatename] = now
	input["Path"] = templatename
	//input["GenTime"] = time.Since(t1)
	if err = t.ExecuteTemplate(buf, templatename, input); err != nil {
		return err
	}
	_, err = output.Write(ParseMarkdown(buf.Bytes()))
	return err

}

func (s *Server) reloadtemplate(templatename string) error {
	strbuf, err := s.gettemplatestring(templatename)
	if err != nil {
		return err
	}
	_, err = s.template.New(templatename).Parse(strbuf)
	return err
	log.Println("reloading template:", templatename)
	if s.template == nil {
		return fmt.Errorf("nil template root")
	}
	t, err := s.template.Clone()
	if err != nil {
		return err
	}
	templatebuf, err := s.gettemplatestring(templatename)
	if err != nil {
		return err
	}
	_, err = t.New(templatename).Parse(templatebuf)
	if err != nil {
		log.Println("err reloading")
		return err
	}
	//s.templates[templatename] = t
	s.template = t
	return nil
}

func (s *Server) errorcode(code int) string {
	return http.StatusText(code)
}
func (s *Server) gettemplatestring(name string) (string, error) {
	if name == "" || name == "/" {
		return "", fmt.Errorf("bad name (empty)")
	}
	for _, v := range RestrictedPathKeywords {
		if strings.Contains(name, v) {
			return "", fmt.Errorf("bad name")
		}
	}
	b, err := ioutil.ReadFile(filepath.Join(s.config.TemplatePath, name))
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (s *Server) ServeTemplate(w http.ResponseWriter, r *http.Request, path string) error {
	data := map[string]interface{}{
		"now": time.Now().UTC(),
		"req": *r,
	}
	if !strings.HasSuffix(path, ".phz") {
		return fmt.Errorf("invalid suffix, router malfunction. shutdown recommended!: %q", path)
	}
	return s.executeTemplate(path, data, w)
}

func (s *Server) loadtemplates() error {
	s.templatelock.Lock()
	defer s.templatelock.Unlock()
	var paths []string
	addpaths := func(path string, info os.FileInfo, err error) error {
		if !strings.HasSuffix(path, ".phz") {
			return nil
		}
		paths = append(paths, path)
		return nil
	}
	if err := filepath.Walk(s.config.TemplatePath, addpaths); err != nil {
		return err
	}

	log.Printf("Loading %v templates", len(paths))
	s.template = template.New(".root").Funcs(s.globalfuncs)
	var err error
	_, err = s.template.ParseGlob(s.config.TemplatePath + "/*.phz")
	if err != nil {
		log.Println(err)
		return err
	}
	log.Println("Includable templates:", s.template.DefinedTemplates())
	for _, path := range paths {
		log.Println("Loading template:", path)
		b, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}
		templatename := strings.TrimPrefix(path, s.config.TemplatePath+"/")
		_, err = s.template.New(templatename).Funcs(s.globalfuncs).Parse(string(b))
		if err != nil {
			return err
		}
		log.Println("parsed template:", path)

	}
	log.Printf("\t%s:%s", s.template.Name(), s.template.DefinedTemplates())
	return nil
}

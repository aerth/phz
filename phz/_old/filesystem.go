package phz

import (
	"log"
	"strings"

	"github.com/fsnotify/fsnotify"
)

// AddWatcher handles changes in the filesystem, reloading templates if needed ( WIP)
func (s *Server) AddWatcher(w *fsnotify.Watcher) {
	go func() {
		for {

			select {
			case event, ok := <-w.Events:
				if !ok {
					log.Println("fsnotify system down (e 101)")
					return
				}
				if containsBadWords(event.Name) || strings.HasSuffix(event.Name, "~") {
					log.Println("Skipping reload:", event.Name)
					continue
				}
				if event.Op&fsnotify.Remove == fsnotify.Remove {
					templatename := strings.TrimPrefix(event.Name, s.config.TemplatePath+"/")
					log.Println("deleting template:", templatename)
					//delete(s.templates, templatename)
					continue
				}
				if event.Op&fsnotify.Write == fsnotify.Write || event.Op&fsnotify.Create == fsnotify.Create || event.Op&fsnotify.Chmod == fsnotify.Chmod {
					templatename := strings.TrimPrefix(event.Name, s.config.TemplatePath+"/")
					log.Println("template modified, reloading:", templatename)
					if err := s.reloadtemplate(templatename); err != nil {
						log.Println("Error reloading template after modification:", err)
						continue
					}
				} else {
					log.Println("new, unhandled fsnotify event:", event)
					continue
				}
			case err, ok := <-w.Errors:
				if !ok {
					log.Println("fsnotify system down (e 100)")
					return
				}
				log.Println("fsnotify going down,  error:", err)
				return
			}

		}

	}()
}

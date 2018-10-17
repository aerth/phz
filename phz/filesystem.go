package phz

import (
	"log"
	"strings"

	"github.com/fsnotify/fsnotify"
)

func (s *Server) AddWatcher(w *fsnotify.Watcher) {
	go func() {
		for {

			select {
			case event, ok := <-w.Events:
				if !ok {
					log.Println("Watcher ded")
					return
				}
				if event.Op&fsnotify.Write == fsnotify.Write {
					if containsBadWords(event.Name) {
						log.Println("Skipping reload:", event.Name)
						continue
					}
					log.Println("Modified file:", event.Name)
					if err := s.reloadtemplate(strings.TrimPrefix(event.Name, s.config.TemplatePath+"/")); err != nil {

						log.Println("Error reloading template after modification:", err)
					}
				}
			case err, ok := <-w.Errors:
				if !ok {
					log.Println("watcher ded?")
					return
				}
				log.Println("fsnotify:", err)
			}

		}

	}()
}

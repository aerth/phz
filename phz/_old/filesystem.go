// Copyright (c) 2018 aerth. All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are
// met:
//
//    * Redistributions of source code must retain the above copyright
// notice, this list of conditions and the following disclaimer.
//    * Redistributions in binary form must reproduce the above
// copyright notice, this list of conditions and the following disclaimer
// in the documentation and/or other materials provided with the
// distribution.
//    * Neither the name of aerth nor the names of this project's
// contributors may be used to endorse or promote products derived from
// this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
// "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
// LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
// A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
// OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
// SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
// LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
// DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
// THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

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

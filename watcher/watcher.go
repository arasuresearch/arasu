// --
// The BSD License (BSD)

// Copyright (c) 2015 Thaniyarasu Kannusamy <thaniyarasu@gmail.com> & Arasu Research Lab Pvt Ltd. All rights reserved.

// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are
// met:

//    * Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.
//    * Redistributions in binary form must reproduce the above copyright notice, this list of
//    conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.
//    * Neither Thaniyarasu Kannusamy <thaniyarasu@gmail.com>. nor ArasuResearch Inc may be used to endorse or promote products derived from this software without specific prior written permission.

// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND AUTHOR
// "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
// LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
// A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
// OWNER OR AUTHOR BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
// SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
// LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
// DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
// THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
// ++

package watcher

import (
	"github.com/arasuresearch/arasu/lib"
	"github.com/howeyc/fsnotify"
	"log"
	"os"
	"path/filepath"
)

type Watcher struct {
	Dir     string
	Dirs    []string
	Pkgs    lib.AS
	Changes lib.AS
	Watcher *fsnotify.Watcher
}

func (w *Watcher) Clean() {
	w.Changes = nil
}
func (w *Watcher) Close() {
	w.Watcher.Close()
}
func (w *Watcher) NoChange() bool {
	return len(w.Changes) == 0
}
func New(dir string) (*Watcher, error) {

	watcher := Watcher{}
	if w, err := fsnotify.NewWatcher(); err == nil {
		go func() {
			for {
				select {
				case ev := <-w.Event:
					fn := ev.Name
					if filepath.Ext(fn) == ".go" {
						watcher.Changes.Add(fn)
						fd, _ := filepath.Split(fn)
						fd = filepath.Clean(fd)
						watcher.Pkgs.Add(fd)
					}

					if ev.IsCreate() {
						watcher.Dirs = append(watcher.Dirs, fn)
						if err := w.Watch(fn); err != nil {
							log.Fatal(err)
						}
					}
					if ev.IsRename() {
						for i, e := range watcher.Dirs {
							if e == fn {
								dirs := watcher.Dirs
								dirs[i] = dirs[len(dirs)-1]
								dirs = dirs[0 : len(dirs)-1]
								watcher.Dirs = dirs

								break
							}
						}
						if err := w.RemoveWatch(fn); err != nil {
							//log.Fatal(err)
							//fmt.Println("remove watching", fn)
						}
					}

				case err := <-w.Error:
					log.Println("error:", err)
				}
			}
		}()

		filepath.Walk(dir, func(src string, info os.FileInfo, err error) error {
			if info.IsDir() {
				watcher.Dirs = append(watcher.Dirs, src)
				if err := w.Watch(src); err != nil {
					log.Fatal(src)
				}
			}
			return nil
		})
		watcher.Watcher = w
		watcher.Dir = dir

		return &watcher, nil
	} else {
		return &watcher, err
	}
}

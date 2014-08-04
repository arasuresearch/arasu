// --
// Copyright (c) 2014 Thaniyarasu Kannusamy <thaniyarasu@gmail.com>.

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.
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

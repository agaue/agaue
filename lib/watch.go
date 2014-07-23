package lib

import (
	"log"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/fsnotify.v0"
)

const (
	eventDelay = 10 * time.Second
)

func launchWatcher() *fsnotify.Watcher {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal("FATAL", err)
	}
	go watch(w)
	// watch posts directory
	if err = w.Add(PostsDir); err != nil {
		w.Close()
		log.Fatal("FATAL", err)
	}
	// watch templates directory
	if err = w.Add(TemplatesDir); err != nil {
		w.Close()
		log.Fatal("FATAL", err)
	}
	return w
}

// TODO: finish the comment
func watch(w *fsnotify.Watcher) {
	var delay <-chan time.Time
	for {
		select {

		case ev := <-w.Events:
			ext := filepath.Ext(ev.Name)
			if strings.HasPrefix(ev.Name, PostsDir) && ext == ".md" {
				delay = time.After(eventDelay)
			} else if strings.HasPrefix(ev.Name, TemplatesDir) && ext == ".html" {
				delay = time.After(eventDelay)
			}

		case err := <-w.Errors:
			log.Println("Watch Error: ", err)

		case <-delay:
			if err := GenerateSite(); err != nil {
				log.Println("Error generating site: ", err)
			} else {
				log.Println("site generated")
			}

		}
	}
}

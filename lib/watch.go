package lib

import (
	"log"
	"path/filepath"
	"strings"
	"time"

	"github.com/howeyc/fsnotify"
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
}

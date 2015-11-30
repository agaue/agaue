package lib

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func ServeBlog() {
	defer launchWatcher().Close()
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatal("FATAL", err)
	}
	fmt.Printf("Agaue is running at http://localhost:%d  Press Ctrl+C to stop.", getConfig(filepath.Join(pwd, "config.json")).port)
	panic(http.ListenAndServe(fmt.Sprintf(":%d", getConfig(filepath.Join(pwd, "config.json")).port), http.FileServer(http.Dir(publicDir))))
}

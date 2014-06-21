package lib

import (
	"fmt"
	"net/http"
)

func ServeBlog() {
	defer launchWatcher().Close()
	fmt.Printf("Agaue is running at http://localhost:%d  Press Ctrl+C to stop.", config.Port)
	panic(http.ListenAndServe(fmt.Sprintf(":%d", config.Port), http.FileServer(http.Dir(PublicDir))))
}

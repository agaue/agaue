package lib

import (
	"net/http"
)

func ServeBlog() {
	defer launchWatcher().Close()
	panic(http.ListenAndServe(":4000", http.FileServer(http.Dir(PublicDir))))
}

package lib

import (
	"net/http"
)

func ServeBlog() {
	panic(http.ListenAndServe(":4000", http.FileServer(http.Dir(PublicDir))))
}

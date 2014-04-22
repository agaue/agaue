package lib

import (
	"fmt"
)

type Site struct {
	Categories []string
	Tags       []string
	PageList   PageSlice
}

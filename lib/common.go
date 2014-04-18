package lib

import (
	"fmt"
	"os"
)

func printError(str string, bigBong ...interface{}) {
	fmt.Fprintln(os.Stderr, str, bigBong)
}

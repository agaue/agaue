package lib

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

func printError(str string, bigBong ...interface{}) {
	fmt.Fprintln(os.Stderr, str, bigBong)
}

func slug(title string) (slug string) {
	//TODO: remove date from final slug
	re, _ := regexp.Compile(`[^\w\s-]`)
	slug = re.ReplaceAllLiteralString(title, "")

	re, _ = regexp.Compile(`[-\s]+`)
	slug = re.ReplaceAllLiteralString(slug, "-")

	slug = strings.ToLower(slug)
	return
}

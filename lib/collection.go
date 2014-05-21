package lib

import (
	"fmt"
	. "github.com/ahmetalpbalkan/go-linq"
)

func getCollection(allPosts []*LongPost) {
	getCategory := func(in T) (T, error) {
		return in.(*LongPost).Category, nil
	}

	res, err := From(allPosts).Select(getCategory).Results()
	if err != nil {
		fmt.Errorf("Error", err)
	} else {
		fmt.Println(res)
	}
}

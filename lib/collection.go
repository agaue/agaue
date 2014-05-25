package lib

import (
	"fmt"
	. "github.com/ahmetalpbalkan/go-linq"
)

func getCollection(allPosts []*LongPost) {
	// getCategory := func(in T) (T, error) {
	// 	return in.(*LongPost).Category, nil
	// }

	res, err := From(allPosts).GroupBy(func(post T) T { return post.(*LongPost).Category }, func(post T) T { return post.(*LongPost).Slug })
	if err != nil {
		fmt.Errorf("Error", err)
	} else {
		fmt.Println(res)
		for key, value := range res {
			fmt.Println(key, value)
		}
	}
}

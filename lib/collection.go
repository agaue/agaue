package lib

import (
	"fmt"
	. "github.com/ahmetalpbalkan/go-linq"
)

type T interface{}

func getCollection(allPosts []*LongPost) map[string][]string {
	res, err := From(allPosts).GroupBy(func(post T) T { return post.(*LongPost).Category }, func(post T) T { return post.(*LongPost).Slug })
	if err != nil {
		fmt.Errorf("Error", err)
	} else {
		fmt.Println(res)
		for key, value := range res {
			fmt.Println(key, value)
		}
		return res
	}
}

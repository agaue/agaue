package lib

import (
	"fmt"
	. "github.com/ahmetalpbalkan/go-linq"
)

// type Collection struct {
// 	Collections map[string][]string
// }

func getCollection(allPosts []*LongPost) map[string][]string {
	collection, err := From(allPosts).GroupBy(func(post T) T { return post.(*LongPost).Category }, func(post T) T { return post.(*LongPost).Slug })
	if err != nil {
		fmt.Errorf("Error", err)
	} else {
		fmt.Println(collection)
		for key, value := range collection {
			fmt.Println(key, value)
		}
		return collection.(map[string][]string)
	}
}

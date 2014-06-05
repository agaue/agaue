package lib

import (
	"fmt"
	. "github.com/ahmetalpbalkan/go-linq"
)

func getCollection(allPosts []*LongPost) map[string][]*LongPost {
	collection, err := From(allPosts).GroupBy(func(post T) T { return post.(*LongPost).Category }, func(post T) T { return post.(*LongPost) })
	if err != nil {
		fmt.Errorf("Error", err)
	}
	return coverMapType(collection)
}

func coverMapType(in map[T][]T) (out map[string][]*LongPost) {
	out = make(map[string][]*LongPost, len(in))
	for k, _ := range in {
		if key, ok := k.(string); ok {
			v := in[k]
			out[key] = make([]*LongPost, 0, len(v))
			for i := range v {
				if value, ok := v[i].(*LongPost); ok {
					out[key] = append(out[key], value)
				}
			}
		}
	}
	return
}

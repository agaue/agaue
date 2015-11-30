package lib

import (
	"fmt"
)

func getCollection(allPosts []longPost) map[string][]longPost {
	collection, err := groupBy(allPosts, func(p longPost) string { return p.Category })
	if err != nil {
		fmt.Println("Error: ", err)
		return nil
	}
	return collection
}

func groupBy(
	posts []longPost,
	keySelector func(c longPost) string,
) (map[string][]longPost, error) {
	var results = make(map[string][]longPost)
	for _, post := range posts {
		key := keySelector(post)
		results[key] = append(results[key], post)
	}
	return results, nil
}

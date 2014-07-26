package lib

import (
	"fmt"
)

func getCollection(allPosts []LongPost) map[string][]LongPost {
	collection, err := groupBy(allPosts, func(p LongPost) string { return p.Category })
	if err != nil {
		fmt.Errorf("Error", err)
		return nil
	}
	return collection
}

func groupBy(
	posts []LongPost,
	keySelector func(c LongPost) string,
) (map[string][]LongPost, error) {
	var results = make(map[string][]LongPost)
	for _, post := range posts {
		key := keySelector(post)
		results[key] = append(results[key], post)
	}
	return results, nil
}

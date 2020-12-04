package main

import (
	"fmt"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

func main() {
	repo, err := git.PlainOpen(".")
	if err != nil {
		panic(err)
	}

	tagsIter, err := repo.Tags()
	if err != nil {
		panic(err)
	}

	tags := []*plumbing.Reference{}

	if tagsIter.ForEach(func(t *plumbing.Reference) error {
		tags = append(tags, t)
		return nil
	}) != nil {
		panic(err)
	}

	fmt.Println(tags)
}

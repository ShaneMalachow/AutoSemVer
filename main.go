package main

import (
	"fmt"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	semver "github.com/shanemalachow/AutoSemVer/SemVer"
)

func main() {
	repo, err := git.PlainOpen(".")
	if err != nil {
		panic(err)
	}

	commits, err := repo.Log(&git.LogOptions{})
	if err != nil {
		panic(err)
	}

	versions := []semver.SemVer{}

	fmt.Println("Trying to find tags in order")
	if commits.ForEach(func(c *object.Commit) error {
		// fmt.Println("Commit: " + c.Hash.String())
		tagsIter, err := repo.Tags()
		err = tagsIter.ForEach(func(t *plumbing.Reference) error {
			// fmt.Println(t)
			tag, err := repo.TagObject(t.Hash())
			if err != nil {
				return err
			}
			tagCommit, err := tag.Commit()
			if err != nil {
				return err
			}
			if tagCommit.Hash.String() == c.Hash.String() {
				ver, err := semver.ParseSemver(strings.Split(t.Name().String(), "/")[2])
				if err != nil {
					return err
				}
				versions = append(versions, ver)
				return nil
			}
			return nil
		})
		return err
	}) != nil {
		panic(err)
	}
	versions[0].Print()
}

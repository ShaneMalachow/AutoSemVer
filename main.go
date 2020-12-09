package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	semver "github.com/shanemalachow/AutoSemVer/SemVer"
)

const commitRegex = `^((?P<type>\w*)(\((?P<area>.+)\))?: ?)?(?P<msg>.*)`

func main() {
	r := regexp.MustCompile(commitRegex)
	major := append([]string{"major", "maj"}, strings.Split(os.Getenv("MAJOR"), ", ")...)
	minor := append([]string{"minor", "min"}, strings.Split(os.Getenv("MINOR"), ", ")...)
	patch := append([]string{"patch", "pat"}, strings.Split(os.Getenv("PATCH"), ", ")...)

	repo, err := git.PlainOpen(".")
	if err != nil {
		panic(err)
	}

	commits, err := repo.Log(&git.LogOptions{})
	if err != nil {
		panic(err)
	}

	versions := []semver.SemVer{}
	commitMsgs := []string{}

	// fmt.Println("Trying to find tags in order")
	if commits.ForEach(func(c *object.Commit) error {
		// fmt.Println("Commit: " + c.Hash.String())
		commitMsgs = append(commitMsgs, c.Message)
		tagsIter, err := repo.Tags()
		if err != nil {
			panic(err)
		}
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

	highestChange := "noop"
	for _, s := range commitMsgs {
		groups := r.FindStringSubmatch(s)
		if groups[2] != "" {
			for _, maj := range major {
				if groups[2] == maj {
					highestChange = "major"
				}
			}
			if highestChange != "major" {
				for _, min := range minor {
					if groups[2] == min {
						highestChange = "minor"
					}
				}
			} else if highestChange != "minor" {
				for _, pat := range patch {
					if groups[2] == pat {
						highestChange = "patch"
					}
				}
			}

		}
	}
	nextVer := versions[0]
	// fmt.Println(nextVer.Version())
	switch highestChange {
	case "major":
		nextVer.Major++
		nextVer.Minor = 0
		nextVer.Patch = 0
	case "minor":
		nextVer.Minor++
		nextVer.Patch = 0
	case "patch":
		nextVer.Patch++
	}
	fmt.Print(nextVer.Version())
}

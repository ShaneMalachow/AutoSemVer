package cmd

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/ShaneMalachow/AutoSemVer/semver"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string
var prerelease string
var buildmetadata string

const commitRegex = `^((?P<type>\w*)(\((?P<area>.+)\))?: ?)?(?P<msg>.*)`

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "semver",
	Short: "AutoSemVer (or simply just semver) allows users to automatically determine the next semantic version based on commit tags.",
	Long:  `AutoSemVer (or simply just semver) allows users to automatically determine the next semantic version based on commit tags.`,
	Run:   determineVersion,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.AutoSemVer.yaml)")

	rootCmd.Flags().StringVarP(&prerelease, "prerelease", "p", "", "Set a pre-release version")
	rootCmd.Flags().StringVarP(&buildmetadata, "build-metadata", "b", "", "Sets build metadata tags")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		viper.AddConfigPath(home)
		viper.SetConfigName(".semver")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func determineVersion(cmd *cobra.Command, args []string) {
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

	version := semver.SemVer{}
	found := false
	commitMsgs := []string{}

	// fmt.Println("Trying to find tags in order")
	err = commits.ForEach(func(c *object.Commit) error {
		// fmt.Println("Commit: " + c.Hash.String())
		if !found {
			commitMsgs = append(commitMsgs, c.Message)
			tagsIter, err := repo.Tags()
			if err != nil {
				panic(err)
			}
			err = tagsIter.ForEach(func(t *plumbing.Reference) error {
				// fmt.Println(t)
				// tag, err := repo.TagObject(t.Hash())
				// if err != nil {
				// 	panic(err)
				// }
				// tagCommit, err := tag.Commit()
				// if err != nil {
				// 	panic(err)
				// }
				if t.Hash().String() == c.Hash.String() {
					version, err = semver.ParseSemver(strings.Split(t.Name().String(), "/")[2])
					if err != nil {
						panic(err)
					}
					found = true
					return nil
				}
				return nil
			})
			return err
		}
		return nil
	})
	if err != nil {
		panic(err)
	}

	highestChange := "noop"
	// fmt.Printf("Major:%v\nMinor:%v\nPatch:%v\n", major, minor, patch)
	for _, s := range commitMsgs {
		groups := r.FindStringSubmatch(s)
		commitType := groups[2]
		if commitType != "" {
			switch {
			case checkSlice(major, commitType):
				highestChange = "major"
			case checkSlice(minor, commitType) && highestChange != "major":
				highestChange = "minor"
			case checkSlice(patch, commitType) && highestChange != "major" && highestChange != "minor":
				highestChange = "patch"
			}
		}
	}
	// fmt.Println(version.Version())
	// fmt.Println(highestChange)
	nextVer := version
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
	if prerelease != "" {
		fmt.Print()
	}
}

func checkSlice(words []string, s string) bool {
	for _, word := range words {
		if s == word {
			return true
		}
	}
	return false
}

package semver_test

import (
	"encoding/json"
	"fmt"
	"testing"

	semver "github.com/shanemalachow/AutoSemVer/SemVer"
)

type SemVerDataset struct {
	str  string
	data semver.SemVer
}

var goodDataset = [...]SemVerDataset{
	{"1.2.3", semver.SemVer{
		Major:         1,
		Minor:         2,
		Patch:         3,
		Prerelease:    "",
		BuildMetadata: "",
	}},
	{"1.1.2-prerelease+meta", semver.SemVer{
		Major:         1,
		Minor:         1,
		Patch:         2,
		Prerelease:    "prerelease",
		BuildMetadata: "meta",
	}},
	{"1.0.0-alpha-a.b-c-somethinglong+build.1-aef.1-its-okay", semver.SemVer{
		Major:         1,
		Minor:         0,
		Patch:         0,
		Prerelease:    "alpha-a.b-c-somethinglong",
		BuildMetadata: "build.1-aef.1-its-okay",
	}},
	{"1.1.2+meta", semver.SemVer{
		Major:         1,
		Minor:         1,
		Patch:         2,
		Prerelease:    "",
		BuildMetadata: "meta",
	}},
	{"1.1.2-prerelease-information", semver.SemVer{
		Major:         1,
		Minor:         1,
		Patch:         2,
		Prerelease:    "prerelease-information",
		BuildMetadata: "",
	}},
}

func TestParseSemver(t *testing.T) {
	for _, tt := range goodDataset {
		dataStr, err := json.Marshal(tt.data)
		if err != nil {
			t.Error("error processing test case")
		}
		testname := fmt.Sprintf("%s,%s", tt.str, dataStr)
		t.Run(testname, func(t *testing.T) {
			ans, err := semver.ParseSemver(tt.str)
			if err != nil {
				t.Errorf("got error parsing %s: %s", tt.str, err)
			}
			answerStr, err := json.Marshal(ans)
			if err != nil {
				t.Error("error marshalling test case result: " + err.Error())
			}
			if ans != tt.data {
				t.Errorf("got %s, want %s", answerStr, dataStr)
			}
		})
	}
}

func TestVersionToString(t *testing.T) {
	for _, tt := range goodDataset {
		dataStr, err := json.Marshal(tt.data)
		if err != nil {
			t.Error("error processing test case")
		}
		testname := fmt.Sprintf("%s,%s", dataStr, tt.str)
		t.Run(testname, func(t *testing.T) {
			ans := tt.data.Version()
			if ans != tt.str {
				t.Errorf("got %s, want %s", ans, tt.str)
			}
		})
	}
}

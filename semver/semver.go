package semver

import (
	"fmt"
	"io"
	"regexp"
	"strconv"
)

// SemVer represents a well formed semantic version (https://semver.org) with fields for the 5 major components of a version: major, minor, and patch versions, as well as prerelease and build metadata.
type SemVer struct {
	Major         int    `json:"major"`
	Minor         int    `json:"minor"`
	Patch         int    `json:"patch"`
	Prerelease    string `json:"prerelease"`
	BuildMetadata string `json:"buildmetadata"`
}

const semverRegex string = `^v?(?P<major>0|[1-9]\d*)\.(?P<minor>0|[1-9]\d*)\.(?P<patch>0|[1-9]\d*)(?:-(?P<prerelease>(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\+(?P<buildmetadata>[0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))?$`

// ParseSemver takes in a raw string and converts it to a SemVer struct for easy access to individual components
func ParseSemver(s string) (SemVer, error) {
	r := regexp.MustCompile(semverRegex)
	match := r.FindStringSubmatch(s)
	maj, err := strconv.Atoi(match[1])
	if err != nil {
		return SemVer{}, fmt.Errorf("cannot convert major version %s to integer", match[1])
	}
	min, err := strconv.Atoi(match[2])
	if err != nil {
		return SemVer{}, fmt.Errorf("cannot convert minor version %s to integer", match[2])
	}
	pat, err := strconv.Atoi(match[3])
	if err != nil {
		return SemVer{}, fmt.Errorf("cannot convert patch version %s to integer", match[3])
	}
	return SemVer{
		Major:         maj,
		Minor:         min,
		Patch:         pat,
		Prerelease:    match[4],
		BuildMetadata: match[5],
	}, nil
}

// Print prints an output of a semantic version to stdout
func (ver SemVer) Print() (int, error) {
	sum := 0
	n, err := fmt.Printf("Major: %v\n", ver.Major)
	sum += n
	if err != nil {
		return sum, err
	}
	n, err = fmt.Printf("Minor: %v\n", ver.Minor)
	sum += n
	if err != nil {
		return sum, err
	}
	n, err = fmt.Printf("Patch: %v\n", ver.Patch)
	sum += n
	if err != nil {
		return sum, err
	}
	n, err = fmt.Printf("Prerelease: %v\n", ver.Prerelease)
	sum += n
	if err != nil {
		return sum, err
	}
	n, err = fmt.Printf("Build Metadata: %v\n", ver.BuildMetadata)
	sum += n
	if err != nil {
		return sum, err
	}
	return sum, nil
}

// Version returns a string representation of the semantic version
func (ver SemVer) Version() string {
	s := fmt.Sprintf("%v.%v.%v", ver.Major, ver.Minor, ver.Patch)
	if ver.Prerelease != "" {
		s += fmt.Sprintf("-%v", ver.Prerelease)
	}
	if ver.BuildMetadata != "" {
		s += fmt.Sprintf("+%v", ver.BuildMetadata)
	}
	return s
}

// Fprint prints an output of a semantic version to a specific io.WriteCloser
func (ver SemVer) Fprint(f io.WriteCloser) (int, error) {
	sum := 0
	n, err := fmt.Fprintf(f, "Major: %v\n", ver.Major)
	sum += n
	if err != nil {
		return sum, err
	}
	n, err = fmt.Fprintf(f, "Minor: %v\n", ver.Minor)
	sum += n
	if err != nil {
		return sum, err
	}
	n, err = fmt.Fprintf(f, "Patch: %v\n", ver.Patch)
	sum += n
	if err != nil {
		return sum, err
	}
	n, err = fmt.Fprintf(f, "Prerelease: %v\n", ver.Prerelease)
	sum += n
	if err != nil {
		return sum, err
	}
	n, err = fmt.Fprintf(f, "Build Metadata: %v\n", ver.BuildMetadata)
	sum += n
	if err != nil {
		return sum, err
	}
	return sum, nil
}

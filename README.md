# AutoSemVer

![Go Build](https://github.com/ShaneMalachow/AutoSemVer/workflows/Go%20Build/badge.svg)

Automatically check Git history to determine next SemVer increment

This library provides methods to automatically determine [semantic version](https://semver.org) increments by processing commits in a modified manner to [Conventional Commits](https://www.conventionalcommits.org/).

## How does it work?

When used within a git repository, it will iterate through the commits until it finds a tag. It will check every commit between the current HEAD and that tag, looking for commits that match the Conventional Commit format (ie. messages starting with `type(function):message`) and calculate whether to increment a major, minor, or patch version, or none at all. By default, the types of commits are `MAJOR`, `MINOR`, and `PATCH`, with any other types not affecting the version calculation. You can pass in environment variables `MAJOR`, `MINOR`, and `PATCH` to include other keywords for the commit types as a comma separated list.

## How can I use it?

Currently, there are binaries available in the releases page as well as a Docker image available. In order to run the Docker image, you need to pass in the working directory of your git repository (the directory containing your .git directory) as a volume mount to `/workdir`.

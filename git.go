package main

import (
	"fmt"
	"os/exec"
)

func git(dir string, args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	if *verbose {
		fmt.Println(cmd.String())
	}
	msg, err := cmd.CombinedOutput()
	cmd.Run()
	return string(msg), err
}

// fetchGitRepository will fetch latest commits and tags
func fetchGitRepository() error {
	_, err := git(*source, "fetch", "--all")
	if err != nil {
		return err
	}
	return nil
}

// getGitTags will query all git tags with created time and subject
func getGitTags() (string, error) {
	tags, err := git(*source, "tag", "-n", "-l", "--sort=creatordate", "--format", "%(refname:short);%(creatordate:short);%(subject)")
	if err != nil {
		return "", err
	}
	return tags, nil
}

// getGitLogs will query commit records between two tags
func getGitLogs(tag1, tag2 string) (string, error) {
	commits, err := git(*source, "log", "--no-merges", "--format=oneline", fmt.Sprintf("%s..%s", tag1, tag2))
	if err != nil {
		return "", err
	}
	return commits, nil
}

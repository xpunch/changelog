package main

import (
	"bytes"
	"fmt"
	"os/exec"
)

func git(dir string, args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Dir = dir
	if *verbose {
		fmt.Println(cmd.String())
	}
	err := cmd.Run()
	return out.String(), err
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
	var notation string
	if len(tag1) > 0 && len(tag2) > 0 {
		notation = fmt.Sprintf("%s..%s", tag1, tag2)
	} else if len(tag1) > 0 {
		notation = tag1
	} else if len(tag2) > 0 {
		notation = tag2
	}
	commits, err := git(*source, "log", "--no-merges", "--format=oneline", notation)
	if err != nil {
		return "", err
	}
	return commits, nil
}

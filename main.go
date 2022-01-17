package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"

	"github.com/hashicorp/go-version"
)

var (
	source  = flag.String("source", "", "--source ~/tmp")
	output  = flag.String("output", "CHANGELOG.md", "--output CHANGELOG.md")
	verbose = flag.Bool("verbose", false, "--verbose")
)

type record struct {
	Version string
	Date    string
	Commits []string
}

func main() {
	flag.Parse()
	if err := fetchGitRepository(); err != nil {
		panic(err)
	}
	gittags, err := getGitTags()
	if err != nil {
		panic(err)
	}
	if *verbose {
		fmt.Println(gittags)
	}
	tags := strings.Split(gittags, "\n")
	records := make([]record, 0, len(tags))
	for _, t := range tags {
		segs := strings.Split(t, ";")
		if len(segs) < 3 || len(segs[0]) == 0 {
			continue
		}
		version := segs[0]
		date := strings.ReplaceAll(segs[1], "-", "/")
		records = append(records, record{Version: version, Date: date})
	}
	sort.Slice(records, func(i, j int) bool {
		s, t := records[i].Version, records[j].Version
		vs, err := version.NewVersion(s)
		if err != nil {
			fmt.Printf("Invalid go version: %s\n", s)
			return s < t
		}
		vt, err := version.NewVersion(t)
		if err != nil {
			fmt.Printf("Invalid go version: %s\n", t)
			return s < t
		}
		return vs.LessThan(vt)
	})
	for i := 0; i < len(records); i++ {
		var v1, v2 string
		if i > 0 {
			v1 = records[i-1].Version
		}
		v2 = records[i].Version
		gitcommits, err := getGitLogs(v1, v2)
		if err != nil {
			fmt.Println(err)
			continue
		}
		if *verbose {
			fmt.Println(string(gitcommits))
		}
		commits := strings.Split(string(gitcommits), "\n")
		records[i].Commits = make([]string, 0, len(commits))
		for _, c := range commits {
			cs := strings.SplitN(c, " ", 2)
			if len(cs) > 1 {
				records[i].Commits = append(records[i].Commits, cs[1])
			} else {
				records[i].Commits = append(records[i].Commits, c)
			}
		}
	}
	var buf bytes.Buffer
	for i := len(records) - 1; i >= 0; i-- {
		r := records[i]
		buf.WriteString(fmt.Sprintf("# %s (%s)\n\n", r.Version, r.Date))
		for _, c := range r.Commits {
			if len(c) > 0 {
				buf.WriteString(fmt.Sprintf("- %s\n", c))
			}
		}
	}
	if err := ioutil.WriteFile(*output, buf.Bytes(), os.ModePerm); err != nil {
		panic(err)
	}
}

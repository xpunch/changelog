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
	fetch   = flag.Bool("fetch", false, "--fetch")
	verbose = flag.Bool("verbose", false, "--verbose")
)

type record struct {
	Version  string
	Date     string
	Features []string
	BugFixes []string
}

func main() {
	flag.Parse()
	if *fetch {
		if err := fetchGitRepository(); err != nil {
			fmt.Println(err)
			return
		}
	}
	gittags, err := getGitTags()
	if err != nil {
		fmt.Println(err)
		return
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
		records[i].Features = make([]string, 0)
		records[i].BugFixes = make([]string, 0)
		for _, c := range commits {
			msg, cs := c, strings.SplitN(c, " ", 2)
			if len(cs) > 1 {
				msg = cs[1]
			}
			if len(msg) == 0 {
				continue
			}
			if strings.HasPrefix(msg, "[fix]") {
				records[i].BugFixes = append(records[i].BugFixes, strings.TrimSpace(strings.TrimPrefix(msg, "[fix]")))
			} else if strings.HasPrefix(msg, "[feat]") {
				records[i].Features = append(records[i].Features, strings.TrimSpace(strings.TrimPrefix(msg, "[feat]")))
			} else if strings.HasPrefix(msg, "[feature]") {
				records[i].Features = append(records[i].Features, strings.TrimSpace(strings.TrimPrefix(msg, "[feature]")))
			} else if strings.Contains(msg, "fix") {
				records[i].BugFixes = append(records[i].BugFixes, msg)
			} else {
				records[i].Features = append(records[i].Features, msg)
			}
		}
	}
	var buf bytes.Buffer
	for i := len(records) - 1; i >= 0; i-- {
		r := records[i]
		buf.WriteString(fmt.Sprintf("# %s (%s)\n", strings.TrimPrefix(r.Version, "v"), r.Date))
		if len(r.Features) > 0 {
			buf.WriteString("\n### Features\n\n")
			for _, c := range r.Features {
				if len(c) > 0 {
					buf.WriteString(fmt.Sprintf("- %s\n", c))
				}
			}
		}
		if len(r.BugFixes) > 0 {
			buf.WriteString("\n### Bug Fixes\n\n")
			for _, c := range r.BugFixes {
				if len(c) > 0 {
					buf.WriteString(fmt.Sprintf("- %s\n", c))
				}
			}
		}
		if i != 0 {
			buf.WriteString("\n")
		}
	}
	if err := ioutil.WriteFile(*output, buf.Bytes(), os.ModePerm); err != nil {
		fmt.Println(err)
		return
	}
}

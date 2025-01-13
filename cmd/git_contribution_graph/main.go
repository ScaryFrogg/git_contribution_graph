package main

import (
	"flag"
	"github.com/ScaryFrogg/git_contribution_graph/internal/draw"
	"github.com/ScaryFrogg/git_contribution_graph/internal/git"
	"github.com/ScaryFrogg/git_contribution_graph/internal/github"
	"time"
)

func main() {
	token := flag.String("token", "", "GitHub token")
	username := flag.String("username", "", "GitHub username")
	fromFlag := flag.String("from", "", "Begin Time for the graph in ISO-8601 format, defaults to the beggining of the current year")
	toFlag := flag.String("to", "", "End Time for the graph in ISO-8601 format, defaults to the current time")
	flag.Parse()

	parseDates(fromFlag, toFlag)

	var contributionMap [][]int
	if *token == "" {
		var err bool
		contributionMap, err = git.GetLocalContributions(*fromFlag, *toFlag)
		if err {
			return
		}
	} else {
		contributionMap = github.FetchContributions(*username, *token, *fromFlag, *toFlag)
	}

	draw.DrawGrid(contributionMap)
}

func parseDates(from *string, to *string) {
	now := time.Now()
	if *from == "" {
		*from = time.Date(now.Year(), 1, 1, 0, 0, 0, 0, time.Local).Format(time.RFC3339)
	}
	if *to == "" {
		*to = now.Format(time.RFC3339)
	}
}

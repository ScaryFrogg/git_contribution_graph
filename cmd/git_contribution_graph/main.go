package main

import (
	"flag"
	"github.com/ScaryFrogg/git_contribution_graph/internal/draw"
	"github.com/ScaryFrogg/git_contribution_graph/internal/github"
)

func main() {
	token := flag.String("token", "", "GitHub token")
	username := flag.String("username", "", "GitHub username")
	fromFlag := flag.String("from", "", "Begin Time for the graph in ISO-8601 format, defaults to the beggining of the current year")
	toFlag := flag.String("to", "", "End Time for the graph in ISO-8601 format, defaults to the current time")
	flag.Parse()

	contributionMap := github.FetchContributions(*username, *token, *fromFlag, *toFlag)
	draw.DrawGrid(contributionMap)
}

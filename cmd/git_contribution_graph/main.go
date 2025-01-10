package main

import (
	"flag"
	"github.com/ScaryFrogg/git_contribution_graph/internal/draw"
	"github.com/ScaryFrogg/git_contribution_graph/internal/github"
)

func main() {
	token := flag.String("token", "", "GitHub token")
	username := flag.String("username", "", "GitHub username")
	flag.Parse()

	contributionMap := github.FetchContributions(*username, *token)
	draw.DrawGrid(contributionMap)
}

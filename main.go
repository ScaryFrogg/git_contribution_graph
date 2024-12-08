package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/shurcooL/graphql"
	"golang.org/x/oauth2"
	"log"
)

const (
	resetColor = "\033[0m"
	white      = "\x1b[48;5;249m"
	c1         = "\x1b[48;5;40m"
	c2         = "\x1b[48;5;34m"
	c3         = "\x1b[48;5;28m"
	c4         = "\x1b[48;5;22m"
)

var colorMap = []string{white, c1, c2, c3, c4}

func drawGrid(activityLevels [][]int) {
	for y := 0; y < len(activityLevels); y++ {
		for x := 0; x < len(activityLevels[y]); x++ {
			var colorIndex int
			cCount := activityLevels[y][x]
			switch {
			case cCount == 0:
				colorIndex = 0
			case cCount < 3:
				colorIndex = 1
			case cCount < 5:
				colorIndex = 2
			default:
				colorIndex = 3
			}
			fmt.Print(colorMap[colorIndex] + "  " + resetColor)
		}
		fmt.Println()
	}
}

func main() {
	token := flag.String("token", "", "GitHub token")
	username := flag.String("username", "", "GitHub username")
	flag.Parse()

	contributionMap := fetchContributions(*username, *token)
	println(len(contributionMap))
	println(len(contributionMap[0]))

	drawGrid(contributionMap)
}
func fetchContributions(username string, token string) [][]int {
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	httpClient := oauth2.NewClient(context.Background(), src)

	client := graphql.NewClient("https://api.github.com/graphql", httpClient)

	// Define the query
	var query Query
	variables := map[string]interface{}{
		"login": graphql.String(username),
	}

	// Execute the query
	err := client.Query(context.Background(), &query, variables)
	if err != nil {
		log.Fatalf("Failed to execute query: %v", err)
	}

	var contributionMap = make([][]int, 7)
	// Print the results
	for _, week := range query.User.ContributionsCollection.ContributionCalendar.Weeks {
		for _, day := range week.ContributionDays {
			contributionMap[int(day.Weekday)] = append(contributionMap[int(day.Weekday)], int(day.ContributionCount))
		}
	}
	return contributionMap
}

type ContributionDay struct {
	Weekday           graphql.Int    `json:"weekday"`
	Date              graphql.String `json:"date"`
	ContributionCount graphql.Int    `json:"contributionCount"`
}

type Week struct {
	ContributionDays []ContributionDay `json:"contributionDays"`
}

type ContributionCalendar struct {
	Weeks []Week `json:"weeks"`
}

type ContributionsCollection struct {
	ContributionCalendar ContributionCalendar `json:"contributionCalendar"`
}

type User struct {
	ContributionsCollection ContributionsCollection `json:"contributionsCollection"`
}

type Query struct {
	User User `graphql:"user(login: $login)"`
}

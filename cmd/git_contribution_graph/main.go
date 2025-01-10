package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"golang.org/x/oauth2"
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
	// OAuth2 client setup
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	httpClient := oauth2.NewClient(context.Background(), src)

	// Build the query
	currentYear := time.Now().Year()
	from := fmt.Sprintf("%d-01-01T00:00:00", currentYear)
	currentTime := time.Now().UTC().Format(time.RFC3339)
	//query ($login: String!, $from: DateTime!, $to: DateTime!) {
	query := `
	query ($login: String!, $to: DateTime!) {
		user(login: $login) {
			contributionsCollection(from: $from, to: $to) {
				contributionCalendar {
					weeks {
						contributionDays {
							weekday
							contributionCount
						}
					}
				}
			}
		}
	}`
	variables := map[string]interface{}{
		"login": username,
		"from":  from,
		"to":    currentTime,
	}
	requestBody := map[string]interface{}{
		"query":     query,
		"variables": variables,
	}
	bodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		log.Fatalf("Failed to marshal request body: %v", err)
	}

	// Send the request
	req, err := http.NewRequest("POST", "https://api.github.com/graphql", bytes.NewBuffer(bodyBytes))
	if err != nil {
		log.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	resp, err := httpClient.Do(req)
	if err != nil {
		log.Fatalf("Failed to execute query: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Unexpected status code: %d", resp.StatusCode)
	}

	// Parse the response
	var result GithubResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Fatalf("Failed to parse response: %v", err)
	}

	if errorMessage := result.Errors; errorMessage != nil {
		log.Fatalln(errorMessage[0])
	}

	// Process the results
	var contributionMap = make([][]int, 7)
	for w, week := range result.Data.User.ContributionsCollection.ContributionCalendar.Weeks {
		fmt.Printf("week:%d", w)
		for _, day := range week.ContributionDays {
			contributionMap[day.Weekday] = append(contributionMap[day.Weekday], day.ContributionCount)
		}
	}
	return contributionMap
}

type GithubResponse struct {
	Data struct {
		User struct {
			ContributionsCollection struct {
				ContributionCalendar struct {
					Weeks []struct {
						ContributionDays []struct {
							Weekday           int `json:"weekday"`
							ContributionCount int `json:"contributionCount"`
						} `json:"contributionDays"`
					} `json:"weeks"`
				} `json:"contributionCalendar"`
			} `json:"contributionsCollection"`
		} `json:"user"`
	} `json:"data"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors"`
}

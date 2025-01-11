package github

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"golang.org/x/oauth2"
)

func FetchContributions(username string, token string) [][]int {
	// OAuth2 client setup
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	httpClient := oauth2.NewClient(context.Background(), src)

	// Build the query
	now := time.Now()
	from := time.Date(now.Year(), 1, 1, 0, 0, 0, 0, time.Local).Format(time.RFC3339)
	currentTime := now.Format(time.RFC3339)
	query := `
	query ($login: String!, $from: DateTime!, $to: DateTime!) {
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

	//Put in the gap when year does not start with first weekday (Sunday)
	var firstWeek = result.Data.User.ContributionsCollection.ContributionCalendar.Weeks[0]
	for gapIndex := 0; gapIndex < 7-len(firstWeek.ContributionDays); gapIndex++ {
		contributionMap[gapIndex] = append(contributionMap[gapIndex], -1)
	}

	//Fill in the rest
	for _, week := range result.Data.User.ContributionsCollection.ContributionCalendar.Weeks {
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

package git

import (
	"fmt"
	"math"
	"os/exec"
	"strings"
	"time"
)

func GetLocalContributions(from string, to string) (matrix [][]int, errorOrNotRepo bool) {

	if !isDirRepo() {
		return nil, true
	}

	gitLog := getGitLog()

	// get days with contributions
	contributionCount := make(map[string]int)
	for _, x := range gitLog {
		contributionCount[x]++
	}
	fmt.Println("contributionCountMap", contributionCount)

	//build matrix
	fromDate, parseError := time.Parse(time.RFC3339, from)
	if parseError != nil {
		//TODO
	}
	days := math.Ceil(time.Now().Sub(fromDate).Hours() / 24)
	contributionMatrix := make([][]int, 7)
	for d := 0; d < int(days); d++ {
		date := fromDate.AddDate(0, 0, d)
		weekday := int(date.Weekday())
		x := date.Format("Jan 2 2006")
		//add gap if selected start day is not first day of the week
		if d == 0 {
			for gapIndex := 0; gapIndex < 6-weekday; gapIndex++ {
				contributionMatrix[gapIndex] = append(contributionMatrix[gapIndex], -1)
			}
		}
		count := 0
		if cc, exists := contributionCount[x]; exists {
			count = cc
		}
		contributionMatrix[weekday] = append(contributionMatrix[weekday], count)
	}

	return contributionMatrix, false
}
func isDirRepo() bool {
	cmd := "git status"
	_, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		return false
	}
	return true
}

func getGitLog() []string {
	//TODO explore advantages of using Pipe and Stdout, Stdin from os stdlib instead
	cmd := "git log --pretty=format:'%cd' | sed -E 's/^.{4}//; s/[[:space:]]\\+.*$//; s/[[:space:]][[:digit:]]{2}:[[:digit:]]{2}:[[:digit:]]{2}//'"

	out, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		println("Failed to execute command: %s", err)
		return make([]string, 0)
	}

	return strings.Split(string(out), "\n")
}

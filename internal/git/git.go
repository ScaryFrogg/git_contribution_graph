package git

import (
	"fmt"
	"log"
	"math"
	"os/exec"
	"strings"
	"time"
)

func GetLocalContributions(from string, to string) (matrix [][]int, errorOrNotRepo bool) {
	if !isDirRepo() {
		return nil, true
	}

	gitLog := getGitLog(from, to)

	// get days with contributions
	contributionCount := make(map[string]int)
	for _, x := range gitLog {
		contributionCount[x]++
	}

	fromDate, parseError := time.Parse(time.RFC3339, from)
	if parseError != nil {
		log.Fatal("Fatal: Unable to parse fromDate")
	}
	toDate, parseError := time.Parse(time.RFC3339, to)
	if parseError != nil {
		log.Fatal("Fatal: Unable to parse toDate")
	}
	//build matrix
	contributionMatrix := make([][]int, 7)
	//add gap if selected start day is not first day of the week
	if d := int(fromDate.Weekday()); d != 0 {
		for gapIndex := 0; gapIndex < d; gapIndex++ {
			contributionMatrix[gapIndex] = append(contributionMatrix[gapIndex], -1)
		}
	}

	days := math.Ceil(toDate.Sub(fromDate).Hours() / 24)
	for d := 0; d < int(days); d++ {
		date := fromDate.AddDate(0, 0, d)
		weekday := int(date.Weekday())
		x := date.Format(time.DateOnly)
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
func getGitLog(from string, to string) []string {
	fromOpt := fmt.Sprintf("--since='%v'", from)
	toOpt := fmt.Sprintf("--until='%v'", to)
	out, err := exec.Command("git", "log", "--pretty=format:%cd", "--date=short", fromOpt, toOpt).Output()
	if err != nil {
		println("Failed to execute command: %s", err)
		return make([]string, 0)
	}
	return strings.Split(string(out), "\n")
}

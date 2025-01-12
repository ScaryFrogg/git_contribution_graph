package git

import (
	"fmt"
	"os/exec"
	"strings"
)

func GetLocalContributions(from string, to string) [][]int {
	gitLog := getGitLog()
	println(len(gitLog))

	// get days with contributions
	contributionCount := make(map[string]int)
	for _, x := range gitLog {
		contributionCount[x]++
	}
	fmt.Println("contributionCountMap", contributionCount)

	//build matrix

	return make([][]int, 7)
}

func getGitLog() []string {
	cmd := "git log --pretty=format:'%cd' | sed -E 's/^.{4}//; s/[[:space:]]\\+.*$//; s/[[:space:]][[:digit:]]{2}:[[:digit:]]{2}:[[:digit:]]{2}//'"

	out, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		println("Failed to execute command: %s", cmd)
		return make([]string, 0)
	}
	return strings.Split(string(out), "\n")
}

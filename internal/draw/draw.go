package draw

import (
	"fmt"
)

const (
	resetColor = "\033[0m"
	white      = "\x1b[48;5;249m"
	c1         = "\x1b[48;5;40m"
	c2         = "\x1b[48;5;34m"
	c3         = "\x1b[48;5;28m"
	c4         = "\x1b[48;5;22m"
)

var colorMap = []string{resetColor, white, c1, c2, c3, c4}

func DrawGrid(activityLevels [][]int) {
	for y := 0; y < len(activityLevels); y++ {
		for x := 0; x < len(activityLevels[y]); x++ {
			var colorIndex int
			cCount := activityLevels[y][x]
			switch {
			case cCount == -1:
				colorIndex = 0
			case cCount == 0:
				colorIndex = 1
			case cCount < 3:
				colorIndex = 2
			case cCount < 5:
				colorIndex = 3
			default:
				colorIndex = 4
			}
			fmt.Print(colorMap[colorIndex] + "  " + resetColor)
		}
		fmt.Println()
	}
}

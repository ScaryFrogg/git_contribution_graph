package draw

import (
	"fmt"
)

const (
	resetColor    = "\033[0m"
	invisibleText = "\033[8;31m"
	white         = "\x1b[38;5;249m"
	c1            = "\x1b[38;5;40m"
	c2            = "\x1b[38;5;34m"
	c3            = "\x1b[38;5;28m"
	c4            = "\x1b[38;5;22m"
)

var colorMap = []string{white, c1, c2, c3, c4}

func DrawGrid(activityLevels [][]int) {
	for y := 0; y < len(activityLevels); y++ {
		for x := 0; x < len(activityLevels[y]); x++ {
			var colorCode string
			cCount := activityLevels[y][x]
			switch {
			case cCount == -1:
				colorCode = invisibleText
			case cCount == 0:
				colorCode = colorMap[0]
			case cCount < 3:
				colorCode = colorMap[1]
			case cCount < 5:
				colorCode = colorMap[2]
			default:
				colorCode = colorMap[3]
			}
			fmt.Print(colorCode + "⬛ " + resetColor)
		}
		fmt.Println()
	}
}

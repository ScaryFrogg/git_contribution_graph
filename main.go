package main

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

var colorMap = []string{white, c1, c2, c3, c4}

func drawGrid(width, height int, activityLevels [][]int) {
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			fmt.Print(colorMap[activityLevels[y][x]] + "  " + resetColor)
		}
		fmt.Println()
	}
}

func main() {
	width := 52
	height := 7
	activityLevels := make([][]int, height)
	for i := range activityLevels {
		activityLevels[i] = make([]int, width)
		for j := range activityLevels[i] {
			activityLevels[i][j] = (i + j) % len(colorMap)
		}
	}

	drawGrid(width, height, activityLevels)
}

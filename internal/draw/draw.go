package draw

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
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

var defaultColorMap = []string{white, c1, c2, c3, c4}

func DrawGrid(activityLevels [][]int, colorSchema string) {
	colorMap := getColorSchema(colorSchema)

	for y := 0; y < len(activityLevels); y++ {
		for x := 0; x < len(activityLevels[y]); x++ {
			var colorCode string
			cCount := activityLevels[y][x]
			switch {
			case cCount == -1:
				fmt.Print("\u3000 ")
				continue
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

func getColorSchema(colorSchema string) []string {
	if colorSchema == "" {
		return defaultColorMap
	}
	re := regexp.MustCompile(`^(\d?\d?\d\,){3}\d?\d?\d$`)
	if !re.MatchString(colorSchema) {
		fmt.Println("color schema not in proper format, falling back to default schema")
		return defaultColorMap
	}

	colorCodes := strings.Split(colorSchema, ",")
	if len(colorCodes) != 4 {
		//Should not come to this because of regexp
		fmt.Println("color schema should contain only 4 color codes, falling back to default schema")
		return defaultColorMap
	}
	colorMap := make([]string, 4)
	for i, color := range colorCodes {
		cc, err := strconv.Atoi(color)
		if err != nil {
			//Should not come to this because of regexp
			fmt.Printf("Not proper value for color code, %v\n", cc)
		}
		if cc < 0 || cc > 255 {
			fmt.Printf("Not proper value for color code, %v values go from 0-255\n", cc)
		}
		colorMap[i] = fmt.Sprintf("\x1b[38;5;%dm", cc)
	}
	return colorMap
}

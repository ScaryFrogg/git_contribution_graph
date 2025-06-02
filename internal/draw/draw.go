package draw

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
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

var (
	defaultColorMap = [5]string{white, c1, c2, c3, c4}
	weekDayNames    = [7]string{"Sun", "Mon", "Tue", "Wen", "Thr", "Fri", "Sat"}
)

func DrawGrid(activityLevels [][]int, colorSchema string, legendEnabled bool) {
	colorMap := getColorSchema(colorSchema)

	for y := 0; y < len(activityLevels); y++ {
		if legendEnabled {
			fmt.Print(weekDayNames[y] + "\u3000")
		}
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
			fmt.Print(colorCode + "\u2588\u2588 " + resetColor)
		}
		fmt.Println()
		fmt.Println()
	}
}

func DrawMonthsLegend(legend bool, from string, intervalLength int) {
	if !legend {
		return
	}
	fromDate, parseError := time.Parse(time.RFC3339, from)
	if parseError != nil {
		return
	}

	//Padding for day of the week text
	fmt.Printf("   \u3000")
	var currMonth = fromDate.Month()
	// Now Align checking to Sunday so we are allways checking first row of squares for next month
	var currDay = int(fromDate.Weekday())
	var firstSunday = fromDate.AddDate(0, 0, -currDay)
	// If next week is not in next month we have space to write down first months name
	if currMonth == firstSunday.AddDate(0, 0, 7).Month() {
		fmt.Printf("%.3s   ", currMonth)
	} else {
		//we skip first cube month will be written in loop later
		fmt.Printf("   ")
	}

	for i := 2; i < intervalLength; i++ {
		if currMonth != firstSunday.AddDate(0, 0, i*7).Month() {
			currMonth = firstSunday.AddDate(0, 0, i*7).Month()
			fmt.Printf("%.3s   ", currMonth)
			i++
		} else {
			fmt.Printf("   ")
		}
	}
	fmt.Println()
}

func getColorSchema(colorSchema string) [5]string {
	if colorSchema == "" {
		return defaultColorMap
	}
	re := regexp.MustCompile(`^(\d?\d?\d\,){4}\d?\d?\d$`)
	if !re.MatchString(colorSchema) {
		fmt.Println("color schema not in proper format, falling back to default schema")
		return defaultColorMap
	}

	colorCodes := strings.Split(colorSchema, ",")
	if len(colorCodes) != 5 {
		//Should not come to this because of regexp
		fmt.Println("color schema should contain only 5 color codes, falling back to default schema")
		return defaultColorMap
	}
	var colorMap [5]string
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

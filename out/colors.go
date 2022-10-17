package out

import (
	"fmt"
	"math"
	"strings"
)

const ResetCode string = "\033[0m"

var Colors = []*Color{
	&Color{"green", "\033[32m", 0},
	&Color{"yellow", "\033[33m", 0},
	&Color{"blue", "\033[34m", 0},
	&Color{"purple", "\033[35m", 0},
	&Color{"cyan", "\033[36m", 0},
	&Color{"white", "\033[37m", 0},
	&Color{"red", "\033[31m", 0},
}

type Color struct {
	Name       string
	Code       string
	usageCount int // used for round robin color assignments
}

// Get color code and increase usage count
func UseColor(colorName string) (string, error) {
	color, err := GetColorByName(colorName)
	if err != nil {
		return "", err
	}

	color.usageCount += 1
	return color.Code, nil
}

func GetColorByName(name string) (*Color, error) {
	lower := strings.ToLower(name)
	for _, color := range Colors {
		if lower == color.Name {
			return color, nil
		}
	}
	return nil, fmt.Errorf("color with name %s not found", name)
}

func UseNextAvailableColorCode() string {
	minIndex := 0
	minUsageCount := math.MaxInt
	for i, color := range Colors {
		if color.usageCount < minUsageCount {
			minIndex = i
			minUsageCount = Colors[i].usageCount
		}
	}
	Colors[minIndex].usageCount += 1
	return Colors[minIndex].Code
}

func ColorText(text, colorCode string) string {
	return colorCode + text + ResetCode
}

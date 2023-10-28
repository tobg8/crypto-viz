package common

// import (
// 	"strconv"
// 	"strings"
// )

// func stringToFloat(s string) float64 {
// 	var value float64
// 	f, err := strconv.ParseFloat(s, 64)
// 	if err != nil {
// 		value = 0
// 	} else {
// 		value = f
// 	}
// 	return value
// }

// func stringToInt(s string) int {
// 	var value int
// 	int, err := strconv.Atoi(s)
// 	if err != nil {
// 		value = 0
// 	} else {
// 		value = int
// 	}
// 	return value
// }

// func splitText(s string, direction int) string {
// 	if s == "" {
// 		return ""
// 	}

// 	// removes spaces and commas
// 	str := strings.ReplaceAll(strings.ReplaceAll(s, " ", ""), ",", ".")

// 	// removes first or last character
// 	if direction == 1 {
// 		str = str[1:]
// 	} else {
// 		str = str[:len(str)-1]
// 	}

// 	return str
// }

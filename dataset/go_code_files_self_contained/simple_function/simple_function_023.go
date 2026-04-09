package main

import (
	"math/rand"
	"strings"
)

func HexColor() string {
	color := make([]byte, 6)
	hashQuestion := []byte("?#")
	for i := 0; i < 6; i++ {
		color[i] = hashQuestion[rand.Intn(2)]
	}

	return "#" + replaceWithLetters(replaceWithNumbers(string(color)))

	// color := ""
	// for i := 1; i <= 6; i++ {
	// 	color += RandString([]string{"?", "#"})
	// }

	// // Replace # with number
	// color = replaceWithNumbers(color)

	// // Replace ? with letter
	// for strings.Count(color, "?") > 0 {
	// 	color = strings.Replace(color, "?", RandString(letters), 1)
	// }

	// return "#" + color
}

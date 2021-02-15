package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {

	str := "sdfdsf, sdpsd[pd, sdosidsoowoow"
	result := strings.Split(str, ",")

	for i := range result {
		fmt.Println(result[i])
		QRgen(result[i], "gen/qr.png")
	}

	// Part 1: open the file and scan it.
	f, _ := os.Open("C:\\programs\\file.txt")
	scanner := bufio.NewScanner(f)

	// Part 2: call Scan in a for-loop.
	for scanner.Scan() {
		line := scanner.Text()

		// Split the line on commas.
		parts := strings.Split(line, ",")

		// Loop over the parts from the string.
		for i := range parts {
			fmt.Println(parts[i])
		}
		// Write a newline.
		fmt.Println()
	}

	QRreader("gen/qr.png")
}

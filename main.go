package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {

	// Part 1: open the file and scan it.
	f, _ := os.Open("file/file.txt")
	scanner := bufio.NewScanner(f)

	// Part 2: call Scan in a for-loop.
	for scanner.Scan() {
		line := scanner.Text()

		// Split the line on commas.
		parts := strings.Split(line, ",")

		// Loop over the parts from the string.
		for i := range parts {
			QRgen(parts[i], "gen/"+fmt.Sprint(i)+".png")
			// QRreader("gen/" + parts[i] + ".png")
			fmt.Println(fmt.Sprint(i) + " " + "---" + parts[i])
		}

	}
	/*
		// Read QRcode image
		qrfiles, err := ioutil.ReadDir("gen/")
		if err != nil {
			log.Fatal(err)
		}

		for _, f := range qrfiles {
			QRreader(f.Name())
			// fmt.Println(f.Name())
		}
	*/
	// QRreader("gen/qr.png")

}

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {

	// Part 1: open the file and scan it.
	f, _ := os.Open("files/readURL.txt")
	scanner := bufio.NewScanner(f)

	// Part 2: call Scan in a for-loop.
	for scanner.Scan() {
		line := scanner.Text()

		// Split the line on commas.
		parts := strings.Split(line, ",")

		// Loop over the parts from the string.
		for i := range parts {
			QRgen(parts[i], "genQR/"+fmt.Sprint(i+1)+".png")
			// QRreader("gen/" + parts[i] + ".png")
			fmt.Println(fmt.Sprint(i+1) + " " + "---" + " " + parts[i])
			CreateTxt(fmt.Sprint(i+1) + ".txt")
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

func CreateTxt(files string) os.File {
	file, err := os.Create(files)

	defer file.Close()

	if err != nil {
		log.Fatal(err)
	}

	/* 	fileName := files

	   	val := "old\nfalcon\nsky\ncup\nforest\n"
	   	data := []byte(val)

	   	err := ioutil.WriteFile(fileName, data, 0644)

	   	if err != nil {
	   		log.Fatal(err)
	   	}

	   	fmt.Println("done")
	*/

	/*
		fileName := "data.txt"

		f, err := os.Create(fileName)

		if err != nil {

		    log.Fatal(err)
		}

		defer f.Close()

		words := []string{"sky", "falcon", "rock", "hawk"}

		for _, word := range words {

		    _, err := f.WriteString(word + "\n")

		    if err != nil {
		        log.Fatal(err)
		    }
		}

		fmt.Println("done")
	*/

	// fmt.Println("file created")
	return *file
}

package HangmanWeb

import (
	"os"
	"log"
	"bufio"
)


func openFile(file string) []string {
	fileLines := []string{}
	// open file
	f, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}
	// remember to close the file at the end of the program
	defer f.Close()

	// read the file line by line using scanner
	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		// do something with a line
		fileLines = append(fileLines, scanner.Text())
	}
	return fileLines
}

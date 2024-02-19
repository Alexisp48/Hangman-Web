package HangmanWeb

import (
	"os"
	"log"
	"bufio"
)


func openFile(file string) []string { // ouvrir un fichier ( pour la liste de mot en .txt )
	fileLines := []string{}
	// open file
	f, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		fileLines = append(fileLines, scanner.Text())
	}
	return fileLines
}

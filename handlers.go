package main

import (
	"bufio"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
    "math/rand"
)

type Player struct {
    Test string
    G *Game
}

type Game struct {
    WordInString string
    WordToFind []string
    WordFind []string
    LetterTested []string
    Win bool
    TryNumber int
}

func (P *Player) Home(w http.ResponseWriter, r *http.Request) {
    P.renderTemplates(w, "home")
    fmt.Fprintf(w, "Hello")
}

func testEq(a, b []string) bool {
    if len(a) != len(b) {
        return false
    }
    for i := range a {
        if a[i] != b[i] {
            return false
        }
    }
    return true
}

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

func (P *Player) Hangman(w http.ResponseWriter, r *http.Request) {
    if len(P.G.WordToFind) == 0 { // initialse le mot
        P.G.WordInString = openFile("worldDataBase.txt")[rand.Intn(len(openFile("worldDataBase.txt")))]
        for i := 0; i < len(P.G.WordInString); i++ {
            P.G.WordToFind = append(P.G.WordToFind, string(P.G.WordInString[i]))
        }
        for i := 0; i < len(P.G.WordToFind); i++ {
            P.G.WordFind = append(P.G.WordFind, "_")
        }
        P.renderTemplates(w, "hangman")
    } else {
        letter := r.FormValue("letter")
        if len(letter) == 1{ // test une lettre
            var letterInWorld = false
            for i := 0; i < len(P.G.WordToFind); i++ {
                if P.G.WordToFind[i] == letter {
                    letterInWorld = true
                    P.G.WordFind[i] = letter
                }
            }
            if !letterInWorld {
                P.G.TryNumber--
            }
            P.G.LetterTested = append(P.G.LetterTested, letter)
        }else if letter == P.G.WordInString { // test un mot et c'est le bon
            P.G.Win = true // win
        } else { // test pas bon lettre ou mot
            P.G.TryNumber--
        }
        if testEq(P.G.WordFind, P.G.WordToFind) {
            P.G.Win = true
        }
        if P.G.Win {
            P.Test = "WINNNNN GG"
        }
        P.renderTemplates(w, "hangman")
    }
}

func (P *Player) renderTemplates(w http.ResponseWriter, tmpl string) {
    t, err := template.ParseFiles("./templates/" + tmpl + ".page.html")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    t.Execute(w, P)
}


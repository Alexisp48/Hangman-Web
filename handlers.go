package main

import (
	"fmt"
	"html/template"
	"net/http"
)

type Player struct {
    Test string
}

var (
    wordInString = ""
    wordToFind = []string{}
    wordFind = []string{}
    letterTested = []string{}
    win = false
    tryNumber = 0
)

func (p *Player) Home(w http.ResponseWriter, r *http.Request) {
    p.renderTemplates(w, "home")
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

func (p *Player) Hangman(w http.ResponseWriter, r *http.Request) {
    p.renderTemplates(w, "hangman")
    if len(wordToFind) == 0 {
        wordToFind = []string{"T", "E", "S", "T"}
        for i := 0; i < len(wordToFind); i++ {
            wordFind = append(wordFind, "_")
        }
    }
    letter := r.FormValue("letter")
    if len(letter) == 1{ // test une lettre
        var letterInWorld = false
        for i := 0; i < len(wordToFind); i++ {
            if wordToFind[i] == letter {
                letterInWorld = true
                wordFind[i] = letter
            }
        }
        if !letterInWorld {
            tryNumber++
        }
        letterTested = append(letterTested, letter)
    }else if letter == wordInString { // test un mot et c'est le bon
        win = true // win
    } else { // test pas bon lettre ou mot
        tryNumber++
    }
    if testEq(wordFind, wordToFind) {
        win = true
    }
    fmt.Fprintf(w, "\n")
    for i := 0; i < len(wordFind); i++ {
        fmt.Fprintf(w, wordFind[i])
    }
    fmt.Fprintf(w, "\n")
    for i := 0; i < len(letterTested); i++ {
        fmt.Fprintf(w, letterTested[i])
    }
    fmt.Fprintf(w, "\n")
    fmt.Fprintln(w, tryNumber)
    if win {
        //fmt.Println(tryNumber)
        wordInString = ""
        wordToFind = []string{}
        wordFind = []string{}
        letterTested = []string{}
        win = false
        tryNumber = 0
    }
}

func (p *Player) renderTemplates(w http.ResponseWriter, tmpl string) {
    t, err := template.ParseFiles("./templates/" + tmpl + ".page.html")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    t.Execute(w, p)
}


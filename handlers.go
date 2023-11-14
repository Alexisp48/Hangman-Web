package main

import (
	"bufio"
	"html/template"
	"log"
	"net/http"
	"os"
    "math/rand"
)

type Player struct {
    Name string
    G *Game
}

type Game struct {
    World string
    WorldFind string
    LetterTested string
    Win string
    TryNumber int
}

func (P *Player) Home(w http.ResponseWriter, r *http.Request) {
    P.renderTemplates(w, "home")
}

func replaceAtIndex(in *string, r rune, i int) {
    out := []rune(*in)
    out[i] = rune(r)
    *in = string(out)
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

func (P *Player) reset() {
    P.G.LetterTested = ""
    P.G.World = ""
    P.G.TryNumber = 10
    P.G.Win = "inGame"
    P.G.WorldFind = ""
}

func (P *Player) Hangman(w http.ResponseWriter, r *http.Request) {
    if P.G.Win != "inGame" {
        P.reset()
    }
    if len(P.G.World) == 0 { // initialse le mot
        P.G.World = openFile("worldDataBase.txt")[rand.Intn(len(openFile("worldDataBase.txt")))]
        for i := 0; i < len(P.G.World); i++ {
            P.G.WorldFind += "_"
        }
        P.renderTemplates(w, "hangman")
    } else {
        letter := r.FormValue("letter")
        if letter != "" {
            if len(letter) == 1{ // test une lettre
                var turn = true
                for i := 0; i < len(P.G.LetterTested)/3; i++ {
                    if []rune(letter)[0] == rune(P.G.LetterTested[i*3]) {
                        turn = false
                    }
                }
                if turn {
                    var letterInWorld = false
                    for i := 0; i < len(P.G.World); i++ {
                        if rune(P.G.World[i]) == []rune(letter)[0] {
                            letterInWorld = true
                            replaceAtIndex(&P.G.WorldFind, []rune(letter)[0], i)
                        }
                    }
                    if !letterInWorld {
                        P.G.TryNumber--
                    }
                    P.G.LetterTested += letter + "  "
                }
            }else if letter == P.G.World { // test un mot et c'est le bon
                P.G.Win = "win" // win
            } else { // test pas bon lettre ou mot
                P.G.TryNumber--
            }
            if P.G.WorldFind == P.G.World {
                P.G.Win = "win"
            }
            if P.G.TryNumber <= 0 {
                P.G.Win = "lose"
            }
            P.renderTemplates(w, "hangman")
        }
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


package main

import (
	"bufio"
	"encoding/json"
	"html/template"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
)

type Player struct {
	Name  string `json:"name"`
	Pwd   string `json:"pwd"`
	Pts   int    `json:"pts"`
	Id    int    `json:"id"`
	G     *Game
	Login bool
}

var Users []Player

type Game struct {
	Word         string
	WordFind     string
	LetterTested string
	LetterColor  []string
	Win          string
	TryNumber    int
}

func (P *Player) Load(filePath string) {

	data, _ := ioutil.ReadFile(filePath)
	err := json.Unmarshal(data, &Users)

	if err != nil {
		log.Fatal(err)
	}

}

func (P *Player) Save(filePath string) {

	Users[P.Id] = *P // Update values

	data, err := json.Marshal(Users)

	if err != nil {
		log.Fatal(err)
	}

	ioutil.WriteFile(filePath, data, 0777)
}

func (P *Player) Home(w http.ResponseWriter, r *http.Request) {

	name := r.FormValue("name")
	pwd := r.FormValue("password")

	if !P.Login && name != "" && pwd != "" {

		UserExist := false
		for i := 0; i < len(Users); i++ {
			if Users[i].Name == name && Users[i].Pwd == pwd {
				P.Name = Users[i].Name
				P.Pwd = Users[i].Pwd
				P.Pts = Users[i].Pts
				P.G = &Game{TryNumber: 10, Win: "inGame"}
				P.Id = Users[i].Id
				UserExist = true
			}
		}

		if !UserExist {
			P.Name = name
			P.Pwd = pwd
			P.Pts = 0
			P.G = &Game{TryNumber: 10, Win: "inGame"}
			P.Id = len(Users)
			Users = append(Users, *P)
		}

		P.G.LetterColor = []string{"none", "none", "none", "none", "none", "none", "none", "none", "none", "none", "none", "none", "none", "none", "none", "none", "none", "none", "none", "none", "none", "none", "none", "none", "none", "none"}

		P.Login = true
		//fmt.Println(P)

	}
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
	P.G.Word = ""
	P.G.TryNumber = 10
	P.G.Win = "inGame"
	P.G.WordFind = ""
	P.G.LetterColor = []string{"none", "none", "none", "none", "none", "none", "none", "none", "none", "none", "none", "none", "none", "none", "none", "none", "none", "none", "none", "none", "none", "none", "none", "none", "none", "none"}
}

func (P *Player) Hangman(w http.ResponseWriter, r *http.Request) {
	if !P.Login {
		
		return
	}
	letter := r.FormValue("letter")
	letterAlphabet := r.FormValue("letterAlphabet")
	if P.G.Win != "inGame" {
		P.reset()
	}
	if len(P.G.Word) == 0 { // initialse le mot
		P.G.Word = openFile("wordDataBase.txt")[rand.Intn(len(openFile("wordDataBase.txt")))]
		for i := 0; i < len(P.G.Word); i++ {
			P.G.WordFind += "_"
		}
		P.renderTemplates(w, "hangman")
	} else {
		if letter != "" || letterAlphabet != "" {
			if letterAlphabet != "" {
				letter = letterAlphabet
			}
			letter = strings.ToUpper(letter)
			if len(letter) == 1 || letterAlphabet != "" { // test une lettre
				var turn = true
				var index = -1
				if 65 <= []rune(letter)[0] && []rune(letter)[0] <= 90 {
					index = int([]rune(letter)[0] - 65)
				} else {
					turn = false
				}
				for i := 0; i < len(P.G.LetterTested)/3; i++ {
					if []rune(letter)[0] == rune(P.G.LetterTested[i*3]) {
						turn = false
					}
				}
				if turn {
					var letterInWord = false
					for i := 0; i < len(P.G.Word); i++ {
						if rune(P.G.Word[i]) == []rune(letter)[0] {
							letterInWord = true
							replaceAtIndex(&P.G.WordFind, []rune(letter)[0], i)
							P.G.LetterColor[index] = "in"
						}
					}
					if !letterInWord {
						P.G.TryNumber--
						P.G.LetterColor[index] = "out"
					}
					P.G.LetterTested += letter + "  "
				}
			} else if letter == P.G.Word { // test un mot et c'est le bon
				P.G.Win = "win" // win
				P.Pts += 11
				P.Save("Users.json")
			} else { // test pas bon lettre ou mot
				P.G.TryNumber--
			}
			if P.G.WordFind == P.G.Word {
				P.G.Win = "win"
				P.Pts += 11
				P.Save("Users.json")
			}
			if P.G.TryNumber <= 0 {
				P.G.Win = "lose"
				P.Pts -= 5
				if P.Pts < 0 {
					P.Pts = 0
				}
				P.Save("Users.json")
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

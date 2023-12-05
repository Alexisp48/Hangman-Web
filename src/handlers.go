package HangmanWeb

import (
	"fmt"
	"html/template"
	"math/rand"
	"net/http"
	"strings"
	"time"
)


func (E *Engine) Home(w http.ResponseWriter, r *http.Request) {

	E.P.CurrentPage = "home"

	name := r.FormValue("name")
	pwd := r.FormValue("password")
	upgrade := r.FormValue("upgrade")
	food := r.FormValue("food")

	if !E.P.Login && name != "" && pwd != "" {

		E.P.Timer = time.Now()

		UserExist := false
		for i := 0; i < len(E.Users); i++ {
			if E.Users[i].Name == name && E.Users[i].Pwd == pwd {
				E.P.Name = E.Users[i].Name
				E.P.Pwd = E.Users[i].Pwd
				E.P.Pts = E.Users[i].Pts
				E.P.G = &Game{TryNumber: 10, Win: "inGame"}
				E.P.R = E.Users[i].R
				E.P.Id = E.Users[i].Id
				E.P.R.AgePrice = []int{0, 20, 47, 68, 95, 112, 130, 143, 155}
				UserExist = true
			}
		}

		if !UserExist {
			E.P.Name = name
			E.P.Pwd = pwd
			E.P.Pts = 0
			E.P.G = &Game{TryNumber: 10, Win: "inGame"}
			E.P.R = &Resource{Food: 0, Age: 1}
			E.P.R.AgePrice = []int{0, 20, 47, 68, 95, 112, 130, 143, 155}
			E.P.Id = len(E.Users)
			E.Users = append(E.Users, E.P)
		}

		E.P.G.LetterColor = []string{"none", "none", "none", "none", "none", "none", "none", "none", "none", "none", "none", "none", "none", "none", "none", "none", "none", "none", "none", "none", "none", "none", "none", "none", "none", "none"}

		E.P.Login = true
		//fmt.Println(P)

		go E.useFood(w)
	} else if E.P.Login && upgrade != "" && E.P.Pts >= E.P.R.AgePrice[E.P.R.Age] {
		E.P.Pts -= E.P.R.AgePrice[E.P.R.Age]
		E.P.R.Age += 1
	} else if E.P.Login && food != "" && E.P.Pts >= 25 {
		E.P.Pts -= 25
		E.P.R.Food += 30
	}

	E.P.renderTemplates(w, "home")

}

func replaceAtIndex(in *string, r rune, i int) {
	out := []rune(*in)
	out[i] = rune(r)
	*in = string(out)
}

func (P *Player) reset() {
	P.G.LetterTested = ""
	P.G.Word = ""
	P.G.TryNumber = 10
	P.G.Win = "inGame"
	P.G.WordFind = ""
	P.G.LetterColor = []string{"none", "none", "none", "none", "none", "none", "none", "none", "none", "none", "none", "none", "none", "none", "none", "none", "none", "none", "none", "none", "none", "none", "none", "none", "none", "none"}
}

func (E *Engine) useFood(w http.ResponseWriter) {
	if time.Since(E.P.Timer).Seconds() >= 30 {
		E.P.Timer = time.Now()
		E.P.R.Food -= int(float32(E.P.R.Age/2) * 6)
		E.Save("data/Users.json")
	}
	go E.useFood(w)
}

func (E *Engine) Hangman(w http.ResponseWriter, r *http.Request) {
	E.P.CurrentPage = "hangman"
	if !E.P.Login {
		fmt.Fprintf(w, "Your are not login")
		fmt.Printf("Your are not login")
		return
	}
	letter := r.FormValue("letter")
	letterAlphabet := r.FormValue("letterAlphabet")
	if E.P.G.Win != "inGame" {
		E.P.reset()
	}
	if len(E.P.G.Word) == 0 { // initialse le mot
		E.P.G.Word = openFile("wordDataBase.txt")[rand.Intn(len(openFile("wordDataBase.txt")))]
		for i := 0; i < len(E.P.G.Word); i++ {
			E.P.G.WordFind += "_"
		}
		E.P.renderTemplates(w, "hangman")
	} else if letter == "" && letterAlphabet == "" {
		E.P.renderTemplates(w, "hangman")
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
				for i := 0; i < len(E.P.G.LetterTested)/3; i++ {
					if []rune(letter)[0] == rune(E.P.G.LetterTested[i*3]) {
						turn = false
					}
				}
				if turn {
					var letterInWord = false
					for i := 0; i < len(E.P.G.Word); i++ {
						if rune(E.P.G.Word[i]) == []rune(letter)[0] {
							letterInWord = true
							replaceAtIndex(&E.P.G.WordFind, []rune(letter)[0], i)
							E.P.G.LetterColor[index] = "in"
						}
					}
					if !letterInWord {
						E.P.G.TryNumber--
						E.P.G.LetterColor[index] = "out"
					}
					E.P.G.LetterTested += letter + "  "
				}
			} else if letter == E.P.G.Word { // test un mot et c'est le bon
				E.P.G.Win = "win" // win
				E.P.Pts += int(float32(E.P.R.Age/2) * 12)
				E.P.R.Food += 5
				E.Save("data/Users.json")
			} else { // test pas bon lettre ou mot
				E.P.G.TryNumber--
			}
			if E.P.G.WordFind == E.P.G.Word {
				E.P.G.Win = "win"
				E.P.Pts += int(float32(E.P.R.Age/2) * 12)
				E.P.R.Food += 5
				E.Save("data/Users.json")
			}
			if E.P.G.TryNumber <= 0 {
				E.P.G.Win = "lose"
				if E.P.Pts < 0 {
					E.P.Pts = 0
				}
				E.Save("data/Users.json")
			}
			E.P.renderTemplates(w, "hangman")
		}
	}
}

func (P *Player) renderTemplates(w http.ResponseWriter, tmpl string) {
	t, err := template.ParseFiles("./serv/templates/" + tmpl + ".page.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	t.Execute(w, P)
}

package HangmanWeb

import (
	"fmt"
	"html/template"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

func (E *Engine) useFood(w http.ResponseWriter) { // utilise 12 de nourriture tout les 20 secondes
	if time.Since(E.P.Timer).Seconds() >= 20 && E.P.CurrentPage == "home" {
		E.P.Timer = time.Now()
		E.P.R.Food -= 12
		if E.P.R.Food <= 0 {
			if E.P.R.Age != 1 {
				E.P.R.Age -= 1
				E.P.R.Food = 30
			} else {
				E.P.R.Food = 0
			}
		}
	}
	go E.useFood(w)
}

func (E *Engine) Home(w http.ResponseWriter, r *http.Request) { // page principale

	E.P.CurrentPage = "home"

	name := r.FormValue("name")
	pwd := r.FormValue("password")
	upgrade := r.FormValue("upgrade")
	food := r.FormValue("food")

	if !E.P.Login && name != "" && pwd != "" {

		E.P.Timer = time.Now()

		UserExist := false
		for i := 0; i < len(E.Users); i++ { // login
			if E.Users[i].Name == name && E.Users[i].Pwd == pwd {
				E.P.Name = E.Users[i].Name
				E.P.Pwd = E.Users[i].Pwd
				E.P.Gold = E.Users[i].Gold
				E.P.G = &Game{TryNumber: 10, Win: "inGame"}
				E.P.R = E.Users[i].R
				E.P.Position = E.Users[i].Position
				E.P.R.AgePrice = []int{0, 20, 47, 68, 95, 112, 130, 143, 155}
				UserExist = true
			}
		}

		if !UserExist { // créer un nouvel utilisateur
			E.P.Name = name
			E.P.Pwd = pwd
			E.P.Gold = 0
			E.P.G = &Game{TryNumber: 10, Win: "inGame"}
			E.P.R = &Resource{Food: 0, Age: 1}
			E.P.R.AgePrice = []int{0, 20, 47, 68, 95, 112, 130, 143, 155}
			E.P.Position = len(E.Users)
			E.Users = append(E.Users, E.P)
		}

		E.P.G.LetterColor = []string{"none", "none", "none", "none", "none", "none", "none", "none", "none", "none", "none", "none", "none", "none", "none", "none", "none", "none", "none", "none", "none", "none", "none", "none", "none", "none"}

		E.P.Login = true

		//fmt.Println(P)

		go E.useFood(w)

	} else if E.P.Login && upgrade != "" && E.P.Gold >= E.P.R.AgePrice[E.P.R.Age] {
		E.P.Gold -= E.P.R.AgePrice[E.P.R.Age]
		E.P.R.Age += 1
	} else if E.P.Login && food != "" && E.P.Gold >= 25 {
		E.P.Gold -= 25
		E.P.R.Food += 30
	}

	E.Save("data/Users.json") // sauvegarde

	E.P.renderTemplates(w, "home") // actualise la page

}

func replaceAtIndex(in *string, r rune, i int) { // fonction utile
	out := []rune(*in)
	out[i] = rune(r)
	*in = string(out)
}

func (P *Player) reset() { // remet toutes les valeurs de base à la fin d'une partie
	P.G.LetterTested = ""
	P.G.Word = ""
	P.G.TryNumber = 10
	P.G.Win = "inGame"
	P.G.WordFind = ""
	P.G.LetterColor = []string{"none", "none", "none", "none", "none", "none", "none", "none", "none", "none", "none", "none", "none", "none", "none", "none", "none", "none", "none", "none", "none", "none", "none", "none", "none", "none"}
}

func (E *Engine) Hangman(w http.ResponseWriter, r *http.Request) { // page hangman
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
	if len(E.P.G.Word) == 0 { // initialse le mot à deviner
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
				for i := 0; i < len(E.P.G.LetterTested)/3; i++ { // ajoute la lettre tester dans la liste de lettre testé
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
				E.P.Gold += 13
				E.P.R.Food += 5
				E.Save("data/Users.json") // save
			} else { // test pas bonne lettre ou mot
				E.P.G.TryNumber--
			}
			if E.P.G.WordFind == E.P.G.Word {
				E.P.G.Win = "win"
				E.P.Gold += 13
				E.P.R.Food += 5
				E.Save("data/Users.json")
			}
			if E.P.G.TryNumber <= 0 { // trop de tentatives c'est perdu
				E.P.G.Win = "lose"
				E.P.Gold -= 4
				if E.P.Gold < 0 {
					E.P.Gold = 0
				}
				E.Save("data/Users.json")
			}
			E.P.renderTemplates(w, "hangman") // actuallise la page
		}
	}
}

func (P *Player) renderTemplates(w http.ResponseWriter, tmpl string) { // affiche la page passer en paramètre
	t, err := template.ParseFiles("./serv/templates/" + tmpl + ".page.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	t.Execute(w, P)
}

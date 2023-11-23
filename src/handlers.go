package HangmanWeb

import (
	"html/template"
	"math/rand"
	"net/http"
	"strings"
)


func (E *Engine) Home(w http.ResponseWriter, r *http.Request) {

	name := r.FormValue("name")
	pwd := r.FormValue("password")

	if !E.P.Login && name != "" && pwd != "" {

		UserExist := false
		for i := 0; i < len(E.Users); i++ {
			if E.Users[i].Name == name && E.Users[i].Pwd == pwd {
				E.P.Name = E.Users[i].Name
				E.P.Pwd = E.Users[i].Pwd
				E.P.Pts = E.Users[i].Pts
				E.P.G = &Game{TryNumber: 10, Win: "inGame"}
				E.P.Id = E.Users[i].Id
				UserExist = true
			}
		}

		if !UserExist {
			E.P.Name = name
			E.P.Pwd = pwd
			E.P.Pts = 0
			E.P.G = &Game{TryNumber: 10, Win: "inGame"}
			E.P.Id = len(E.Users)
			E.Users = append(E.Users, E.P)
		}

		E.P.G.LetterColor = []string{"none", "none", "none", "none", "none", "none", "none", "none", "none", "none", "none", "none", "none", "none", "none", "none", "none", "none", "none", "none", "none", "none", "none", "none", "none", "none"}

		E.P.Login = true
		//fmt.Println(P)

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

func (E *Engine) Hangman(w http.ResponseWriter, r *http.Request) {
	if !E.P.Login {
		
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
				E.P.Pts += 11
				E.Save("data/Users.json")
			} else { // test pas bon lettre ou mot
				E.P.G.TryNumber--
			}
			if E.P.G.WordFind == E.P.G.Word {
				E.P.G.Win = "win"
				E.P.Pts += 11
				E.Save("data/Users.json")
			}
			if E.P.G.TryNumber <= 0 {
				E.P.G.Win = "lose"
				E.P.Pts -= 5
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

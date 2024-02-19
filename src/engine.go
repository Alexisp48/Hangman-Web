package HangmanWeb

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

type Player struct { // données du joueur
	Name        string    `json:"name"`
	Pwd         string    `json:"pwd"`
	Gold        int       `json:"pts"`
	Position    int       `json:"id"`
	R           *Resource `json:"R"`
	CurrentPage string
	Timer       time.Time
	G           *Game
	Login       bool
}

type Resource struct {// données des resources du joueur
	Food     int
	Age      int // Avancement de la civilisation
	AgePrice []int
}

type Game struct { // données du jeu
	Word         string
	WordFind     string
	LetterTested string
	LetterColor  []string
	Win          string
	TryNumber    int
}

type Engine struct { // données du site en général
	P     Player
	Users []Player
	Port  string
}

func (E *Engine) Init() {  // initialise les données
	rand.Seed(time.Now().UnixNano())

	E.Load("data/Users.json")

	p := E.Users[0]
	p.Login = false

	E.P = p

	E.Port = ":8080"
}

func (E *Engine) Run() { // lancement du site

	E.Init()

	fs := http.FileServer(http.Dir("serv"))
	http.Handle("/serv/", http.StripPrefix("/serv/", fs))

	http.HandleFunc("/", E.Home)
	http.HandleFunc("/hangman", E.Hangman)

	fmt.Println("(http://localhost:8080) - Serveur started on port", E.Port)

	http.ListenAndServe(E.Port, nil) // créer le serveur local au bon port
}

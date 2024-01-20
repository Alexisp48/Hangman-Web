package HangmanWeb

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

type Player struct {
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

type Resource struct {
	Food     int
	Age      int // Avancement de la civilisation
	AgePrice []int
}

type Game struct {
	Word         string
	WordFind     string
	LetterTested string
	LetterColor  []string
	Win          string
	TryNumber    int
}

type Engine struct {
	P     Player
	Users []Player
	Port  string
}

func (E *Engine) Init() {
	rand.Seed(time.Now().UnixNano())

	E.Load("data/Users.json")

	p := E.Users[0]
	p.Login = false

	E.P = p

	E.Port = ":8080"
}

func (E *Engine) Run() {

	E.Init()

	fs := http.FileServer(http.Dir("serv"))
	http.Handle("/serv/", http.StripPrefix("/serv/", fs))

	http.HandleFunc("/", E.Home)
	http.HandleFunc("/hangman", E.Hangman)

	fmt.Println("(http://localhost:8080) - Serveur started on port", E.Port)

	http.ListenAndServe(E.Port, nil)
}

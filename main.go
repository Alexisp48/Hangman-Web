package main

import (
	"math/rand"
	"fmt"
	"net/http"
	"time"
)

const (
    port = ":8080"
)

func main() {

    rand.Seed(time.Now().UnixNano())

    p := Player{"PAS ENCORE WIN T NUL", &Game{}}

    fs := http.FileServer(http.Dir("templates"))
	http.Handle("/templates/", http.StripPrefix("/templates/", fs))

    http.HandleFunc("/", p.Home)
    http.HandleFunc("/hangman", p.Hangman)

    fmt.Println("(http://localhost:8080) - Serveur started on port", port)

    http.ListenAndServe(port, nil)
}

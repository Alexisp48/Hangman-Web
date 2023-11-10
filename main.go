package main

import (
	"fmt"
	"net/http"
)

const (
    port = ":8080"
)

func main() {
    p := Player{"WORLD"}
    http.HandleFunc("/", p.Home)
    http.HandleFunc("/hangman", p.Hangman)

    fmt.Println("(http://localhost:8080) - Serveur started on port", port)

    http.ListenAndServe(port, nil)
}

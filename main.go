package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"time"
    HangmanWeb "HangmanWeb/src"
)

const (
    port = ":8080"
)

func main() {

    rand.Seed(time.Now().UnixNano())

    var unLog []HangmanWeb.Player

    data, _ := ioutil.ReadFile("data/Users.json")
    err := json.Unmarshal(data, &unLog)

    if err != nil {
        log.Fatal(err)
    }

    p := unLog[0]
    p.Login = false

    var e HangmanWeb.Engine

    e.P = p
    e.Users = unLog

    fs := http.FileServer(http.Dir("serv"))
	http.Handle("/serv/", http.StripPrefix("/serv/", fs))

    http.HandleFunc("/", e.Home)
    http.HandleFunc("/hangman", e.Hangman)

    fmt.Println("(http://localhost:8080) - Serveur started on port", port)

    http.ListenAndServe(port, nil)
}

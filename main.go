package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"time"
)

const (
    port = ":8080"
)

func main() {

    rand.Seed(time.Now().UnixNano())

    var unLog []Player

    data, _ := ioutil.ReadFile("Users.json")
    err := json.Unmarshal(data, &unLog)

    if err != nil {
        log.Fatal(err)
    }

    p := unLog[0]
    Users = unLog
    p.Login = false

    fs := http.FileServer(http.Dir("templates"))
	http.Handle("/templates/", http.StripPrefix("/templates/", fs))

    http.HandleFunc("/", p.Home)
    http.HandleFunc("/hangman", p.Hangman)

    fmt.Println("(http://localhost:8080) - Serveur started on port", port)

    http.ListenAndServe(port, nil)
}

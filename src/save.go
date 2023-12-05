package HangmanWeb

import (
	"log"
	"encoding/json"
	"io/ioutil"
)

func (E *Engine) Load(filePath string) {

	data, _ := ioutil.ReadFile(filePath)
	err := json.Unmarshal(data, &E.Users)

	if err != nil {
		log.Fatal(err)
	}

}

func (E *Engine) Save(filePath string) {

	E.Users[E.P.Id] = E.P // Update values //

	data, err := json.Marshal(E.Users)

	if err != nil {
		log.Fatal(err)
	}

	ioutil.WriteFile(filePath, data, 0777)
}
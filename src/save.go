package HangmanWeb

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

func (E *Engine) Load(filePath string) {

	data, _ := ioutil.ReadFile(filePath)
	err := json.Unmarshal(data, &E.Users)

	if err != nil {
		log.Fatal(err)
	}

}

func (E *Engine) Save(filePath string) {

	if E.P.Login {
		E.Users[E.P.Id] = E.P // Update values //

		for i := 1; i < len(E.Users); i++ {
			for j := 1; j < len(E.Users); j++ {
				if E.Users[i].Pts > E.Users[j].Pts {
					E.Users[i].Name, E.Users[j].Name = E.Users[j].Name, E.Users[i].Name
					E.Users[i].Pwd, E.Users[j].Pwd = E.Users[j].Pwd, E.Users[i].Pwd
					E.Users[i].Pts, E.Users[j].Pts = E.Users[j].Pts, E.Users[i].Pts
					E.Users[i].R, E.Users[j].R = E.Users[j].R, E.Users[i].R
					E.Users[i].CurrentPage, E.Users[j].CurrentPage = E.Users[j].CurrentPage, E.Users[i].CurrentPage
					E.Users[i].Timer, E.Users[j].Timer = E.Users[j].Timer, E.Users[i].Timer
					E.Users[i].G, E.Users[j].G = E.Users[j].G, E.Users[i].G
					E.Users[i].Login, E.Users[j].Login = E.Users[j].Login, E.Users[i].Login
				}
			}
		}

		for i := 0; i < len(E.Users); i++ {
			if E.Users[i].Name == E.P.Name && E.Users[i].Pwd == E.P.Pwd {
				E.P.Id = E.Users[i].Id
			}
		}





		data, err := json.Marshal(E.Users)

		if err != nil {
			log.Fatal(err)
		}

		ioutil.WriteFile(filePath, data, 0777)
	}


	
}
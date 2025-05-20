package register

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type UserAdder interface {
	AddUser(username, email, password string) error
}

type Request struct {
	Login             string `json:"login"`
	Email             string `json:"email"`
	Unhashed_password string `json:"password"`
}

func New(adder UserAdder) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("New request: ", r.RequestURI)
		reqBody, err := io.ReadAll(r.Body)
		if err != nil {
			log.Print("Reading body has been failed", "\n")
			fmt.Fprintf(w, "Error")
			return
		}

		var dat Request
		err = json.Unmarshal(reqBody, &dat)
		if err != nil {
			log.Print("Marshaling has been failed", "\n")
			fmt.Fprintf(w, "Error")
			return
		}

		err = adder.AddUser(dat.Login, dat.Email, dat.Unhashed_password)
		if err != nil {
			log.Print("Error: can not add user", "\n")
			fmt.Fprintf(w, "Error: can not add user")
			return
		}
	}
}

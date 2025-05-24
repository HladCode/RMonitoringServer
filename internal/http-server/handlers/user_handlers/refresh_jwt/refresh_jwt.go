package refresh_jwt

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	api_jwt "github.com/HladCode/RMonitoringServer/internal/lib/api/jwt"
	"github.com/HladCode/RMonitoringServer/internal/lib/e"
)

type refresherTokenValidator interface {
	IfRefreshTokenValid(unhashed_token string, user_id int) (bool, error)
	GetUserID(username string) (int, error)
}

type Request struct {
	Username      string `json:"username"`
	UnhashedToken string `json:"token"`
}

func New(refresher refresherTokenValidator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("New request: ", r.RequestURI)
		reqBody, err := io.ReadAll(r.Body)
		if err != nil {
			log.Print("Reading body has been failed", "\n")
			fmt.Fprintf(w, e.Wrap("Error", err).Error())
			return
		}

		var dat Request
		err = json.Unmarshal(reqBody, &dat)
		if err != nil {
			log.Print("Marshaling has been failed", "\n")
			fmt.Fprintf(w, e.Wrap("Error", err).Error())
			return
		}

		id, err := refresher.GetUserID(dat.Username)
		if err != nil {
			log.Print("Error: can not get user id", "\n")
			fmt.Fprintf(w, "Error: can not get user id")
			return
		}
		check, err := refresher.IfRefreshTokenValid(dat.UnhashedToken, id)
		if err != nil {
			log.Print("Error: can not check for refresh token", "\n")
			fmt.Fprintf(w, e.Wrap("Error: can not check for refresh token", err).Error())
			return
		}

		if check {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			accessToken, err := api_jwt.GenerateJWT(id)
			if err != nil {
				log.Print("Error: can not make JWT", "\n")
				fmt.Fprintf(w, "Error: can not make JWT")
				return
			}

			json.NewEncoder(w).Encode(map[string]string{
				"access_token": accessToken,
			})
			fmt.Fprintf(w, "200")
		} else {
			fmt.Fprint(w, "Invalid Token. Please login")
		}
	}
}

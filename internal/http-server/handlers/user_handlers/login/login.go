package login

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	api_jwt "github.com/HladCode/RMonitoringServer/internal/lib/api/jwt"
	"github.com/HladCode/RMonitoringServer/internal/storage"
	"golang.org/x/crypto/bcrypt"
)

type UserGetter interface {
	GetUser(username string) (storage.User_data, error)
	SaveRefreshToken(hashed_token string, user_id int) error // TODO: даты создания и просрочки сделать в самой функции
}

type Request struct {
	Login             string `json:"Login"`
	Unhashed_password string `json:"password"`
}

func New(getter UserGetter) http.HandlerFunc {
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

		user, err := getter.GetUser(dat.Login)
		if err != nil {
			log.Print("Error: can not find user", "\n")
			fmt.Fprintf(w, "Error: can not find user")
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(dat.Unhashed_password), bcrypt.DefaultCost)
		if err != nil {
			log.Print("Error: can not hash", "\n")
			fmt.Fprintf(w, "Error: can not hash")
			return
		}

		if user.Hashed_password == string(hashedPassword) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			accessToken, err := api_jwt.GenerateJWT(user.ID)
			if err != nil {
				log.Print("Error: can not make JWT", "\n")
				fmt.Fprintf(w, "Error: can not make JWT")
				return
			}
			refreshToken, hashedRefreshToken, err := api_jwt.GenerateRefreshToken()
			if err != nil {
				log.Print("Error: can not make refresh toket", "\n")
				fmt.Fprintf(w, "Error: can not make refresh toket")
				return
			}

			err = getter.SaveRefreshToken(hashedRefreshToken, user.ID)
			if err != nil {
				log.Print("Error: can not add refresh toket", "\n")
				fmt.Fprintf(w, "Error: can add make refresh toket")
				return
			}

			json.NewEncoder(w).Encode(map[string]string{
				"access_token":  accessToken,
				"refresh_token": refreshToken,
			})
			fmt.Fprintf(w, "200")
		} else {
			fmt.Fprintf(w, "401")
		}
	}
}

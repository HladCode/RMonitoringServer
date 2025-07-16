package getdevices

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/HladCode/RMonitoringServer/internal/lib/e"
)

type DeviceGetter interface {
	GetDevicesFromUserInJson(username string) (string, error) // map[string]string, []string
}

type Request struct {
	Username string `json:"Username"`
}

func New(getter DeviceGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("New request: ", r.RequestURI)
		reqBody, err := io.ReadAll(r.Body)
		if err != nil {
			log.Print("Reading body has been failed", "\n")
			fmt.Fprintf(w, "Error")
			return
		}

		//log.Println(string(reqBody))

		var dat Request
		err = json.Unmarshal(reqBody, &dat)
		if err != nil {
			log.Print("Marshaling has been failed", "\n", reqBody, "\n")
			fmt.Fprintf(w, "Error")
			return
		}

		resp, err := getter.GetDevicesFromUserInJson(dat.Username)
		if err != nil {
			log.Println(err.Error())
			fmt.Fprintf(w, e.Wrap("Error: ", err).Error())
			return
		}

		fmt.Fprintf(w, resp)
	}
}

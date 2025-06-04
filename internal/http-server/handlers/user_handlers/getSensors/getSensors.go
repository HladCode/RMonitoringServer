package getSensors

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type SensorsGetter interface {
	GetSensorsFromDeviceIdInJson(device_id string) (string, error)
}

type Request struct {
	Device_id string `json:"Device_id"`
}

func New(getter SensorsGetter) http.HandlerFunc {
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
			log.Print("Marshaling has been failed", "\n", reqBody, "\n")
			fmt.Fprintf(w, "Error")
			return
		}

		resp, err := getter.GetSensorsFromDeviceIdInJson(dat.Device_id)
		if err != nil {
			log.Println(err.Error())
		}

		log.Println(string(reqBody), resp)

		fmt.Fprintf(w, resp)
	}
}

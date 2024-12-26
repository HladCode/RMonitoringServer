package getData

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	re "github.com/HladCode/RMonitoringServer/internal/lib/api/response"
)

type Request struct {
	Timestamp   string `json:"timestamp"`
	ObjectName  string `json:"object_name"`
	PhoneNumber string `json:"phone_number"`
	Temperature string `json:"tempreature"`
}

type DataSaver interface {
	SaveTemperature(timestamp, phoneNumber, refrigeratorPath string, temp float64) error
}

func New(saver DataSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//TODO: make normal log
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

		//log.Println(dat.Tempreature, "°C", ", ", dat.Path)
		t, _ := strconv.ParseFloat(dat.Temperature, 32)

		if err = saver.SaveTemperature(dat.Timestamp, dat.PhoneNumber, dat.ObjectName, t); err != nil {
			log.Println(err.Error())
		}

		fmt.Fprintf(w, re.ArduinoOk().Status)
	}
}

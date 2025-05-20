package getData_old

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	re "github.com/HladCode/RMonitoringServer/internal/lib/api/response"
)

type Request struct {
	ID              string  `json:"id"`
	Purpose         string  `json:"p"`
	SensorPinNumber string  `json:"n"` // TODO: in esp32 code change n in request json
	Timestamp       string  `json:"t"`
	SensorValue     float64 `json:"v"`
}

type DataSaver interface {
	SaveData(ID, Purpose, SensorPinNumber, Timestamp string, data float64) error
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

		//log.Println(reqBody)

		var dat Request
		err = json.Unmarshal(reqBody, &dat)
		if err != nil {
			log.Print("Marshaling has been failed", "\n", reqBody, "\n")
			fmt.Fprintf(w, "Error")
			return
		}

		//log.Println(dat.Tempreature, "°C", ", ", dat.Path)
		//t, _ := strconv.ParseFloat(dat.Temperature, 32)

		if err = saver.SaveData(dat.ID, dat.Purpose, dat.SensorPinNumber, dat.Timestamp, dat.SensorValue); err != nil {
			log.Println(err.Error())
		}

		fmt.Fprintf(w, re.ArduinoOk().Status)
	}
}

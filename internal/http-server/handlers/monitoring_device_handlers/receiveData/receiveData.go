package receivedata

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	re "github.com/HladCode/RMonitoringServer/internal/lib/api/response"
	"github.com/HladCode/RMonitoringServer/internal/storage"
)

type DataSaver interface {
	AddNewData(readings []storage.Data_unit) error
}

type Request struct {
	AllCurrentData []storage.Data_unit `json:"AllCurrentData"`
}

func New(saver DataSaver) http.HandlerFunc {
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

		if err = saver.AddNewData(dat.AllCurrentData); err != nil {
			log.Println(err.Error())
		}

		fmt.Fprintf(w, re.ArduinoOk().Status)
	}
}

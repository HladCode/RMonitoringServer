package getData

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"

	re "github.com/HladCode/RMonitoringServer/internal/lib/api/response"
)

type Request struct {
	Temperature string `json:"tempreature"`
	PhoneNumber string `json:"phone_number"`
	ObjectName  string `json:"object_name"`

	//for this arduino must have RTC module
	// time        string `json:"time"`
}

// type Data struct {
// 	Temperature string `json:"t"`
// 	Pressure    string `json:"p"`
// }

type DataSaver interface {
	SaveTemperature(timestamp, refrigeratorPath string, temp float64) error
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
		//t, _ := strconv.ParseFloat(dat.Temperature, 32)
		var mu sync.Mutex

		mu.Lock()
		//saver.SaveTemperature(time.Now().Format("2006-01-02 15:04:05"), dat.Path, t)
		mu.Unlock()
		fmt.Fprintf(w, re.ArduinoOk().Status)
	}
}

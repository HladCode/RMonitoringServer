package senddatafromday

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type DataGetter interface {
	GetDataFromDay(ID string, sensor_ID, day, month, year int) (string, error)
}

type Request struct {
	ID        string `json:"ID"`
	Sensor_ID int    `json:"sensor_ID"`
	Day       int    `json:"day"`
	Month     int    `json:"month"`
	Year      int    `json:"year"`
}

func New(getter DataGetter) http.HandlerFunc {
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

		resp, err := getter.GetDataFromDay(dat.ID, dat.Sensor_ID, dat.Day, dat.Month, dat.Year)
		if err != nil {
			log.Println(err.Error())
		}

		fmt.Fprintf(w, resp)
	}
}

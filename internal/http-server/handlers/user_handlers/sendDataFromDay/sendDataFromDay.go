package senddatafromday

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type DataGetter interface {
	GetDataFromDay(ID string, sensor_ID int, fromDateTime, ToDateTime string) (string, error)
}

type Request struct {
	ID           string `json:"ID"`
	Sensor_ID    int    `json:"sensor_ID"`
	FromDateTime string `json:"from"`
	ToDateTime   string `json:"to"`
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

		resp, err := getter.GetDataFromDay(dat.ID, dat.Sensor_ID, dat.FromDateTime, dat.ToDateTime)
		if err != nil {
			log.Println(err.Error())
		}

		fmt.Fprintf(w, resp)
	}
}

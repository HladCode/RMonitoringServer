package getTime

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func New() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("New request: ", r.RequestURI)
		currentTime := time.Now()
		fmt.Fprintf(w, currentTime.Format("2006 01 02 15 04 05"))
	}
}

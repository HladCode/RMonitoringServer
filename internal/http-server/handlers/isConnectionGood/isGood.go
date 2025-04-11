package isConnectionGood

import (
	"fmt"
	"log"
	"net/http"

	re "github.com/HladCode/RMonitoringServer/internal/lib/api/response"
)

func New() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("New request: ", r.RequestURI)
		fmt.Fprintf(w, re.ArduinoOk().Status)
	}
}

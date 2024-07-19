package isConnectionGood

import (
	"fmt"
	"net/http"

	re "github.com/HladCode/RMonitoringServer/internal/lib/api/response"
)

func New() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, re.ArduinoOk().Status)
	}
}

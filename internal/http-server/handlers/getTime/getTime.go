package getTime

import (
	"fmt"
	"net/http"
	"time"
)

func New() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		currentTime := time.Now()
		fmt.Fprintf(w, currentTime.Format("2006 01 02 15 04 05"))
	}
}

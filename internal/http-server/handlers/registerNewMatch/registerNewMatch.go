package registerNewMatch

import "net/http"

type Request struct {
	ID string `json:"match_id"`
}

func New() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

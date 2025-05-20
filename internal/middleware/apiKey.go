package middleware

import "net/http"

func APIKeyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// key := r.Header.Get("X-API-Key")
		// if !ValidateAPIKey(key) { // TODO: ValidateAPIKey
		// 	http.Error(w, "Forbidden", http.StatusForbidden)
		// 	return
		// }
		next.ServeHTTP(w, r)
	})
}

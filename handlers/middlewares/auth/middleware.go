package authmiddleware

import (
	"net/http"
	"os"
)

func Inıt(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		xApiKey := r.Header.Get("X-API-KEY")
		if xApiKey != os.Getenv("X-API-KEY") {
			http.Error(w, "Access is denied", http.StatusUnauthorized)
			return
		}

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		next.ServeHTTP(w, r)
	})
}

package auth

import (
	"net/http"
)

func SignupH(w http.ResponseWriter, r *http.Request) {
	
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Signup endpoint"))
}

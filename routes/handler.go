package routes

import (
	"fmt"
	"net/http"
)

// HomeHandler handles the requests to the root endpoint.
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Hello World Rate Limited API by Joneco!")
}

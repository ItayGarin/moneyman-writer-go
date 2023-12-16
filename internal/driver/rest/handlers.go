package rest

import "net/http"

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

func handleGcsTransactionsUploadedEvent(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

package utils

import (
	"log"
	"net/http"
)

func Error(w http.ResponseWriter, msg string, status int) {
	if len(msg) == 0 {
		msg = http.StatusText(status)
	}

	log.Printf("[ERROR] [%d] %s\n", status, msg)

	if status >= http.StatusInternalServerError {
		http.Error(w, http.StatusText(status), status)
		return
	}

	http.Error(w, msg, status)
	return
}

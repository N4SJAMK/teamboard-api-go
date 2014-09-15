package utils

import (
	"encoding/json"
	"net/http"
	"strconv"
)

func ReadJSON(req *http.Request, v interface{}) error {
	defer req.Body.Close()
	return json.NewDecoder(req.Body).Decode(v)
}

func WriteJSON(res http.ResponseWriter, s int, v interface{}) {
	content, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	res.Header().Set("Content-Type", "application/json")
	res.Header().Set("Content-Length", strconv.Itoa(len(content)))

	res.WriteHeader(s)
	res.Write(content)
	return
}

package utils

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func NewResponse(w http.ResponseWriter, responseData Response) {
	w.Header().Set("Content-Type", "application/json")
	data, _ := json.Marshal(responseData)
	w.WriteHeader(responseData.Status)
	w.Write(data)
}

func DecodeBody(r *http.Request, data interface{}) error {
	return json.NewDecoder(r.Body).Decode(data)
}

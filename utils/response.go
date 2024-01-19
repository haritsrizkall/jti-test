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

func NewResponse(w http.ResponseWriter, status int, message string, data interface{}) {
	responseData := Response{
		Status:  status,
		Message: message,
		Data:    data,
	}
	w.Header().Set("Content-Type", "application/json")
	respJson, _ := json.Marshal(responseData)
	w.WriteHeader(responseData.Status)
	w.Write(respJson)
}

func DecodeBody(r *http.Request, data interface{}) error {
	return json.NewDecoder(r.Body).Decode(data)
}

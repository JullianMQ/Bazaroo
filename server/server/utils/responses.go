package utils

import (
	"encoding/json"
	"net/http"
)

func ErrorRes(res http.ResponseWriter, status int, mess string) {
	res.WriteHeader(status)
	json.NewEncoder(res).Encode(ErrorResponse{
		Error: mess,
	})
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type OkResponse struct {
	Message string `json:"message"`
}

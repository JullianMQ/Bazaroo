package server

import (
	"encoding/json"
	"net/http"
)

func getRoot(res http.ResponseWriter, req *http.Request) {
	myMess := "Message"
	res.Header().Set("Content-Type", "application/json")
	json.NewEncoder(res).Encode(myMess)
}

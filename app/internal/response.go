package internal

import (
	"encoding/json"
	"net/http"
)

type APIResponse struct {
	Code int32       `json:"code"`
	Body interface{} `json:"body"`
}
type FibItem struct {
	Position uint32 `json:"position"`
	Item     uint32 `json:"item"`
}

func ResponseOk(ar []FibItem, w http.ResponseWriter) {
	resp := APIResponse{
		Code: http.StatusOK,
		Body: ar,
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func ResponseBad(msg string, w http.ResponseWriter) {
	resp := APIResponse{
		Code: http.StatusBadRequest,
		Body: msg,
	}
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(resp)
}

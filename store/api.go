package store

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type APIHandler struct {
	db *sql.DB
}

type POSTTokens struct {
	Tokens []string `json:"tokens"`
}

type GETToken struct {
	Token string `json:"token"`
}

func NewAPIHandler(db *sql.DB) *APIHandler {
	return &APIHandler{db}
}

func (h *APIHandler) TokenEndpoint(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		// get a token
		body, status := h.GetToken(r)
		w.WriteHeader(status)
		w.Write(body)
	case "POST":
		// insert tokens
		body, status := h.AddTokens(r)
		w.WriteHeader(status)
		w.Write(body)
	}
}

func (h *APIHandler) AddTokens(r *http.Request) ([]byte, int) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return []byte("Unable to read payload"), 400
	}
	payload := &POSTTokens{}
	err = json.Unmarshal(body, payload)
	if err != nil {
		return []byte("Invalid payload received"), 400
	}
	insertTokens(payload.Tokens, h.db)
	return []byte("OK"), 200
}

func (h *APIHandler) GetToken(r *http.Request) ([]byte, int) {
	token, err := getSingleToken(h.db)
	if err != nil {
		return []byte(fmt.Sprintf("Error getting token from db: %v", err)), 500
	}

	resp, err := json.Marshal(GETToken{token})
	if err != nil {
		return []byte(fmt.Sprintf("Error marshalling token resp %v", err)), 500
	}
	return resp, 200
}

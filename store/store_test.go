package store

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB
var apiHost = "0.0.0.0:1234"
var apiLocation = "http://" + apiHost + "/tokens"

func TestMain(m *testing.M) {
	// Set up a testing DB
	os.Remove("testing.db")
	log.Println("Creating test DB")
	file, err := os.Create("testing.db")
	if err != nil {
		log.Fatal(err.Error())
	}
	file.Close()
	log.Println("testing.db created")
	DB, err = sql.Open("sqlite3", "testing.db")
	if err != nil {
		log.Println("Error opening db", err)
	}
	log.Println("database:", DB)
	createStoreTables(DB)

	// set up api routes
	h := APIHandler{DB}
	http.HandleFunc("/tokens", h.TokenEndpoint)

	// start the http server
	go http.ListenAndServe("0.0.0.0:1234", nil)

	os.Exit(m.Run())
}

func TestSubmitTokens(t *testing.T) {
	data := &POSTTokens{
		Tokens: []string{
			"A",
			"B",
			"C",
			"D",
		},
	}
	buffer := new(bytes.Buffer)
	json.NewEncoder(buffer).Encode(data)

	resp, err := http.Post(apiLocation, "application/json", buffer)
	if err != nil {
		t.Error(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}

	if string(body) != "OK" {
		t.Error("unexpected body: ", string(body))
	}

	if resp.StatusCode != 200 {
		t.Error("unexpected status code: ", resp.StatusCode)
	}
}
func TestGetToken(t *testing.T) {
	_, err := DB.Exec("INSERT INTO tokens (token) VALUES ('A'),('B'),('C')")
	if err != nil {
		t.Error(err)
	}

	getToken := func() string {
		resp, err := http.Get(apiLocation)
		if err != nil {
			t.Error(err)
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Error(err)
		}
		tokenResp := &GETToken{}
		err = json.Unmarshal(body, tokenResp)
		if err != nil {
			log.Println(string(body))
			t.Error(err)
		}
		return tokenResp.Token
	}

	// fetch and delete all tokens in db
	for i := 0; i < 3; i++ {
		token := getToken()
		log.Println(token)
	}

	// attempt a fetch when all tokens are gone
	token := getToken()
	log.Println(token)
}

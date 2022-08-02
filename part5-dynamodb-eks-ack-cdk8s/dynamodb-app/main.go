package main

import (
	"dynamodb-url-shortener/db"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/gorilla/mux"
)

var client *dynamodb.Client
var table string
var region string

func init() {
	table = os.Getenv("TABLE_NAME")
	if table == "" {
		log.Fatal("environment variable TABLE_NAME missing")
	}

	region = os.Getenv("AWS_REGION")
	if region == "" {
		log.Fatal("environment variable AWS_REGION missing")
	}

	log.Println("Table", table, "Region", region)

	client = dynamodb.NewFromConfig(*&aws.Config{Region: region})
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/", createShortCode).Methods(http.MethodPost)
	router.HandleFunc("/{shortcode}", accessURLWithShortCode).Methods(http.MethodGet)

	log.Println("starting server...")
	http.ListenAndServe(":8080", router)
}

func createShortCode(rw http.ResponseWriter, req *http.Request) {
	b, err := ioutil.ReadAll(req.Body)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	defer req.Body.Close()

	url := string(b)
	log.Println("URL", url)

	shortCode, err := db.SaveURL(url)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(rw).Encode(CreateShortCodeResponse{ShortCode: shortCode})
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	rw.WriteHeader(http.StatusCreated)
	log.Println("short code for", url, shortCode)
}

type CreateShortCodeResponse struct {
	ShortCode string
}

func accessURLWithShortCode(rw http.ResponseWriter, req *http.Request) {
	shortCode := mux.Vars(req)["shortcode"]

	url, err := db.GetLongURL(shortCode)

	if err != nil {
		if errors.Is(err, db.ErrUrlNotFound) {
			http.Error(rw, err.Error(), http.StatusNotFound)
		} else {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	rw.Header().Add("Location", url)
	rw.WriteHeader(http.StatusFound)

	log.Println("found short code for", url, shortCode)
}

package main

import (
	mysqlclient "MakeUpCalendar/client/mysql"
	smsclient "MakeUpCalendar/client/sms"
	"context"
	"encoding/json"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {

	var sms smsclient.Client
	{
		source := os.Getenv("SMS_CLIENT_SOURCE")
		username := os.Getenv("SMS_CLIENT_USERNAME")
		password := os.Getenv("SMS_CLIENT_PASSWORD")
		url := os.Getenv("SMS_CLIENT_URL")

		sms = smsclient.NewClient(source, username, password, url)
	}

	var db mysqlclient.Client
	{
		port := os.Getenv("DATABASE_PORT")
		username := os.Getenv("DATABASE_USER_NAME")
		password := os.Getenv("DATABASE_USER_PASSWORD")
		hostname := os.Getenv("DATABASE_HOST")
		dbname := os.Getenv("DATABASE_NAME")

		db = mysqlclient.NewClient(username, password, hostname, dbname, port)
	}

	var s Service
	{
		s = NewService(sms, db)
	}

	r := mux.NewRouter()

	r.Use(Auth)

	r.HandleFunc("/api/categories", func(w http.ResponseWriter, r *http.Request) {
		s.CategoriesHandler(w, r)
	})

	r.HandleFunc("/api/artists", func(w http.ResponseWriter, r *http.Request) {
		s.MakeUpArtistsHandler(w, r)
	})

	r.HandleFunc("/api/artist/{short-key}", func(w http.ResponseWriter, r *http.Request) {
		s.MakeUpArtistHandler(w, r)
	})

	r.HandleFunc("/api/appointment", func(w http.ResponseWriter, r *http.Request) {
		s.AppointmentHandler(w, r)
	})

	http.Handle("/", r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal("Server couldn't start: ", err)
	}
}

type Service interface {
	CategoriesHandler(w http.ResponseWriter, _ *http.Request)
	MakeUpArtistsHandler(w http.ResponseWriter, r *http.Request)
	MakeUpArtistHandler(w http.ResponseWriter, r *http.Request)
	AppointmentHandler(w http.ResponseWriter, r *http.Request)
}

type service struct {
	sms smsclient.Client
	db  mysqlclient.Client
}

func NewService(sms smsclient.Client, db mysqlclient.Client) Service {
	cli := &service{
		sms: sms,
		db:  db,
	}

	return cli
}

func (s service) CategoriesHandler(w http.ResponseWriter, _ *http.Request) {
	ctx := context.Background()
	categories := s.db.GetCategories(ctx)
	writeJSONResponse(w, categories, http.StatusOK)
}

func (s service) MakeUpArtistsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	makeUpArtists, err := s.db.GetMakeUpArtists(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSONResponse(w, makeUpArtists, http.StatusOK)
}

func (s service) MakeUpArtistHandler(w http.ResponseWriter, r *http.Request) {
	shortKey := mux.Vars(r)["short-key"]
	if shortKey == "" {
		http.Error(w, "This method expects short-key parameters", http.StatusBadRequest)
		return
	}

	makeUpArtist, err := s.db.GetMakeUpArtist(context.Background(), shortKey)
	if err != nil {
		writeJSONResponse(w, nil, http.StatusNotFound)
		return
	}

	writeJSONResponse(w, makeUpArtist, http.StatusOK)
}

func (s service) AppointmentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			writeJSONResponse(w, "Error reading request body", http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()

		var ar AppointmentRequest
		err = json.Unmarshal(body, &ar)
		if err != nil {
			writeJSONResponse(w, "Error decoding JSON", http.StatusInternalServerError)
			return
		}

		_, err = json.Marshal(ar)
		if err != nil {
			writeJSONResponse(w, "Error encoding JSON", http.StatusInternalServerError)
			return
		}

		dueDate, err := time.Parse("2006-01-02 15:04:05", ar.DueDate)
		if err != nil {
			writeJSONResponse(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = s.db.CreateAppointment(context.Background(), ar.MakeUpArtistID, ar.FirstName, ar.LastName, ar.PhoneNumber, ar.Email, ar.District, ar.Message, ar.PreferTimePeriodInDay, dueDate)
		if err != nil {
			writeJSONResponse(w, err.Error(), http.StatusInternalServerError)
			return
		}

		writeJSONResponse(w, "success", http.StatusOK)
		return
	}

	writeJSONResponse(w, "this method allow only POST request", http.StatusNotFound)
}

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		xApiKey := r.Header.Get("X-API-KEY")
		if xApiKey != os.Getenv("X-API-KEY") {
			http.Error(w, "Access is denied", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func writeJSONResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	r := Response{
		Data: data,
	}

	if err := json.NewEncoder(w).Encode(r); err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
	}
}

type (
	Response struct {
		Data interface{} `json:"data"`
	}

	AppointmentRequest struct {
		MakeUpArtistID        int    `json:"makeUpArtistID"`
		FirstName             string `json:"firstName"`
		LastName              string `json:"lastName"`
		PhoneNumber           string `json:"phoneNumber"`
		Email                 string `json:"email"`
		District              string `json:"district"`
		Message               string `json:"message"`
		PreferTimePeriodInDay string `json:"preferTimePeriodInDay"`
		DueDate               string `json:"dueDate"`
	}
)

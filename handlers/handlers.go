package handlers

import (
	mysqlclient "MakeUpCalendar/client/mysql"
	smsclient "MakeUpCalendar/client/sms"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"math/rand"
	"net/http"
	"os"
	"text/template"
	"time"
)

type Service interface {
	CategoriesHandler(w http.ResponseWriter, _ *http.Request)
	MakeUpArtistsHandler(w http.ResponseWriter, r *http.Request)
	MakeUpArtistHandler(w http.ResponseWriter, r *http.Request)
	AppointmentHandler(w http.ResponseWriter, r *http.Request)
	AppointmentViewHandler(w http.ResponseWriter, r *http.Request)
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
		ctx := context.Background()

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

		transactionId := generateTransactionID(8)

		err = s.db.CreateAppointment(context.Background(), ar.MakeUpArtistID, ar.FirstName, ar.LastName, ar.PhoneNumber, ar.Email, ar.District, ar.Message, ar.PreferTimePeriodInDay, transactionId, dueDate)
		if err != nil {
			writeJSONResponse(w, err.Error(), http.StatusInternalServerError)
			return
		}

		text := fmt.Sprintf("Talep: %s%s", os.Getenv("APPOINTMENT_VIEW_URL"), transactionId)
		s.sms.Send(ctx, text)

		writeJSONResponse(w, "success", http.StatusOK)
		return
	}

	writeJSONResponse(w, "this method allow only POST request", http.StatusNotFound)
}

func (s service) AppointmentViewHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	transactionID := mux.Vars(r)["id"]
	if transactionID == "" {
		http.Error(w, "This method expects transactionID parameters", http.StatusBadRequest)
		return
	}

	a, err := s.db.GetAppointment(ctx, transactionID)
	if err != nil {
		writeJSONResponse(w, nil, http.StatusNotFound)
		return
	}

	tmpl := `
		<!DOCTYPE html>
		<html>
		<head>
			<title>Veri Görüntüleme</title>
			<meta name="viewport" content="width=device-width, initial-scale=1">
			<link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/4.5.2/css/bootstrap.min.css">
		</head>
		<body>
			<div class="container-fluid">
					<h1 class="text-center">Kişisel Bilgiler</h1>
					<div class="row justify-content-center">
						<div class="col-md-6">
							<p><strong>Ad:</strong> {{.FirstName}} {{.LastName}}</p>
							<p><strong>E-posta:</strong> {{.Email}}</p>
							<p><strong>Telefon:</strong> {{.PhoneNumber}}</p>
							<p><strong>Bölge:</strong> {{.District}}</p>
						</div>
					</div>
					<div class="row justify-content-center">
						<div class="col-md-6">
							<h2 class="mt-4">Ek Bilgiler</h2>
							<p><strong>Tercih Edilen Zaman Dilimi:</strong> {{.PreferTimePeriodInDay}}</p>
							<p><strong>Oluşturulma Tarihi:</strong> {{.CreatedAt}}</p>
							<p><strong>Talep Edilen Tarih:</strong> {{.DueDate}}</p>
							<p><strong>Mesaj:</strong> {{.Message}}</p>
						</div>
					</div>
			</div>
		</body>
		</html>
	`

	tmplObj := template.Must(template.New("data").Parse(tmpl))

	err = tmplObj.Execute(w, a)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func writeJSONResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	w.WriteHeader(statusCode)

	r := Response{
		Data: data,
	}

	if err := json.NewEncoder(w).Encode(r); err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
	}
}

func generateTransactionID(length int) string {
	characters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	shuffledCharacters := make([]byte, length)
	for i := range shuffledCharacters {
		shuffledCharacters[i] = characters[r.Intn(len(characters))]
	}

	return string(shuffledCharacters)
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

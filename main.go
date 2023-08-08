package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	r := mux.NewRouter()

	r.Use(Auth)

	r.HandleFunc("/api/categories", CategoriesHandler)
	r.HandleFunc("/api/artists", MakeUpArtistsHandler)
	r.HandleFunc("/api/artist/{short-key}", MakeUpArtistHandler)
	r.HandleFunc("/api/appointment", AppointmentHandler)

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

func CategoriesHandler(w http.ResponseWriter, _ *http.Request) {
	categories := NewClient().GetCategories(context.Background())
	writeJSONResponse(w, categories, http.StatusOK)
}

func MakeUpArtistsHandler(w http.ResponseWriter, r *http.Request) {
	makeUpArtists, err := NewClient().GetMakeUpArtists(context.Background())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSONResponse(w, makeUpArtists, http.StatusOK)
}

func MakeUpArtistHandler(w http.ResponseWriter, r *http.Request) {
	shortKey := mux.Vars(r)["short-key"]
	if shortKey == "" {
		http.Error(w, "This method expects short-key parameters", http.StatusBadRequest)
		return
	}

	makeUpArtist, err := NewClient().GetMakeUpArtist(context.Background(), shortKey)
	if err != nil {
		writeJSONResponse(w, nil, http.StatusNotFound)
		return
	}

	writeJSONResponse(w, makeUpArtist, http.StatusOK)
}

func AppointmentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			writeJSONResponse(w, "Error reading request body", http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()

		var ar Appointment
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

		err = NewClient().CreateAppointment(context.Background(), ar.MakeUpArtistID, ar.FirstName, ar.LastName, ar.PhoneNumber, ar.Email, ar.District, ar.Message, ar.PreferTimePeriodInDay, dueDate)
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

type Response struct {
	Data interface{} `json:"data"`
}

type Category struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Priority int    `json:"priority"`
}

type MakeUpArtist struct {
	ID            string `json:"id"`
	FirstName     string `json:"firstName"`
	LastName      string `json:"lastName"`
	Image         string `json:"image"`
	ShortKey      string `json:"shortKey"`
	ServiceParams string `json:"services"`
}

type Appointment struct {
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

type Client interface {
	GetCategories(ctx context.Context) []Category
	GetMakeUpArtists(ctx context.Context) ([]MakeUpArtist, error)
	GetMakeUpArtist(ctx context.Context, ShortKey string) (*MakeUpArtist, error)
	CreateAppointment(ctx context.Context, makeUpArtistId int, firstName, lastName, phoneNumber, email, district, message, preferTimePeriodInDay string, dueDate time.Time) error
}

type client struct {
	connection string
}

func NewClient() Client {
	return &client{
		connection: generateDBConnection(),
	}
}

func (c client) GetCategories(_ context.Context) []Category {
	return []Category{
		{ID: "tumu", Name: "Tümü", Priority: 1},
		{ID: "dugun", Name: "Düğün Majyajı", Priority: 2},
		{ID: "nisan", Name: "Nişan Majyajı", Priority: 3},
		{ID: "mezuniyet", Name: "Mezuniyet Majyajı", Priority: 4},
	}
}

func (c client) GetMakeUpArtists(_ context.Context) ([]MakeUpArtist, error) {
	db, err := sql.Open("mysql", c.connection)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := "SELECT id, first_name, last_name, profile_photo_url, short_key, service_params  FROM makeup_artists WHERE is_active = 1 AND is_deleted = 0"
	rows, err := db.Query(query)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var makeUpArtists []MakeUpArtist
	for rows.Next() {
		var makeUpArtist MakeUpArtist
		err = rows.Scan(&makeUpArtist.ID, &makeUpArtist.FirstName, &makeUpArtist.LastName, &makeUpArtist.Image, &makeUpArtist.ShortKey, &makeUpArtist.ServiceParams)
		if err != nil {
			return nil, err
		}
		makeUpArtists = append(makeUpArtists, makeUpArtist)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return makeUpArtists, nil
}

func (c client) GetMakeUpArtist(_ context.Context, ShortKey string) (*MakeUpArtist, error) {
	db, err := sql.Open("mysql", c.connection)
	if err != nil {
		return nil, err
	}

	defer db.Close()

	var makeUpArtist MakeUpArtist

	query := "SELECT id, first_name, last_name, profile_photo_url, short_key FROM makeup_artists WHERE short_key = ?"
	err = db.QueryRow(query, ShortKey).Scan(&makeUpArtist.ID, &makeUpArtist.FirstName, &makeUpArtist.LastName, &makeUpArtist.Image, &makeUpArtist.ShortKey, &makeUpArtist.ServiceParams)
	if err != nil {
		return nil, err
	}

	return &makeUpArtist, nil
}

func (c client) CreateAppointment(ctx context.Context, makeUpArtistId int, firstName, lastName, phoneNumber, email, district, message, preferTimePeriodInDay string, dueDate time.Time) error {
	db, err := sql.Open("mysql", c.connection)
	if err != nil {
		return err
	}
	defer db.Close()

	query := `INSERT INTO appointments (makeup_artist_id, first_name, last_name, phone_number, email, district, due_date, message, prefer_time_period_in_day) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`

	args := []interface{}{
		makeUpArtistId,
		firstName,
		lastName,
		phoneNumber,
		email,
		district,
		dueDate,
		message,
		preferTimePeriodInDay,
	}

	_, err = db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}

func generateDBConnection() string {

	username := os.Getenv("DATABASE_USER_NAME")
	password := os.Getenv("DATABASE_USER_PASSWORD")
	hostname := os.Getenv("DATABASE_HOST")
	port := os.Getenv("DATABASE_PORT")
	dbname := os.Getenv("DATABASE_NAME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", username, password, hostname, port, dbname)

	return dsn
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

package main

import (
	mysqlclient "MakeUpCalendar/client/mysql"
	smsclient "MakeUpCalendar/client/sms"
	"MakeUpCalendar/handlers"
	authmiddleware "MakeUpCalendar/handlers/middlewares/auth"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
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

	var s handlers.Service
	{
		s = handlers.NewService(sms, db)
	}

	r := mux.NewRouter()

	r.Handle("/t/{id}", http.HandlerFunc(s.AppointmentViewHandler))
	r.Handle("/api/categories", authmiddleware.Inıt(http.HandlerFunc(s.CategoriesHandler)))
	r.Handle("/api/artists", authmiddleware.Inıt(http.HandlerFunc(s.MakeUpArtistsHandler)))
	r.Handle("/api/artist/{short-key}", authmiddleware.Inıt(http.HandlerFunc(s.MakeUpArtistHandler)))
	r.Handle("/api/appointment", authmiddleware.Inıt(http.HandlerFunc(s.AppointmentHandler)))

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "application/index.html")
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

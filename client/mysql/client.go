package mysqlclient

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

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

type Client interface {
	GetCategories(ctx context.Context) []Category
	GetMakeUpArtists(ctx context.Context) ([]MakeUpArtist, error)
	GetMakeUpArtist(ctx context.Context, shortKey string) (*MakeUpArtist, error)
	CreateAppointment(ctx context.Context, makeUpArtistId int, firstName, lastName, phoneNumber, email, district, message, preferTimePeriodInDay string, dueDate time.Time) error
}

type client struct {
	connection string
}

func NewClient(userName, password, hostName, dbName, port string) Client {
	return &client{
		connection: fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", userName, password, hostName, port, dbName),
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

func (c client) GetMakeUpArtist(_ context.Context, shortKey string) (*MakeUpArtist, error) {
	db, err := sql.Open("mysql", c.connection)
	if err != nil {
		return nil, err
	}

	defer db.Close()

	var makeUpArtist MakeUpArtist

	query := "SELECT id, first_name, last_name, profile_photo_url, short_key, service_params FROM makeup_artists WHERE short_key = ?"
	err = db.QueryRow(query, shortKey).Scan(&makeUpArtist.ID, &makeUpArtist.FirstName, &makeUpArtist.LastName, &makeUpArtist.Image, &makeUpArtist.ShortKey, &makeUpArtist.ServiceParams)
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

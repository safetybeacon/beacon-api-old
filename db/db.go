package db

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	_ "github.com/lib/pq"
	"github.com/t4ke0/locations_api/pkg/api"
)

var (
	ErrUserEmailAlreadyExist error = errors.New("user email already exists")
)

// Db
type Db struct {
	conn *sql.DB
}

func NewDB(postgresDSN string) (Db, error) {
	conn, err := sql.Open("postgres", postgresDSN)
	if err != nil {
		return Db{}, err
	}

	return Db{
		conn: conn,
	}, nil
}

// Close closes postgresql connection.
func (d Db) Close() error {
	return d.conn.Close()
}

// CreateTables create database tables.
func (d Db) CreateTables() (err error) {
	_, err = d.conn.Exec(`
CREATE TABLE IF NOT EXISTS users (
	id serial primary key,
	email varchar not null,
	password varchar not null,
	firstname varchar not null,
	lastname varchar not null,
	city varchar not null,
	country varchar not null,
	timecreated timestamp
);

CREATE TABLE IF NOT EXISTS tokens (
	id SERIAL PRIMARY KEY,
	user_id INTEGER not null,
	token varchar not null,
	name varchar not null,
	timecreated timestamp
);

CREATE TABLE IF NOT EXISTS locations (
	id SERIAL PRIMARY KEY,
	user_id integer not null,
	token_id integer not null,
	private boolean default false,
	latitude float not null,
	longitude float not null,
	time timestamp
);
	`)

	return
}

func (d Db) NewUser(
	email, password,
	firstname, lastname,
	city, country string,
) error {

	u, err := d.GetUser(email)
	if err != nil {
		return err
	}

	if u.Id != 0 {
		return ErrUserEmailAlreadyExist
	}

	result, err := d.conn.Exec(`
INSERT INTO users(email, password, firstname, lastname, city, country, timecreated)
VALUES($1, $2, $3, $4, $5, $6, $7)`, email, password, firstname, lastname, city, country, time.Now())
	if err != nil {
		return err
	}

	if n, _ := result.RowsAffected(); n == 0 {
		return fmt.Errorf("failed to insert new user into db")
	}

	return nil
}

type User struct {
	Id          int
	Email       string
	Password    string
	Firstname   string
	Lastname    string
	City        string
	Country     string
	TimeCreated time.Time
}

// GetUser
func (d Db) GetUser(email string) (u User, err error) {
	err = d.conn.QueryRow(`
SELECT id, email, password, firstname, lastname, city, country, timecreated
FROM users
WHERE email = $1`, email).Scan(
		&u.Id,
		&u.Email,
		&u.Password,
		&u.Firstname,
		&u.Lastname,
		&u.City,
		&u.Country,
		&u.TimeCreated,
	)
	if err != nil && err != sql.ErrNoRows {
		return
	}
	err = nil
	return
}

// NewToken
func (d Db) NewToken(userID int, token string, name string) (int, error) {
	result, err := d.conn.Exec(`
INSERT INTO tokens(user_id, token, name, timecreated)
VALUES($1, $2, $3, $4)`, userID, token, name, time.Now())
	if err != nil {
		return 0, err
	}

	if n, _ := result.RowsAffected(); n == 0 {
		return 0, fmt.Errorf("couldn't insert new token into the database table")
	}

	var id int
	if err := d.conn.QueryRow(`SELECT last_value FROM tokens_id_seq`).Scan(&id); err != nil {
		return id, err
	}

	return id, nil
}

// DeleteToken deletes token from the database where userID = user_id
func (d Db) DeleteToken(userID int, tokenID int) (int, error) {
	result, err := d.conn.Exec(`
DELETE FROM tokens WHERE user_id = $1 AND id = $2`, userID, tokenID)
	if err != nil {
		return 0, err
	}

	if n, _ := result.RowsAffected(); n == 0 {
		return 0, fmt.Errorf("couldn't delete token from the database")
	}

	return tokenID, nil
}

// GetTokenID
func (d Db) GetTokenID(userID int, token string) (id int, err error) {
	err = d.conn.QueryRow(`
SELECT id
FROM tokens
WHERE user_id = $1
AND token = $2`, userID, token).Scan(&id)
	if err != nil && err != sql.ErrNoRows {
		return
	}
	err = nil
	return
}

// CheckUserToken
func (d Db) CheckUserToken(userID int, token string) (bool, error) {
	var count int
	err := d.conn.QueryRow(`
SELECT count(token) 
FROM tokens
WHERE user_id = $1 AND token = $2`, userID, token).Scan(&count)
	if err != nil && err != sql.ErrNoRows {
		return false, err
	}

	return (count == 1), nil
}

// IsTokenExist
func (d Db) IsTokenExist(token string) (bool, error) {
	var count int
	err := d.conn.QueryRow(`
SELECT count(token) from tokens where token = $1`, token).Scan(&count)
	if err != nil && err != sql.ErrNoRows {
		return false, err
	}

	return (count == 1), nil
}

// NewLocation adds new location into the database.
func (d Db) NewLocation(userID, tokenID int, timestamp time.Time, latitude, longitude float64, private bool) (int, error) {
	result, err := d.conn.Exec(`
INSERT INTO locations(user_id, token_id, private, latitude, longitude, time)
VALUES($1, $2, $3, $4, $5, $6)`, userID, tokenID, private, latitude, longitude, timestamp)
	if err != nil {
		return 0, err
	}

	if n, _ := result.RowsAffected(); n == 0 {
		return 0, fmt.Errorf("cannot insert new location into the database")
	}

	var id int
	if err := d.conn.QueryRow(`SELECT last_value FROM locations_id_seq`).Scan(&id); err != nil {
		return id, err
	}
	return id, nil
}

// GetLocations
func (d Db) GetLocations() ([]api.Location, error) {
	rows, err := d.conn.Query(`
SELECT l.private, u.firstname, u.lastname, l.latitude, l.longitude, l.time
FROM locations l, users u 
WHERE l.user_id = u.id
	`)

	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	defer rows.Close()

	var locations []api.Location

	for rows.Next() {
		var location api.Location
		if err := rows.Scan(
			&location.Private,
			&location.Firstname,
			&location.Lastname,
			&location.Position.Latitude,
			&location.Position.Longitude,
			&location.Timestamp,
		); err != nil {
			return nil, err
		}

		if location.Private {
			location.Firstname = ""
			location.Lastname = ""
		}

		locations = append(locations, location)
	}

	return locations, nil

}

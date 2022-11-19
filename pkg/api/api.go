package api

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

// RegisterRequest register JSON request.
type RegisterRequest struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	City      string `json:"city"`
	Country   string `json:"country"`
}

// HashPassword hash user password.
func (rr *RegisterRequest) HashPassword() error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(rr.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	rr.Password = string(hashed)
	return nil
}

// Validate validates the register JSON request.
func (rr RegisterRequest) Validate() bool {
	return !(rr.Email == "" || rr.Password == "" ||
		rr.Firstname == "" || rr.Lastname == "" ||
		rr.City == "" || rr.Country == "")
}

// RegisterResponse register JSON response.
type RegisterResponse struct {
	Id int `json:"id"`
}

// LoginRequest login JSON request.
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Device   string `json:"device"`
}

// Validate validates the login JSON request.
func (lr LoginRequest) Validate() bool {
	return !(lr.Email == "" || lr.Password == "" || lr.Device == "")
}

// CompareHashAndPassword compares hashed password and user login request
// password. true means password is valid otherwise it's invalid.
func (lr LoginRequest) CompareHashAndPassword(hashedPassword string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(lr.Password)) == nil
}

// LoginResponse login JSON response.
type LoginResponse struct {
	Id    int    `json:"id"`
	Token string `json:"token"`
}

// LogoutRequest logout JSON request.
type LogoutRequest struct {
	Tokens []int `json:"tokens"`
}

// LogoutResponse logout JSON response.
type LogoutResponse struct {
	Invalidated []int `json:"invalidated"`
}

// Position latitude and longitude of the user.
type Position struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// AddLocationRequest add new location JSON request.
type AddLocationRequest struct {
	Timestamp time.Time `json:"timestamp"`
	Position  Position  `json:"position"`
	Private   bool      `json:"private"`
}

// Validate validates add location JSON request.
func (al AddLocationRequest) Validate() bool {
	return !(al.Timestamp == time.Time{} || al.Position.Latitude == 0 || al.Position.Longitude == 0)
}

// AddLocationResponse add location JSON response.
type AddLocationResponse struct {
	Id int `json:"id"`
}

// Location structure represents a location.
type Location struct {
	AddLocationRequest
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

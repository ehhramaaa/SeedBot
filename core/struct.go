package core

import "net/http"

type Client struct {
	username   string
	authToken  string
	proxy      string
	httpClient *http.Client
}

type Account struct {
	QueryId        string
	UserId         int
	Username       string
	FirstName      string
	LastName       string
	AuthDate       string
	Hash           string
	AllowWriteToPm bool
	LanguageCode   string
	QueryData      string
}

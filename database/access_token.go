package database

import (
	"database/sql"
	"time"

	uuid "github.com/satori/go.uuid"
)

// Token ...
type Token interface {
	Find() (*interface{}, error)
	Generate() (int, error)
}

// AccessTokenResponse ...
type AccessTokenResponse struct {
	ClientID     int
	AtokenID     int
	AtokenSecret string
}

// GiveAccess ...
func GiveAccess(t *Token, userID int) (AccessTokenResponse, error) {
	var response AccessTokenResponse

	client := Client{Type: "confidential"}
	clientID, err := client.Generate()
	if err != nil {
		return response, err
	}

	response.ClientID = clientID

	atoken := AccessToken{ClientID: clientID, Scopes: "", UserID: userID}
	atokenID, err := atoken.Generate()
	if err != nil {
		return response, err
	}

	response.AtokenID = atokenID

	return response, nil
}

// AccessToken type
type AccessToken struct {
	ID           int       `json:"id"`
	ClientID     int       `json:"client_id"`
	Expires      time.Time `json:"expires"`
	RefreshToken string    `json:"refresh_token"`
	Scopes       string    `json:"scopes"`
	Secret       string    `json:"secret"`
	UserID       int       `json:"user_id"`
	RedirectURI  string    `json:"redirect_uri"`
}

// Find returns an AccessToken
func (*AccessToken) Find(id int) (*AccessToken, error) {
	res, err := QueryRow("SELECT * from access_tokens WHERE id = $1", id)
	if err != nil {
		return nil, err
	}

	var atoken AccessToken

	if err = res.Scan(&atoken); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return &atoken, nil
}

// Generate saves an AccessToken in database
func (a *AccessToken) Generate() (int, error) {
	defaultScopes := "read-write"

	if a.Scopes == "" {
		a.Scopes = defaultScopes
	}

	a.Secret = uuid.NewV4().String()

	row, err := QueryRow(`
		INSERT INTO access_tokens (client_id, scopes, secret, user_id)
    	VALUES ($1, $2, $3, $4)
		RETURNING id`, a.ClientID, a.Scopes, a.Secret, a.UserID)
	if err != nil {
		return 0, err
	}

	var id int

	if err = row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

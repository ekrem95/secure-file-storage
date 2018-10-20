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

// AccessToken type
type AccessToken struct {
	ID           int            `json:"id"`
	ClientID     int            `json:"client_id"`
	Expires      time.Time      `json:"expires"`
	RefreshToken sql.NullString `json:"refresh_token"`
	Scopes       string         `json:"scopes"`
	Secret       string         `json:"secret"`
	UserID       int            `json:"user_id"`
	RedirectURI  sql.NullString `json:"redirect_uri"`
}

// Find returns an AccessToken
func (*AccessToken) Find(id float64) (*AccessToken, error) {
	res, err := QueryRow("SELECT id, client_id, expires, refresh_token, scopes, secret, user_id, redirect_uri from access_tokens WHERE id = $1", id)
	if err != nil {
		return nil, err
	}

	var atoken AccessToken

	if err = res.Scan(&atoken.ID, &atoken.ClientID, &atoken.Expires, &atoken.RefreshToken, &atoken.Scopes, &atoken.Secret, &atoken.UserID, &atoken.RedirectURI); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return &atoken, nil
}

// Generate saves an AccessToken in database
func (a *AccessToken) Generate() (int, string, error) {
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
		return 0, "", err
	}

	var id int

	if err = row.Scan(&id); err != nil {
		return 0, "", err
	}

	return id, a.Secret, nil
}

package database

import "database/sql"

// Client type
type Client struct {
	ID                 int    `json:"id"`
	Type               string `json:"type"`
	RedirectUris       string `json:"redirect_uris"`
	DefaultRedirectURI string `json:"default_redirect_uri"`
	AllowedGrantTypes  string `json:"allowed_grant_types"`
}

// Find returns a Client
func (*Client) Find(id int) (*Client, error) {
	res, err := QueryRow("SELECT * from clients WHERE id = $1", id)
	if err != nil {
		return nil, err
	}

	var client Client

	if err = res.Scan(&client); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return &client, nil
}

// Generate saves a Client in database
func (c *Client) Generate() (int, error) {
	var t string

	switch c.Type {
	case "confidential":
		t = c.Type
	default:
		t = "public"
	}

	row, err := QueryRow(`INSERT INTO clients (type) VALUES ($1) RETURNING id`, t)
	if err != nil {
		return 0, err
	}

	var id int

	if err = row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

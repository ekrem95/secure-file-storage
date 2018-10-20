package database

import (
	"errors"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

var errMalformedToken error = errors.New("Malformed token")

// CustomClaims ...
type CustomClaims struct {
	AtokenID int `json:"atoken_id"`
	ClientID int `json:"client_id"`
	jwt.StandardClaims
}

// NewJWTWithClaims creating a token using a custom claims type.
// The StandardClaim is embedded in the custom type to allow
// for easy encoding, parsing and validation of standard claims.
func NewJWTWithClaims(t *Token, userID int) (string, error) {
	client := Client{Type: "confidential"}
	cid, err := client.Generate()
	if err != nil {
		return "", err
	}

	atoken := AccessToken{ClientID: cid, Scopes: "", UserID: userID}
	aid, secret, err := atoken.Generate()
	if err != nil {
		return "", err
	}

	mySigningKey := []byte(secret)

	// Generate the Claims
	claims := CustomClaims{
		AtokenID: aid,
		ClientID: cid,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(mySigningKey)
	if err != nil {
		return "", err
	}

	return ss, nil

}

// ParseJWT is parsing the error types using bitfield checks
func ParseJWT(tokenString string) (*jwt.Token, jwt.MapClaims, error) {
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, &claims, mapClaims)

	if ve, ok := err.(*jwt.ValidationError); ok {
		if ve.Errors&jwt.ValidationErrorMalformed != 0 {
			return nil, nil, errMalformedToken
		}
	} else {
		return nil, nil, err
	}

	return token, claims, nil
}

func mapClaims(token *jwt.Token) (interface{}, error) {
	return []byte(""), nil
}

// CheckAuth ...
func CheckAuth(bearer string) (int, error) {
	token, claims, err := ParseJWT(bearer)
	if err != nil {
		return 0, err
	}

	cid := claims["client_id"]
	aid := claims["atoken_id"]

	if cid == nil || aid == nil {
		return 0, errMalformedToken
	}

	var client Client
	var atoken AccessToken

	cli, err := client.Find(cid.(float64))
	if err != nil {
		return 0, err
	}

	atk, err := atoken.Find(aid.(float64))
	if err != nil {
		return 0, err
	}

	if cli.ID != atk.ClientID {
		return 0, errMalformedToken
	}

	if time.Now().After(atk.Expires) {
		return 0, errors.New("Access token is expired")
	}

	// bearer token must be equal to signed string
	ss, err := token.SignedString([]byte(atk.Secret))
	if err != nil || ss != bearer {
		return 0, errors.New("Invalid token")
	}

	return atk.UserID, nil
}

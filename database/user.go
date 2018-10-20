package database

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

// User type
type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
}

// CheckPasswordHash compares a bcrypt hashed password with its possible plaintext equivalent.
func (u *User) CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// HashPassword returns the bcrypt hash of the password at the given cost.
func (u *User) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
}

// Save user
func (u *User) Save() error {
	hash, err := u.HashPassword(u.Password)
	if err != nil {
		return err
	}

	row, err := QueryRow(`INSERT INTO users (name, email, password) values($1, $2, $3) RETURNING id`, u.Name, u.Email, hash)
	if err != nil {
		return err
	}

	if err = row.Scan(&u.ID); err != nil {
		return err
	}

	u.Password = hash

	return nil
}

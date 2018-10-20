package database

import (
	"time"
)

// File type
type File struct {
	ID         int       `json:"id"`
	Name       string    `json:"name"`
	Path       string    `json:"path,omitempty"`
	Ext        string    `json:"ext,omitempty"`
	Algorithms string    `json:"algorithms"`
	UserID     int       `json:"user_id,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
}

// Files type
type Files []File

// Find all files of a user
func (files *Files) Find(uid int) error {
	result, err := Query(`SELECT id, name, algorithms, created_at FROM files WHERE user_id = $1`, uid)
	if err != nil {
		return err
	}

	var file File

	for result.Next() {
		if err = result.Scan(&file.ID, &file.Name, &file.Algorithms, &file.CreatedAt); err != nil {
			return err
		}

		*files = append(*files, file)
	}

	return nil
}

// Save file info
func (f *File) Save() error {
	_, err := Exec(`INSERT INTO files (name, path, ext, algorithms, user_id) values($1, $2, $3, $4, $5)`, f.Name, f.Path, f.Ext, f.Algorithms, f.UserID)
	if err != nil {
		return err
	}

	return nil
}

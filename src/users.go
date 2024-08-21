package main

import (
	"database/sql"
	"net/http"
	"path/filepath"
	"strings"
	"time"
)

type User struct {
	Email       string    `json:"email"`
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Image       string    `json:"image"`
	Location    string    `json:"location"`
	Type        string    `json:"type"`
	Birthday    time.Time `json:"birthday"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	IsDelete    bool      `json:"is_delete"`
	IsBan       bool      `json:"is_ban"`
	IsOfficial  bool      `json:"is_official"`
}

func (s *Server) HandleUsers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		s.handleUsersGet(w, r)
		return
	}
	respondHTTPErr(w, r, http.StatusNotFound)
}

func (s *Server) handleUsersGet(w http.ResponseWriter, r *http.Request) {
	sub := strings.TrimPrefix(r.URL.Path, "/users")
	_, id := filepath.Split(sub)
	if id != "" {
		user, err := s.findUserByID(id)
		if err != nil {
			if err == sql.ErrNoRows {
				respondErr(w, r, http.StatusNotFound, "User not found")
			} else {
				respondErr(w, r, http.StatusInternalServerError, "Failed to fetch user")
			}
			return
		}
		respond(w, r, http.StatusOK, user)
	} else {
		users, err := s.getAllUsers()
		if err != nil {
			respondErr(w, r, http.StatusInternalServerError, "Failed to fetch users")
			return
		}
		respond(w, r, http.StatusOK, users)
	}
}

func (s *Server) getAllUsers() ([]User, error) {
	query := `
		SELECT email, id, name, image, location, type, birthday, description, created_at, updated_at, is_delete, is_ban, is_official
		FROM users
	`
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(
			&user.Email,
			&user.ID,
			&user.Name,
			&user.Image,
			&user.Location,
			&user.Type,
			&user.Birthday,
			&user.Description,
			&user.CreatedAt,
			&user.UpdatedAt,
			&user.IsDelete,
			&user.IsBan,
			&user.IsOfficial,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

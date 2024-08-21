package main

import (
	"net/http"
	"time"
)

type MePutRequestType struct {
	Email        string `json:"email"`
	Name         string `json:"name"`
	Image        string
	favoriteTeam string
	Location     string
	Type         string
	Birthday     time.Time
	Description  string
	UpdatedAt    time.Time
}

func (s *Server) HandleMe(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		s.handleMeGet(w, r)
		return
	case "PUT":
		s.handleMePut(w, r)
		return
	}
	respondHTTPErr(w, r, http.StatusNotFound)
}

func (s *Server) handleMeGet(w http.ResponseWriter, r *http.Request) {
	userID, err := s.GetCurrentUser(w, r)
	if err != nil {
		respondErr(w, r, http.StatusUnauthorized, "Unauthorized")
		return
	}
	user, err := s.findUserByID(userID)
	if err != nil {
		respondErr(w, r, http.StatusInternalServerError, "Failed to get user")
		return
	}

	respond(w, r, http.StatusOK, user)
}

func (s *Server) handleMePut(w http.ResponseWriter, r *http.Request) {
	var user MePutRequestType
	if err := decodeBody(r, &user); err != nil {
		respondErr(w, r, http.StatusBadRequest, "Invalid request body")
		return
	}
	userID, err := s.GetCurrentUser(w, r)
	if err != nil {
		respondErr(w, r, http.StatusUnauthorized, "Unauthorized")
		return
	}

	err = s.updateUser(userID, user)
	if err != nil {
		respondErr(w, r, http.StatusInternalServerError, "Failed to update user")
		return
	}

	updatedUser, err := s.findUserByID(userID)
	if err != nil {
		respondErr(w, r, http.StatusInternalServerError, "Failed to get updated user")
		return
	}

	respond(w, r, http.StatusOK, updatedUser)
}

func (s *Server) updateUser(userID string, update MePutRequestType) error {
	query := `
		UPDATE users 
		SET email = ?, name = ?, image = ?, favorite_team = ?, location = ?, type = ?, birthday = ?, description = ?, updated_at = ?
		WHERE id = ?
	`
	_, err := s.db.Exec(query,
		update.Email, update.Name, update.Image, update.favoriteTeam,
		update.Location, update.Type, update.Birthday, update.Description,
		time.Now(), userID,
	)
	return err
}

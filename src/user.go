package main

import "net/http"

func (s *Server) handleUser(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.handleGetUser(w, r)
	default:
		respondHTTPErr(w, r, http.StatusMethodNotAllowed)
	}
}

func (s *Server) handleGetUser(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		respondErr(w, r, http.StatusBadRequest, "Missing user ID")
		return
	}

	user, err := s.findUserByID(id)
	if err != nil {
		respondErr(w, r, http.StatusInternalServerError, "Failed to get user")
		return
	}
	respond(w, r, http.StatusOK, user)
}

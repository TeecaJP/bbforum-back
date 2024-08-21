package main

import (
	"net/http"
	"time"
)

func (s *Server) HandleSignOut(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "auth",
		Value:    "/",
		Secure:   true,
		HttpOnly: true,
		Path:     "/",
		SameSite: 4,
		Expires:  time.Now().Add(365 * 24 * time.Hour),
	})
	respond(w, r, http.StatusOK, "")
}

package main

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("/app/.env")
	if err != nil {
		log.Panicln(err)
	}
	s := &Server{}
	if err := s.InitDB(); err != nil {
		log.Panicln(err)
	}
	defer s.CloseDB()
	GenerateBoardsFromJSON(s.db, "./players_by_position.json")

	ticker := time.NewTicker(5 * time.Second)
	go func() {
		for range ticker.C {
			// read json and generate board
			GenerateBoardsFromJSON(s.db, "./players_by_position.json")
		}
	}()

	mux := http.NewServeMux()
	mux.HandleFunc("/signin", withCORS(s.HandleSignIn))
	mux.HandleFunc("/signout", withCORS(s.HandleSignOut))
	mux.HandleFunc("/signup", withCORS(s.HandleSignUp))
	mux.HandleFunc("/me", withCORS(s.HandleMe))
	mux.HandleFunc("/users", withCORS(s.HandleUsers))
	mux.HandleFunc("/users/", withCORS(s.HandleUsers))
	mux.HandleFunc(("/posts"), withCORS(s.HandlePosts))
	mux.HandleFunc(("/post/"), withCORS(NeedToken(s.HandlePosts)))
	log.Println("start server")
	http.ListenAndServe(":8080", mux)
}

type Server struct {
	db *sql.DB
}

func withCORS(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-Requested-With, Origin, X-Csrftoken, Accept, Cookie")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		switch r.Method {
		case "OPTIONS":
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			return
		}
		fn(w, r)
	}
}

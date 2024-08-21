package main

import (
	"database/sql"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

func generateUUID() string {
	return uuid.New().String()
}

func isValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}

func Now() time.Time {
	nowFormatted := time.Now().Format(time.RFC3339)
	now, _ := time.Parse(time.RFC3339, nowFormatted)
	return now
}

func UpdateField(v interface{}) bson.M {
	val := reflect.ValueOf(v)
	typeOfT := val.Type()
	var updateField bson.M = bson.M{}

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		if !field.IsZero() {
			updateField[strings.ToLower(typeOfT.Field(i).Name)] = field.Interface()
		}
	}
	return updateField
}

func GetQueries(r *http.Request, params ...string) map[string]string {
	query := map[string]string{}
	for _, v := range params {
		query[v] = r.URL.Query().Get(v)
	}
	return query
}

func queryRows(db *sql.DB, query string, args ...interface{}) (*sql.Rows, error) {
	return db.Query(query, args...)
}

func queryRow(db *sql.DB, query string, args ...interface{}) *sql.Row {
	return db.QueryRow(query, args...)
}

func execQuery(db *sql.DB, query string, args ...interface{}) (sql.Result, error) {
	return db.Exec(query, args...)
}

func withTransaction(db *sql.DB, fn func(*sql.Tx) error) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := fn(tx); err != nil {
		return err
	}

	return tx.Commit()
}
func (s *Server) findUserByEmail(email string) (User, error) {
	var user User
	query := "SELECT id, email, name, image, type, created_at, updated_at FROM users WHERE email = ?"
	err := s.db.QueryRow(query, email).Scan(&user.ID, &user.Email, &user.Name, &user.Image, &user.Type, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func (s *Server) findUserByID(id string) (User, error) {
	var user User
	query := "SELECT id, email, name, image, type, created_at, updated_at FROM users WHERE id = ?"
	err := s.db.QueryRow(query, id).Scan(&user.ID, &user.Email, &user.Name, &user.Image, &user.Type, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

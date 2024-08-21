package main

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func (s *Server) InitDB() error {
	var err error
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?parseTime=true",
		os.Getenv("MYSQL_USER"),
		os.Getenv("MYSQL_PASSWORD"),
		os.Getenv("MYSQL_HOST"),
		os.Getenv("MYSQL_DATABASE"))

	s.db, err = sql.Open("mysql", dataSourceName)
	if err != nil {
		return fmt.Errorf("error opening database: %v", err)
	}

	s.db.SetMaxOpenConns(25)
	s.db.SetMaxIdleConns(25)
	s.db.SetConnMaxLifetime(5 * time.Minute)

	err = s.db.Ping()
	if err != nil {
		return fmt.Errorf("error connecting to the database: %v", err)
	}

	fmt.Println("Successfully connected to the database")
	return nil
}

func (s *Server) CloseDB() {
	if s.db != nil {
		s.db.Close()
	}
}

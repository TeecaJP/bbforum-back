package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
)

type Board struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	Tags      []string  `json:"tags"`
}

type Manager struct {
	Number string `json:"number"`
	Name   string `json:"name"`
}

type Player struct {
	Number string `json:"number"`
	Name   string `json:"name"`
}

type Roster struct {
	Manager    []Manager `json:"監督"`
	Pitcher    []Player  `json:"投手"`
	Catcher    []Player  `json:"捕手"`
	Infielder  []Player  `json:"内野手"`
	Outfielder []Player  `json:"外野手"`
}

type Team struct {
	Roster   Roster `json:"roster"`
	Optioned Roster `json:"optioned"`
}

type Teams struct {
	Buffaloes Team `json:"buffaloes"`
	Marines   Team `json:"marines"`
	Hawks     Team `json:"hawks"`
	Eagles    Team `json:"eagles"`
	Lions     Team `json:"lions"`
	Fighters  Team `json:"fighters"`
	Giants    Team `json:"giants"`
	Tigers    Team `json:"tigers"`
	Dragons   Team `json:"dragons"`
	Swallows  Team `json:"swallows"`
	Carp      Team `json:"carp"`
	BayStars  Team `json:"baystars"`
}

func BoardExists(db *sql.DB, name string) (bool, error) {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM boards WHERE name = ?)"
	err := db.QueryRow(query, name).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("error checking board existence: %v", err)
	}
	return exists, nil
}

func CreateBoard(db *sql.DB, name string, tags []string) (*Board, error) {
	exists, err := BoardExists(db, name)
	if err != nil {
		return nil, err
	}
	if exists {
		log.Printf("board %s already exists\n", name)
		return nil, nil
	}

	board := &Board{
		ID:        uuid.New().String(),
		Name:      name,
		CreatedAt: time.Now(),
		Tags:      tags,
	}

	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	query := `INSERT INTO boards (id, name, created_at) VALUES (?, ?, ?)`
	_, err = tx.Exec(query, board.ID, board.Name, board.CreatedAt)
	if err != nil {
		return nil, err
	}

	for _, tag := range tags {
		if err := addBoardTag(tx, board.ID, tag); err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	log.Printf("created board %s\n", name)

	return board, nil
}

func addBoardTag(tx *sql.Tx, boardID string, tagName string) error {
	var tagID string
	err := tx.QueryRow("SELECT id FROM tags WHERE name = ?", tagName).Scan(&tagID)
	if err == sql.ErrNoRows {
		tagID = uuid.New().String()
		_, err = tx.Exec("INSERT INTO tags (id, name) VALUES (?, ?)", tagID, tagName)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	_, err = tx.Exec("INSERT INTO board_tags (board_id, tag_id) VALUES (?, ?)", boardID, tagID)
	return err
}

func GenerateBoardsFromJSON(db *sql.DB, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	byteValue, _ := ioutil.ReadAll(file)

	var teams Teams
	if err := json.Unmarshal(byteValue, &teams); err != nil {
		return fmt.Errorf("error unmarshaling JSON: %v", err)
	}

	teamMap := map[string]Team{
		"buffaloes": teams.Buffaloes,
		"marines":   teams.Marines,
		"hawks":     teams.Hawks,
		"eagles":    teams.Eagles,
		"lions":     teams.Lions,
		"fighters":  teams.Fighters,
		"giants":    teams.Giants,
		"tigers":    teams.Tigers,
		"dragons":   teams.Dragons,
		"swallows":  teams.Swallows,
		"carp":      teams.Carp,
		"baystars":  teams.BayStars,
	}

	for teamName, team := range teamMap {
		_, err := CreateBoard(db, teamName, []string{"team", teamName})
		if err != nil {
			return fmt.Errorf("error creating board for team %s: %v", teamName, err)
		}

		for _, manager := range team.Roster.Manager {
			_, err := CreateBoard(db, manager.Name, []string{"manager", teamName})
			if err != nil {
				return fmt.Errorf("error creating board for manager %s: %v", manager.Name, err)
			}
		}

		playerTypes := []struct {
			players []Player
			tag     string
		}{
			{team.Roster.Pitcher, "pitcher"},
			{team.Roster.Catcher, "catcher"},
			{team.Roster.Infielder, "infielder"},
			{team.Roster.Outfielder, "outfielder"},
		}

		for _, pt := range playerTypes {
			for _, player := range pt.players {
				_, err := CreateBoard(db, player.Name, []string{pt.tag, teamName})
				if err != nil {
					return fmt.Errorf("error creating board for player %s: %v", player.Name, err)
				}
			}
		}
	}

	return nil
}

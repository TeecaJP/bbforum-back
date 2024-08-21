package main

import (
	"database/sql"
	"net/http"
	"time"
)

type Post struct {
	ID        string    `json:"id"`
	BoardID   string    `json:"board_id"`
	UserID    string    `json:"user_id"`
	Content   string    `json:"content"`
	ReplyTo   *string   `json:"reply_to"`
	CreatedAt time.Time `json:"created_at"`
	Tags      []string  `json:"tags"`
}

func (s *Server) HandlePost(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/posts/"):]
	if !isValidUUID(id) {
		respondErr(w, r, http.StatusBadRequest, "Invalid post ID")
		return
	}

	switch r.Method {
	case http.MethodGet:
		s.handlePostGet(w, r, id)
	case http.MethodPost:
		s.handlePostCreate(w, r)
	case http.MethodPut:
		s.handlePostUpdate(w, r, id)
	case http.MethodDelete:
		s.handlePostDelete(w, r, id)
	default:
		respondHTTPErr(w, r, http.StatusMethodNotAllowed)
	}
}

func (s *Server) handlePostCreate(w http.ResponseWriter, r *http.Request) {
	var post Post
	if err := decodeBody(r, &post); err != nil {
		respondErr(w, r, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := post.Create(s.db); err != nil {
		respondErr(w, r, http.StatusInternalServerError, "Failed to create post")
		return
	}

	respond(w, r, http.StatusCreated, post)
}

func (s *Server) handlePostGet(w http.ResponseWriter, r *http.Request, id string) {
	post, err := GetPostByID(s.db, id)
	if err != nil {
		respondErr(w, r, http.StatusInternalServerError, "Failed to get post")
		return
	}
	if post == nil {
		respondHTTPErr(w, r, http.StatusNotFound)
		return
	}
	respond(w, r, http.StatusOK, post)
}

func (s *Server) handlePostUpdate(w http.ResponseWriter, r *http.Request, id string) {
	var post Post
	if err := decodeBody(r, &post); err != nil {
		respondErr(w, r, http.StatusBadRequest, "Invalid request body")
		return
	}

	post.ID = id

	if err := post.Update(s.db); err != nil {
		respondErr(w, r, http.StatusInternalServerError, "Failed to update post")
		return
	}

	respond(w, r, http.StatusOK, post)
}

func (s *Server) handlePostDelete(w http.ResponseWriter, r *http.Request, id string) {
	post := &Post{ID: id}
	if err := post.Delete(s.db); err != nil {
		respondErr(w, r, http.StatusInternalServerError, "Failed to delete post")
		return
	}

	respond(w, r, http.StatusNoContent, nil)
}

func (p *Post) Create(db *sql.DB) error {
	p.ID = generateUUID()
	return withTransaction(db, func(tx *sql.Tx) error {
		query := `INSERT INTO posts (id, board_id, user_id, content, reply_to, created_at) VALUES (?, ?, ?, ?, ?, ?)`
		_, err := tx.Exec(query, p.ID, p.BoardID, p.UserID, p.Content, p.ReplyTo, time.Now())
		if err != nil {
			return err
		}

		return p.updateTags(tx)
	})
}

func (p *Post) Update(db *sql.DB) error {
	return withTransaction(db, func(tx *sql.Tx) error {
		query := `UPDATE posts SET content = ?, reply_to = ? WHERE id = ?`
		_, err := tx.Exec(query, p.Content, p.ReplyTo, p.ID)
		if err != nil {
			return err
		}

		return p.updateTags(tx)
	})
}

func (p *Post) Delete(db *sql.DB) error {
	return withTransaction(db, func(tx *sql.Tx) error {
		if _, err := tx.Exec("DELETE FROM post_tags WHERE post_id = ?", p.ID); err != nil {
			return err
		}

		_, err := tx.Exec("DELETE FROM posts WHERE id = ?", p.ID)
		return err
	})
}

func (p *Post) updateTags(tx *sql.Tx) error {
	if _, err := tx.Exec("DELETE FROM post_tags WHERE post_id = ?", p.ID); err != nil {
		return err
	}

	for _, tagName := range p.Tags {
		if err := addPostTag(tx, p.ID, tagName); err != nil {
			return err
		}
	}
	return nil
}

func addPostTag(tx *sql.Tx, postID string, tagName string) error {
	var tagID string
	err := tx.QueryRow("SELECT id FROM tags WHERE name = ?", tagName).Scan(&tagID)
	if err == sql.ErrNoRows {
		tagID = generateUUID()
		_, err := tx.Exec("INSERT INTO tags (id, name) VALUES (?, ?)", tagID, tagName)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	_, err = tx.Exec("INSERT INTO post_tags (post_id, tag_id) VALUES (?, ?)", postID, tagID)
	return err
}

func getPostTags(db *sql.DB, p *Post) error {
	query := `
		SELECT t.name
		FROM tags t
		JOIN post_tags pt ON t.id = pt.tag_id
		WHERE pt.post_id = ?
	`
	rows, err := queryRows(db, query, p.ID)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var tagName string
		if err := rows.Scan(&tagName); err != nil {
			return err
		}
		p.Tags = append(p.Tags, tagName)
	}

	return nil
}

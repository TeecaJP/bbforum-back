package main

import (
	"database/sql"
	"net/http"
)

func (s *Server) HandlePosts(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.handlePostsGet(w, r)
	default:
		respondHTTPErr(w, r, http.StatusMethodNotAllowed)
	}
}

func (s *Server) handlePostsGet(w http.ResponseWriter, r *http.Request) {
	boardID := r.URL.Query().Get("board_id")
	tag := r.URL.Query().Get("tag")

	var boardIDPtr *string
	if boardID != "" {
		if !isValidUUID(boardID) {
			respondErr(w, r, http.StatusBadRequest, "Invalid board ID")
			return
		}
		boardIDPtr = &boardID
	}

	var tagPtr *string
	if tag != "" {
		tagPtr = &tag
	}

	posts, err := GetPosts(s.db, boardIDPtr, tagPtr)
	if err != nil {
		respondErr(w, r, http.StatusInternalServerError, "Failed to get posts")
		return
	}

	respond(w, r, http.StatusOK, posts)
}

func GetPosts(db *sql.DB, boardID *string, tag *string) ([]Post, error) {
	query, args := buildGetPostsQuery(boardID, tag)

	rows, err := queryRows(db, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanPosts(db, rows)
}

func GetPostByID(db *sql.DB, id string) (*Post, error) {
	query := `SELECT id, board_id, user_id, content, reply_to, created_at FROM posts WHERE id = ?`
	row := queryRow(db, query, id)

	post, err := scanPost(row)
	if err != nil {
		return nil, err
	}

	if err := getPostTags(db, post); err != nil {
		return nil, err
	}

	return post, nil
}

func buildGetPostsQuery(boardID *string, tag *string) (string, []interface{}) {
	query := `
		SELECT DISTINCT p.id, p.board_id, p.user_id, p.content, p.reply_to, p.created_at
		FROM posts p
	`
	var args []interface{}

	if tag != nil {
		query += `
			JOIN post_tags pt ON p.id = pt.post_id
			JOIN tags t ON pt.tag_id = t.id
		`
	}

	query += " WHERE 1=1"

	if boardID != nil {
		query += " AND p.board_id = ?"
		args = append(args, *boardID)
	}

	if tag != nil {
		query += " AND t.name = ?"
		args = append(args, *tag)
	}

	query += " ORDER BY p.created_at DESC"

	return query, args
}

func scanPosts(db *sql.DB, rows *sql.Rows) ([]Post, error) {
	var posts []Post
	for rows.Next() {
		post, err := scanPost(rows)
		if err != nil {
			return nil, err
		}

		if err := getPostTags(db, post); err != nil {
			return nil, err
		}

		posts = append(posts, *post)
	}
	return posts, nil
}

func scanPost(scanner interface{ Scan(...interface{}) error }) (*Post, error) {
	var p Post
	var replyTo sql.NullString
	err := scanner.Scan(&p.ID, &p.BoardID, &p.UserID, &p.Content, &replyTo, &p.CreatedAt)
	if err != nil {
		return nil, err
	}
	if replyTo.Valid {
		p.ReplyTo = &replyTo.String
	}
	return &p, nil
}

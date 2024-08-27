package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/Nikhils-179/connection-pool/db"
)

type Response struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	UserName string `json:"username"`
	Bio      string `json:"bio"`
}

type Follow struct {
	SourceID int
	TargetID int
}

func Handler(w http.ResponseWriter, r *http.Request) {
	sqlDB := db.OpenDBConnection()
	userID, err := strconv.Atoi(r.Header.Get("x-user-id"))
	if err != nil {
		log.Println("Error fetching user-id")
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	users, err := GetListFollows(r.Context(), userID, sqlDB)
	if err != nil {
		resp := &Response{
			Status:  http.StatusInternalServerError,
			Message: "ERROR",
		}
		bytes, _ := json.Marshal(resp)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(bytes)
		return
	}

	resp := &Response{
		Status:  http.StatusOK,
		Message: "SUCCESS",
		Data:    users,
	}

	bytes, _ := json.Marshal(resp)
	w.WriteHeader(http.StatusOK)
	w.Write(bytes)
}

func GetListFollows(ctx context.Context, id int, db *sql.DB) ([]*User, error) {
	users := []*User{}
	query := "SELECT source_id, target_id FROM follows WHERE source_id = ?"
	rows, err := db.QueryContext(ctx, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		follow := &Follow{}
		if err := rows.Scan(&follow.SourceID, &follow.TargetID); err != nil {
			return nil, err
		}

		user, err := GetUserDetails(ctx, follow.TargetID, db)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}
	return users, nil
}

func HandlerFix(w http.ResponseWriter, r *http.Request) {
	sqlDB := db.OpenDBConnection()
	userID, err := strconv.Atoi(r.Header.Get("x-user-id"))
	if err != nil {
		log.Println("Error fetching user-id")
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	follows, err := GetListFollowsFix(r.Context(), userID, sqlDB)
	if err != nil {
		resp := &Response{
			Status:  http.StatusInternalServerError,
			Message: "ERROR",
		}
		bytes, _ := json.Marshal(resp)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(bytes)
		return
	}

	users := []*User{}
	for _, follow := range follows {
		user, err := GetUserDetails(r.Context(), follow.TargetID, sqlDB)
		if err != nil {
			resp := &Response{
				Status:  http.StatusInternalServerError,
				Message: "ERROR",
			}
			bytes, _ := json.Marshal(resp)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(bytes)
			return
		}
		users = append(users, user)
	}

	resp := &Response{
		Status:  http.StatusOK,
		Message: "SUCCESS",
		Data:    users,
	}
	bytes, _ := json.Marshal(resp)
	w.WriteHeader(http.StatusOK)
	w.Write(bytes)
}

func GetListFollowsFix(ctx context.Context, id int, db *sql.DB) ([]*Follow, error) {
	follows := []*Follow{}
	query := "SELECT source_id, target_id FROM follows WHERE source_id = ?"
	rows, err := db.QueryContext(ctx, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		follow := &Follow{}
		if err := rows.Scan(&follow.SourceID, &follow.TargetID); err != nil {
			return nil, err
		}
		follows = append(follows, follow)
	}
	return follows, nil
}

func GetUserDetails(ctx context.Context, id int, db *sql.DB) (*User, error) {
	user := &User{}
	query := "SELECT id, name, username, bio FROM users WHERE id = ?"
	row := db.QueryRowContext(ctx, query, id)
	err := row.Scan(&user.ID, &user.Name, &user.UserName, &user.Bio)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}
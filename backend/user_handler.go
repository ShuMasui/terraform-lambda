package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// --- 型定義ブロック ---
type UserScanner interface {
	scanUsers(ctx context.Context) ([]User, error)
	putUser(ctx context.Context, user *User) error
}

type Service struct {
	Table UserScanner
}

type ApiUsers struct {
	UserID   string `json:"user_id"`
	UserName string `json:"user_name"`
}

func NewService(us *TableBasic) *Service {
	if us == nil {
		return nil
	}

	return &Service{Table: us}
}

func (s *Service) GetUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	users, err := s.Table.scanUsers(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("%v", err)))
		return
	}

	response := make([]ApiUsers, 0, len(users))
	for _, u := range users {
		response = append(response, ApiUsers{
			UserID:   u.UserID,
			UserName: u.UserName,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *Service) PostUsers(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	defer r.Body.Close()

	// --- 変数宣言ブロック ---
	var user ApiUsers
	var err error

	if err = json.NewDecoder(r.Body).Decode(&user); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error": "Invalid Json format: %s"}`, err.Error())
		return
	}

	if err = s.Table.putUser(ctx, &User{UserID: user.UserID, UserName: user.UserName}); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error": "Failed Store for DynamoDB: %s"}`, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

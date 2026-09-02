package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

// モック用の構造体
type mockUserScanner struct {
	scanUsersFunc func(ctx context.Context) ([]User, error)
	putUserFunc   func(ctx context.Context, user *User) error
}

func (m *mockUserScanner) scanUsers(ctx context.Context) ([]User, error) {
	if m.scanUsersFunc != nil {
		return m.scanUsersFunc(ctx)
	}
	return nil, nil
}

func (m *mockUserScanner) putUser(ctx context.Context, user *User) error {
	if m.putUserFunc != nil {
		return m.putUserFunc(ctx, user)
	}
	return nil
}

func TestService_GetUsers(t *testing.T) {
	tests := []struct {
		name           string
		mockScan       func(ctx context.Context) ([]User, error)
		expectedStatus int
		expectedBody   []ApiUsers
		expectedHeader string
	}{
		{
			name: "正常系：複数ユーザーが存在する場合、200 OKとユーザー一覧JSONを返す",
			mockScan: func(ctx context.Context) ([]User, error) {
				return []User{
					{UserID: "usr-001", UserName: "Alice"},
					{UserID: "usr-002", UserName: "Bob"},
				}, nil
			},
			expectedStatus: http.StatusOK,
			expectedBody: []ApiUsers{
				{UserID: "usr-001", UserName: "Alice"},
				{UserID: "usr-002", UserName: "Bob"},
			},
			expectedHeader: "application/json",
		},
		{
			name: "正常系：ユーザーが0件の場合、200 OKと空配列JSONを返す",
			mockScan: func(ctx context.Context) ([]User, error) {
				return []User{}, nil
			},
			expectedStatus: http.StatusOK,
			expectedBody:   []ApiUsers{},
			expectedHeader: "application/json",
		},
		{
			name: "異常系：DBエラーが発生した場合、500 Internal Server Errorを返す",
			mockScan: func(ctx context.Context) ([]User, error) {
				return nil, errors.New("dynamodb scan error")
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Service にモックを注入
			service := &Service{
				Table: &mockUserScanner{
					scanUsersFunc: tt.mockScan,
				},
			}

			// テスト用リクエストとレコーダーを作成
			req := httptest.NewRequest(http.MethodGet, "/users", nil)
			rec := httptest.NewRecorder()

			// ハンドラーを実行
			service.GetUsers(rec, req)

			// 1. ステータスコードの検証
			if rec.Code != tt.expectedStatus {
				t.Errorf("status code = %d, want %d", rec.Code, tt.expectedStatus)
			}

			// 2. ヘッダーの検証（指定がある場合）
			if tt.expectedHeader != "" {
				contentType := rec.Header().Get("Content-Type")
				if contentType != tt.expectedHeader {
					t.Errorf("Content-Type = %s, want %s", contentType, tt.expectedHeader)
				}
			}

			// 3. レスポンスボディの検証
			if tt.expectedBody != nil {
				var actualBody []ApiUsers
				if err := json.NewDecoder(rec.Body).Decode(&actualBody); err != nil {
					t.Fatalf("failed to decode response body: %v", err)
				}

				if !reflect.DeepEqual(actualBody, tt.expectedBody) {
					t.Errorf("response body = %+v, want %+v", actualBody, tt.expectedBody)
				}
			}
		})
	}
}

func TestService_PostUsers(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    string
		mockPut        func(ctx context.Context, user *User) error
		expectedStatus int
		expectedBody   *ApiUsers
		expectedHeader string
	}{
		{
			name:        "正常系：ユーザー作成成功時、200 OKと作成されたユーザーJSONを返す",
			requestBody: `{"user_id": "usr-001", "user_name": "Alice"}`,
			mockPut: func(ctx context.Context, user *User) error {
				if user.UserID != "usr-001" || user.UserName != "Alice" {
					return errors.New("unexpected user payload")
				}
				return nil
			},
			expectedStatus: http.StatusOK,
			expectedBody: &ApiUsers{
				UserID:   "usr-001",
				UserName: "Alice",
			},
			expectedHeader: "application/json",
		},
		{
			name:           "異常系：無効なJSONフォーマットの場合、400 Bad Requestを返す",
			requestBody:    `{invalid json}`,
			mockPut:        nil,
			expectedStatus: http.StatusBadRequest,
			expectedHeader: "application/json",
		},
		{
			name:        "異常系：DB保存時にエラーが発生した場合、500 Internal Server Errorを返す",
			requestBody: `{"user_id": "usr-001", "user_name": "Alice"}`,
			mockPut: func(ctx context.Context, user *User) error {
				return errors.New("dynamodb put error")
			},
			expectedStatus: http.StatusInternalServerError,
			expectedHeader: "application/json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &Service{
				Table: &mockUserScanner{
					putUserFunc: tt.mockPut,
				},
			}

			req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBufferString(tt.requestBody))
			rec := httptest.NewRecorder()

			service.PostUsers(rec, req)

			// 1. ステータスコードの検証
			if rec.Code != tt.expectedStatus {
				t.Errorf("status code = %d, want %d", rec.Code, tt.expectedStatus)
			}

			// 2. ヘッダーの検証（指定がある場合）
			if tt.expectedHeader != "" {
				contentType := rec.Header().Get("Content-Type")
				if contentType != tt.expectedHeader {
					t.Errorf("Content-Type = %s, want %s", contentType, tt.expectedHeader)
				}
			}

			// 3. レスポンスボディの検証
			if tt.expectedBody != nil {
				var actualBody ApiUsers
				if err := json.NewDecoder(rec.Body).Decode(&actualBody); err != nil {
					t.Fatalf("failed to decode response body: %v", err)
				}

				if !reflect.DeepEqual(actualBody, *tt.expectedBody) {
					t.Errorf("response body = %+v, want %+v", actualBody, *tt.expectedBody)
				}
			}
		})
	}
}

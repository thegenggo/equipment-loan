package apitest

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/thegenggo/equipment-loan/api/internal/config"
	"github.com/thegenggo/equipment-loan/api/internal/model"
	"github.com/thegenggo/equipment-loan/api/internal/router"
	"golang.org/x/crypto/bcrypt"
)

func newTestApp(t *testing.T) *gin.Engine {
	t.Helper()

	gin.SetMode(gin.TestMode)
	truncateAll(t)

	cfg := &config.Config{
		JWTSecret: "test-secret-not-used-anywhere-else",
		JWTTTL:    time.Hour,
	}

	return router.Setup(cfg, testDB)
}

func truncateAll(t *testing.T) {
	t.Helper()

	statements := []string{
		"SET FOREIGN_KEY_CHECKS = 0",
		"TRUNCATE TABLE loan_requests",
		"TRUNCATE TABLE equipments",
		"TRUNCATE TABLE users",
		"SET FOREIGN_KEY_CHECKS = 1",
	}
	for _, statement := range statements {
		if _, err := testDB.Exec(statement); err != nil {
			t.Fatalf("truncate: %s: %v", statement, err)
		}
	}
}

func seedUser(t *testing.T, email, password, role string) int64 {
	t.Helper()

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		t.Fatalf("hash password: %v", err)
	}

	result, err := testDB.Exec(
		`INSERT INTO users (email, password_hash, name, role) VALUES (?, ?, ?, ?)`,
		email, string(hash), email, role,
	)
	if err != nil {
		t.Fatalf("seed user %s: %v", email, err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		t.Fatalf("seed user id: %v", err)
	}

	return id
}

func seedEquipment(t *testing.T, code, status string) int64 {
	t.Helper()

	result, err := testDB.Exec(
		`INSERT INTO equipments (code, name, category, status) VALUES (?, ?, ?, ?)`,
		code, code, "test", status,
	)

	if err != nil {
		t.Fatalf("seed equipment %s: %v", code, err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		t.Fatalf("seed equipment id: %v", err)
	}

	return id
}

func signIn(t *testing.T, app *gin.Engine, email, password string) string {
	t.Helper()

	response := do(t, app, http.MethodPost, "/api/v1/auth/login", "", map[string]string{
		"email":    email,
		"password": password,
	})
	if response.Code != http.StatusOK {
		t.Fatalf("sign in as %s: got %d, body %s", email, response.Code, response.Body)
	}

	var payload struct {
		Token string `json:"token"`
	}
	if err := json.Unmarshal(response.Body.Bytes(), &payload); err != nil {
		t.Fatalf("decode login response: %v", err)
	}

	return payload.Token
}

func do(t *testing.T, app *gin.Engine, method, path, token string, body any) *httptest.ResponseRecorder {
	t.Helper()

	var reader *bytes.Reader
	if body == nil {
		reader = bytes.NewReader(nil)
	} else {
		encoded, err := json.Marshal(body)
		if err != nil {
			t.Fatalf("encode request body: %v", err)
		}
		reader = bytes.NewReader(encoded)
	}

	request := httptest.NewRequest(method, path, reader)
	if body != nil {
		request.Header.Set("Content-Type", "application/json")
	}
	if token != "" {
		request.Header.Set("Authorization", "Bearer "+token)
	}

	recorder := httptest.NewRecorder()
	app.ServeHTTP(recorder, request)

	return recorder
}

func borrow(t *testing.T, app *gin.Engine, token string, equipmentID int64) *httptest.ResponseRecorder {
	t.Helper()

	return do(t, app, http.MethodPost, "/api/v1/loans", token, map[string]any{
		"equipment_id": equipmentID,
		"purpose":      "integration test",
	})
}

func staffAndAdmin(t *testing.T, app *gin.Engine) (staffToken, adminToken string) {
	t.Helper()

	seedUser(t, "staff@test.local", "password123", model.RoleStaff)
	seedUser(t, "admin@test.local", "password123", model.RoleAdmin)

	return signIn(t, app, "staff@test.local", "password123"),
		signIn(t, app, "admin@test.local", "password123")
}

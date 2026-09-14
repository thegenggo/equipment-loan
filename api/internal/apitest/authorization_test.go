package apitest

import (
	"encoding/json"
	"net/http"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/thegenggo/equipment-loan/api/internal/model"
)

func TestStaffCannotReachAdminEndpoint(t *testing.T) {
	app := newTestApp(t)
	staffToken, adminToken := staffAndAdmin(t, app)
	equipmentID := seedEquipment(t, "NB-200", model.EquipmentAvailable)

	loanID := createLoanID(t, app, staffToken, equipmentID)

	cases := []struct {
		name   string
		method string
		path   string
		body   any
	}{
		{"create equipment", http.MethodPost, "/api/v1/equipments", map[string]string{"code": "X-1", "name": "x", "category": "x"}},
		{"update equipment", http.MethodPut, "/api/v1/equipments/1", map[string]string{"code": "X-1", "name": "x", "category": "x", "status": "available"}},
		{"delete equipment", http.MethodDelete, "/api/v1/equipments/1", nil},
		{"approve loan", http.MethodPatch, "/api/v1/loans/1/approve", nil},
		{"reject loan", http.MethodPatch, "/api/v1/loans/1/reject", nil},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			response := do(t, app, testCase.method, testCase.path, staffToken, testCase.body)

			if response.Code != http.StatusForbidden {
				t.Fatalf("want 403, got %d, body %s", response.Code, response.Body)
			}
		})

		t.Run(testCase.name+" without a token", func(t *testing.T) {
			response := do(t, app, testCase.method, testCase.path, "", testCase.body)

			if response.Code != http.StatusUnauthorized {
				t.Fatalf("want 401, got %d, body %s", response.Code, response.Body)
			}
		})
	}

	t.Run("an administrator is allowed through", func(t *testing.T) {
		response := do(t, app, http.MethodPatch,
			"/api/v1/loans/"+itoa(loanID)+"/approve", adminToken, nil)

		if response.Code != http.StatusOK {
			t.Fatalf("want 200, got %d, body %s", response.Code, response.Body)
		}
	})
}

func TestStaffCannotTouchAnotherStaffsRequest(t *testing.T) {
	app := newTestApp(t)
	seedUser(t, "owner@test.local", "password123", model.RoleStaff)
	seedUser(t, "stranger@test.local", "password123", model.RoleStaff)
	seedUser(t, "admin@test.local", "password123", model.RoleAdmin)
	ownerToken := signIn(t, app, "owner@test.local", "password123")
	strangerToken := signIn(t, app, "stranger@test.local", "password123")
	adminToken := signIn(t, app, "admin@test.local", "password123")

	equipmentID := seedEquipment(t, "NB-300", model.EquipmentAvailable)
	loanID := createLoanID(t, app, ownerToken, equipmentID)
	path := "/api/v1/loans/" + itoa(loanID)

	t.Run("reading it answers not found, not forbidden", func(t *testing.T) {
		response := do(t, app, http.MethodGet, path, strangerToken, nil)

		if response.Code != http.StatusNotFound {
			t.Fatalf("want 404, got %d, body %s", response.Code, response.Body)
		}
	})

	t.Run("it is absent from their listing", func(t *testing.T) {
		response := do(t, app, http.MethodGet, "/api/v1/loans", strangerToken, nil)

		if response.Code != http.StatusOK {
			t.Fatalf("want 200, got %d", response.Code)
		}
		if got := response.Body.String(); got != "[]" {
			t.Fatalf("want an empty list, got %s", got)
		}
	})

	t.Run("returning it is refused", func(t *testing.T) {
		if response := do(t, app, http.MethodPatch,
			path+"/approve", adminToken, nil); response.Code != http.StatusOK {
			t.Fatalf("approve as admin: want 200, got %d, body %s", response.Code, response.Body)
		}

		response := do(t, app, http.MethodPatch, path+"/return", strangerToken, nil)

		if response.Code != http.StatusForbidden {
			t.Fatalf("want 403, got %d, body %s", response.Code, response.Body)
		}
		assertOpenRequestCount(t, equipmentID, 1)
	})

	t.Run("the owner can still return it", func(t *testing.T) {
		response := do(t, app, http.MethodPatch, path+"/return", ownerToken, nil)

		if response.Code != http.StatusOK {
			t.Fatalf("want 200, got %d, body %s", response.Code, response.Body)
		}
		assertOpenRequestCount(t, equipmentID, 0)
	})
}

func createLoanID(t *testing.T, app *gin.Engine, token string, equipmentID int64) int64 {
	t.Helper()

	response := borrow(t, app, token, equipmentID)
	if response.Code != http.StatusCreated {
		t.Fatalf("create loan: want 201, got %d, body %s", response.Code, response.Body)
	}

	var payload struct {
		ID int64 `json:"id"`
	}
	if err := json.Unmarshal(response.Body.Bytes(), &payload); err != nil {
		t.Fatalf("decode loan response: %v", err)
	}

	return payload.ID
}

func itoa(id int64) string {
	return strconv.FormatInt(id, 10)
}

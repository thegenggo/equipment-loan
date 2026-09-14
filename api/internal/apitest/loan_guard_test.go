package apitest

import (
	"net/http"
	"sync"
	"testing"

	"github.com/thegenggo/equipment-loan/api/internal/model"
)

func TestCannotBorrowUnavailableEquipment(t *testing.T) {
	t.Run("equipment already out on loan", func(t *testing.T) {
		app := newTestApp(t)
		staffToken, _ := staffAndAdmin(t, app)
		equipmentID := seedEquipment(t, "NB-100", model.EquipmentBorrowed)

		response := borrow(t, app, staffToken, equipmentID)

		if response.Code != http.StatusConflict {
			t.Fatalf("want 409, got %d, body %s", response.Code, response.Body)
		}
	})

	t.Run("equipment already requested by someone else", func(t *testing.T) {
		app := newTestApp(t)
		seedUser(t, "first@test.local", "password123", model.RoleStaff)
		seedUser(t, "second@test.local", "password123", model.RoleStaff)
		firstToken := signIn(t, app, "first@test.local", "password123")
		secondToken := signIn(t, app, "second@test.local", "password123")
		equipmentID := seedEquipment(t, "NB-101", model.EquipmentAvailable)

		if first := borrow(t, app, firstToken, equipmentID); first.Code != http.StatusCreated {
			t.Fatalf("first request: want 201, got %d, body %s", first.Code, first.Body)
		}

		second := borrow(t, app, secondToken, equipmentID)

		if second.Code != http.StatusConflict {
			t.Fatalf("second request: want 409, got %d, body %s", second.Code, second.Body)
		}
		assertOpenRequestCount(t, equipmentID, 1)
	})

	t.Run("two people borrowing at the same moment", func(t *testing.T) {
		app := newTestApp(t)
		seedUser(t, "racer-a@test.local", "password123", model.RoleStaff)
		seedUser(t, "racer-b@test.local", "password123", model.RoleStaff)
		tokenA := signIn(t, app, "racer-a@test.local", "password123")
		tokenB := signIn(t, app, "racer-b@test.local", "password123")
		equipmentID := seedEquipment(t, "NB-102", model.EquipmentAvailable)

		start := make(chan struct{})
		codes := make([]int, 2)

		var group sync.WaitGroup
		group.Add(2)

		for index, token := range []string{tokenA, tokenB} {
			go func() {
				defer group.Done()
				<-start
				codes[index] = borrow(t, app, token, equipmentID).Code
			}()
		}

		close(start)
		group.Wait()

		created, conflicted := 0, 0
		for _, code := range codes {
			switch code {
			case http.StatusCreated:
				created++
			case http.StatusConflict:
				conflicted++
			default:
				t.Fatalf("unexpected status %d", code)
			}
		}

		if created != 1 || conflicted != 1 {
			t.Fatalf("want exactly one 201 and one 409, got %v", codes)
		}
		assertOpenRequestCount(t, equipmentID, 1)
	})
}

// assertOpenRequestCount fails unless the item has exactly the expected number
// of pending or approved requests against it.
func assertOpenRequestCount(t *testing.T, equipmentID int64, want int) {
	t.Helper()

	var got int
	err := testDB.Get(&got,
		`SELECT COUNT(*) FROM loan_requests
		WHERE equipment_id = ? AND status IN (?, ?)`,
		equipmentID, model.LoanPending, model.LoanApproved,
	)
	if err != nil {
		t.Fatalf("count open requests: %v", err)
	}

	if got != want {
		t.Fatalf("open requests for equipment %d: want %d, got %d", equipmentID, want, got)
	}
}

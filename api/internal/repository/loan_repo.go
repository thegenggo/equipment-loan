package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/thegenggo/equipment-loan/api/internal/model"
)

var ErrLoanNotFound = errors.New("loan request not found")

type LoanFilter struct {
	UserID int64
	Status string
}

type LoanRepository struct {
	db *sqlx.DB
}

func NewLoanRepository(db *sqlx.DB) *LoanRepository {
	return &LoanRepository{db: db}
}

const detailQuery = `
	SELECT l.id, l.user_id, l.equipment_id, l.purpose, l.status,
		l.requested_at, l.approved_by, l.approved_at, l.returned_at,
		e.code AS equipment_code, e.name AS equipment_name,
		u.name AS user_name
	FROM loan_requests l
	JOIN equipments e ON e.id = l.equipment_id
	JOIN users u ON u.id = l.user_id`

func (r *LoanRepository) List(ctx context.Context, filter LoanFilter) ([]model.LoanRequestDetail, error) {
	query := detailQuery
	args := []any{}

	conditions := []string{}
	if filter.UserID != 0 {
		conditions = append(conditions, "l.user_id = ?")
		args = append(args, filter.UserID)
	}
	if filter.Status != "" {
		conditions = append(conditions, "l.status = ?")
		args = append(args, filter.Status)
	}
	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	query += " ORDER BY l.requested_at DESC"

	loans := []model.LoanRequestDetail{}
	if err := r.db.SelectContext(ctx, &loans, query, args...); err != nil {
		return nil, fmt.Errorf("select loan requests: %w", err)
	}

	return loans, nil
}

func (r *LoanRepository) FindDetailByID(ctx context.Context, id int64) (*model.LoanRequestDetail, error) {
	var loan model.LoanRequestDetail
	if err := r.db.GetContext(ctx, &loan, detailQuery+" WHERE l.id = ?", id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrLoanNotFound
		}
		return nil, fmt.Errorf("select loan requests by id: %w", err)
	}

	return &loan, nil
}

func (r *LoanRepository) FindByIDForUpdate(ctx context.Context, exec Executor, id int64) (*model.LoanRequest, error) {
	const query = `
		SELECT id, user_id, equipment_id, purpose, status,
			requested_at, approved_by, approved_at, returned_at
		FROM loan_requests
		WHERE id = ?
		FOR UPDATE`

	var loan model.LoanRequest
	if err := exec.GetContext(ctx, &loan, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrLoanNotFound
		}
		return nil, fmt.Errorf("select loan request for update: %w", err)
	}

	return &loan, nil
}

func (r *LoanRepository) HasActiveForEquipment(ctx context.Context, exec Executor, equipmentID int64) (bool, error) {
	const query = `
		SELECT EXISTS (
			SELECT 1 FROM loan_requests
			WHERE equipment_id = ? AND status IN (?, ?)
		)`

	var exists bool
	if err := exec.GetContext(ctx, &exists, query, equipmentID, model.LoanPending, model.LoanApproved); err != nil {
		return false, fmt.Errorf("check active loan: %w", err)
	}

	return exists, nil
}

func (r *LoanRepository) Create(ctx context.Context, exec Executor, loan *model.LoanRequest) (int64, error) {
	const query = `
		INSERT INTO loan_requests (user_id, equipment_id, purpose, status, requested_at)
		VALUES (:user_id, :equipment_id, :purpose, :status, :requested_at)`

	result, err := exec.NamedExecContext(ctx, query, loan)
	if err != nil {
		return 0, fmt.Errorf("insert loan request: %w", translate(err))
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("last insert id: %w", err)
	}

	return id, nil
}

func (r *LoanRepository) Decide(ctx context.Context, exec Executor, id, adminID int64, status string, at time.Time) error {
	const query = `
		UPDATE loan_requests
		SET status = ?, approved_by = ?, approved_at = ?
		WHERE id = ? AND status = ?`

	result, err := exec.ExecContext(ctx, query, status, adminID, at, id, model.LoanPending)
	if err != nil {
		return fmt.Errorf("update loan decision: %w", err)
	}

	return expectOneRow(result)
}

func (r *LoanRepository) MarkReturned(ctx context.Context, exec Executor, id int64, at time.Time) error {
	const query = `
		UPDATE loan_requests
		SET status = ?, returned_at = ?
		WHERE id = ? AND status = ?`

	result, err := exec.ExecContext(ctx, query, model.LoanReturned, at, id, model.LoanApproved)
	if err != nil {
		return fmt.Errorf("update loan return: %w", err)
	}

	return expectOneRow(result)
}

func expectOneRow(result sql.Result) error {
	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}
	if affected == 0 {
		return ErrLoanNotFound
	}

	return nil
}

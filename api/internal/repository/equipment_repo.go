package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/thegenggo/equipment-loan/api/internal/model"
)

var ErrEquipmentNotFound = errors.New("equipment not found")

type EquipmentRepository struct {
	db *sqlx.DB
}

func NewEquipmentRepository(db *sqlx.DB) *EquipmentRepository {
	return &EquipmentRepository{db: db}
}

func (r *EquipmentRepository) List(ctx context.Context, status string) ([]model.Equipment, error) {
	query := `
		SELECT id, code, name, category, status
		FROM equipments`

	args := []any{}

	if status != "" {
		query += ` WHERE status = ?`
		args = append(args, status)
	}

	query += ` ORDER BY code`

	equipments := []model.Equipment{}
	if err := r.db.SelectContext(ctx, &equipments, query, args...); err != nil {
		return nil, fmt.Errorf("select equipments: %w", err)
	}

	return equipments, nil
}

func (r *EquipmentRepository) FindByID(ctx context.Context, id int64) (*model.Equipment, error) {
	const query = `
		SELECT id, code, name, category, status
		FROM equipments
		WHERE id = ?`

	var equipment model.Equipment
	if err := r.db.GetContext(ctx, &equipment, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrEquipmentNotFound
		}
		return nil, fmt.Errorf("select equipment by id: %w", err)
	}

	return &equipment, nil
}

func (r *EquipmentRepository) Create(ctx context.Context, equipment *model.Equipment) (int64, error) {
	const query = `
		INSERT INTO equipments (code, name, category, status)
		VALUES (:code, :name, :category, :status)`

	result, err := r.db.NamedExecContext(ctx, query, equipment)
	if err != nil {
		return 0, fmt.Errorf("insert equipment: %w", translate(err))
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("last insert id: %w", err)
	}

	return id, nil
}

func (r *EquipmentRepository) Update(ctx context.Context, equipment *model.Equipment) error {
	const query = `
		UPDATE equipments
		SET code = :code, name = :name, category = :category, status = :status
		WHERE id = :id`

	if _, err := r.db.NamedExecContext(ctx, query, equipment); err != nil {
		return fmt.Errorf("update equipment: %w", translate(err))
	}

	return nil
}

func (r *EquipmentRepository) Delete(ctx context.Context, id int64) error {
	const query = `DELETE FROM equipments WHERE id = ?`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete equipment: %w", translate(err))
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}
	if affected == 0 {
		return ErrEquipmentNotFound
	}

	return nil
}

func (r *EquipmentRepository) FindByIDForUpdate(ctx context.Context, exec Executor, id int64) (*model.Equipment, error) {
	const query = `
		SELECT id, code, name, category, status
		FROM equipments
		WHERE id = ?
		FOR UPDATE`

	var equipment model.Equipment
	if err := exec.GetContext(ctx, &equipment, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrEquipmentNotFound
		}
		return nil, fmt.Errorf("select equipment for update: %w", err)
	}

	return &equipment, nil
}

func (r *EquipmentRepository) UpdateStatus(ctx context.Context, exec Executor, id int64, status string) error {
	const query = `UPDATE equipments SET status = ? WHERE id = ?`

	if _, err := exec.ExecContext(ctx, query, status, id); err != nil {
		return fmt.Errorf("update equipment status: %w", err)
	}

	return nil
}

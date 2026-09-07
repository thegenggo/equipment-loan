package service

import (
	"context"
	"errors"

	"github.com/thegenggo/equipment-loan/api/internal/model"
	"github.com/thegenggo/equipment-loan/api/internal/repository"
)

var (
	ErrEquipmentNotFound  = errors.New("equipment not found")
	ErrEquipmentCodeTaken = errors.New("equipment code already exists")
	ErrEquipmentInUse     = errors.New("equipment has loan history")
	ErrEquipmentBorrowed  = errors.New("equipment is currently borrowed")
	ErrInvalidStatus      = errors.New("invalid equipment status")
)

type EquipmentService struct {
	equipments *repository.EquipmentRepository
}

func NewEquipmentService(equipments *repository.EquipmentRepository) *EquipmentService {
	return &EquipmentService{equipments: equipments}
}

func (s *EquipmentService) List(ctx context.Context, status string) ([]model.Equipment, error) {
	if status != "" && !model.IsValidEquipmentStatus(status) {
		return nil, ErrInvalidStatus
	}

	return s.equipments.List(ctx, status)
}

func (s *EquipmentService) Get(ctx context.Context, id int64) (*model.Equipment, error) {
	equipment, err := s.equipments.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrEquipmentNotFound) {
			return nil, ErrEquipmentNotFound
		}
		return nil, err
	}

	return equipment, nil
}

func (s *EquipmentService) Create(ctx context.Context, code, name, category string) (*model.Equipment, error) {
	equipment := &model.Equipment{
		Code:     code,
		Name:     name,
		Category: category,
		Status:   model.EquipmentAvailable,
	}

	id, err := s.equipments.Create(ctx, equipment)
	if err != nil {
		if errors.Is(err, repository.ErrDuplicateEntry) {
			return nil, ErrEquipmentCodeTaken
		}
		return nil, err
	}
	equipment.ID = id

	return equipment, nil
}

func (s *EquipmentService) Update(ctx context.Context, id int64, code, name, category, status string) (*model.Equipment, error) {
	current, err := s.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	if current.Status == model.EquipmentBorrowed {
		return nil, ErrEquipmentBorrowed
	}

	current.Code = code
	current.Name = name
	current.Category = category
	current.Status = status

	if err := s.equipments.Update(ctx, current); err != nil {
		if errors.Is(err, repository.ErrDuplicateEntry) {
			return nil, ErrEquipmentCodeTaken
		}
		return nil, err
	}

	return current, nil
}

func (s *EquipmentService) Delete(ctx context.Context, id int64) error {
	err := s.equipments.Delete(ctx, id)
	switch {
	case errors.Is(err, repository.ErrEquipmentNotFound):
		return ErrEquipmentNotFound
	case errors.Is(err, repository.ErrStillReferenced):
		return ErrEmailTaken
	}

	return err
}

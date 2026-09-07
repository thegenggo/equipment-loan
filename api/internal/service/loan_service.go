package service

import (
	"context"
	"errors"
	"time"

	"github.com/thegenggo/equipment-loan/api/internal/model"
	"github.com/thegenggo/equipment-loan/api/internal/repository"
)

var (
	ErrLoanNotFound          = errors.New("loan request not found")
	ErrLoanNotPending        = errors.New("loan request is no longer pending")
	ErrLoanNotApproved       = errors.New("loan request is not on loan")
	ErrEquipmentNotAvailable = errors.New("equipment is not available")
	ErrEquipmentRequested    = errors.New("equipment already has on open request")
	ErrNotOwner              = errors.New("loan request belongs to someone else")
)

type LoanService struct {
	tx         *repository.TxManager
	loans      *repository.LoanRepository
	equipments *repository.EquipmentRepository
}

func NewLoanService(
	tx *repository.TxManager,
	loans *repository.LoanRepository,
	equipments *repository.EquipmentRepository,
) *LoanService {
	return &LoanService{tx: tx, loans: loans, equipments: equipments}
}

func (s *LoanService) Create(ctx context.Context, userID, equipmentID int64, purpose string) (*model.LoanRequestDetail, error) {
	var loanID int64

	err := s.tx.Run(ctx, func(exec repository.Executor) error {
		equipment, err := s.equipments.FindByIDForUpdate(ctx, exec, equipmentID)
		if err != nil {
			if errors.Is(err, repository.ErrEquipmentNotFound) {
				return ErrEquipmentNotFound
			}
			return err
		}

		if equipment.Status != model.EquipmentAvailable {
			return ErrEquipmentNotAvailable
		}

		active, err := s.loans.HasActiveForEquipment(ctx, exec, equipmentID)
		if err != nil {
			return err
		}
		if active {
			return ErrEquipmentRequested
		}

		loanID, err = s.loans.Create(ctx, exec, &model.LoanRequest{
			UserID:      userID,
			EquipmentID: equipmentID,
			Purpose:     purpose,
			Status:      model.LoanPending,
			RequestedAt: time.Now().UTC(),
		})

		return err
	})
	if err != nil {
		return nil, err
	}

	return s.loans.FindDetailByID(ctx, loanID)
}

func (s *LoanService) List(ctx context.Context, actorID int64, actorRole, status string) ([]model.LoanRequestDetail, error) {
	if status != "" && !isValidLoanStatus(status) {
		return nil, ErrInvalidStatus
	}

	filter := repository.LoanFilter{Status: status}
	if actorRole != model.RoleAdmin {
		filter.UserID = actorID
	}

	return s.loans.List(ctx, filter)
}

func (s *LoanService) Get(ctx context.Context, id, actorID int64, actorRole string) (*model.LoanRequestDetail, error) {
	loan, err := s.loans.FindDetailByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrLoanNotFound) {
			return nil, ErrLoanNotFound
		}
		return nil, err
	}

	if actorRole != model.RoleAdmin && loan.UserID != actorID {
		return nil, ErrLoanNotFound
	}

	return loan, nil
}

func (s *LoanService) Approve(ctx context.Context, id, adminID int64) (*model.LoanRequestDetail, error) {
	err := s.tx.Run(ctx, func(exec repository.Executor) error {
		loan, err := s.loanForUpdate(ctx, exec, id)
		if err != nil {
			return err
		}
		if loan.Status != model.LoanPending {
			return ErrLoanNotPending
		}

		equipment, err := s.equipments.FindByIDForUpdate(ctx, exec, loan.EquipmentID)
		if err != nil {
			return err
		}
		if equipment.Status != model.EquipmentAvailable {
			return ErrEquipmentNotAvailable
		}

		if err := s.equipments.UpdateStatus(ctx, exec, equipment.ID, model.EquipmentBorrowed); err != nil {
			return err
		}

		return s.loans.Decide(ctx, exec, id, adminID, model.LoanApproved, time.Now().UTC())
	})
	if err != nil {
		return nil, err
	}

	return s.loans.FindDetailByID(ctx, id)
}

func (s *LoanService) Reject(ctx context.Context, id, adminID int64) (*model.LoanRequestDetail, error) {
	loan, err := s.loans.FindDetailByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrLoanNotFound) {
			return nil, ErrLoanNotFound
		}
		return nil, err
	}
	if loan.Status != model.LoanPending {
		return nil, ErrLoanNotPending
	}

	err = s.loans.Decide(ctx, s.tx.DB(), id, adminID, model.LoanRejected, time.Now().UTC())
	if err != nil {
		if errors.Is(err, repository.ErrLoanNotFound) {
			return nil, ErrLoanNotPending
		}
		return nil, err
	}

	return s.loans.FindDetailByID(ctx, id)
}

func (s *LoanService) Return(ctx context.Context, id, actorID int64) (*model.LoanRequestDetail, error) {
	err := s.tx.Run(ctx, func(exec repository.Executor) error {
		loan, err := s.loanForUpdate(ctx, exec, id)
		if err != nil {
			return err
		}

		if loan.UserID != actorID {
			return ErrNotOwner
		}
		if loan.Status != model.LoanApproved {
			return ErrLoanNotApproved
		}

		if err := s.loans.MarkReturned(ctx, exec, id, time.Now().UTC()); err != nil {
			return err
		}

		return s.equipments.UpdateStatus(ctx, exec, loan.EquipmentID, model.EquipmentAvailable)
	})
	if err != nil {
		return nil, err
	}

	return s.loans.FindDetailByID(ctx, id)
}

func (s *LoanService) loanForUpdate(ctx context.Context, exec repository.Executor, id int64) (*model.LoanRequest, error) {
	loan, err := s.loans.FindByIDForUpdate(ctx, exec, id)
	if err != nil {
		if errors.Is(err, repository.ErrLoanNotFound) {
			return nil, ErrLoanNotFound
		}
		return nil, err
	}

	return loan, nil
}

func isValidLoanStatus(s string) bool {
	switch s {
	case model.LoanPending, model.LoanApproved, model.LoanRejected, model.LoanReturned:
		return true
	default:
		return false
	}
}

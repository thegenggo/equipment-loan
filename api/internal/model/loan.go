package model

import "time"

const (
	LoanPending  = "pending"
	LoanApproved = "approved"
	LoanRejected = "rejected"
	LoanReturned = "returned"
)

type LoanRequest struct {
	ID          int64      `db:"id" json:"id"`
	UserID      int64      `db:"user_id" json:"user_id"`
	EquipmentID int64      `db:"equipment_id" json:"equipment_id"`
	Purpose     string     `db:"purpose" json:"purpose"`
	Status      string     `db:"status" json:"status"`
	RequestedAt time.Time  `db:"requested_at" json:"requested_at"`
	ApprovedBy  *int64     `db:"approved_by" json:"approved_by"`
	ApprovedAt  *time.Time `db:"approved_at" json:"approved_at"`
	ReturnedAt  *time.Time `db:"returned_at" json:"returned_at"`
}

type LoanRequestDetail struct {
	LoanRequest
	EquipmentCode string `db:"equipment_code" json:"equipment_code"`
	EquipmentName string `db:"equipment_name" json:"equipment_name"`
	UserName      string `db:"user_name" json:"user_name"`
}

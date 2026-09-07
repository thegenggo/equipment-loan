package model

const (
	EquipmentAvailable = "available"
	EquipmentBorrowed  = "borrowed"
	EquipmentRepair    = "repair"
)

type Equipment struct {
	ID       int64  `db:"id" json:"id"`
	Code     string `db:"code" json:"code"`
	Name     string `db:"name" json:"name"`
	Category string `db:"category" json:"category"`
	Status   string `db:"status" json:"status"`
}

func IsValidEquipmentStatus(s string) bool {
	switch s {
	case EquipmentAvailable, EquipmentBorrowed, EquipmentRepair:
		return true
	default:
		return false
	}
}

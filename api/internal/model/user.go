package model

const (
	RoleStaff = "staff"
	RoleAdmin = "admin"
)

type User struct {
	ID           int64  `db:"id" json:"id"`
	Email        string `db:"email" json:"email"`
	PasswordHash string `db:"password_hash" json:"-"`
	Name         string `db:"name" json:"name"`
	Role         string `db:"role" json:"role"`
}

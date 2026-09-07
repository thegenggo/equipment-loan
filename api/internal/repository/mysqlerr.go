package repository

import (
	"errors"

	"github.com/go-sql-driver/mysql"
)

const (
	codeDuplicateEntry  = 1062
	codeRowIsReferenced = 1451
)

var (
	ErrDuplicateEntry  = errors.New("duplicate entry")
	ErrStillReferenced = errors.New("row is still referenced")
)

func translate(err error) error {
	var mysqlErr *mysql.MySQLError
	if !errors.As(err, &mysqlErr) {
		return err
	}

	switch mysqlErr.Number {
	case codeDuplicateEntry:
		return ErrDuplicateEntry
	case codeRowIsReferenced:
		return ErrStillReferenced
	default:
		return err
	}
}

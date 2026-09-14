package apitest

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/jmoiron/sqlx"

	_ "github.com/go-sql-driver/mysql"
)

var testDB *sqlx.DB

func TestMain(m *testing.M) {
	dsn := os.Getenv("TEST_DB_DSN")
	if dsn == "" {
		fmt.Println("apitest: TEST_DB_DSN is not set, skipping integration tests")
		os.Exit(0)
	}

	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		fmt.Printf("apitest: cannot reach the test database: %v\n", err)
		os.Exit(1)
	}
	testDB = db

	if err := rebuildSchema(db); err != nil {
		fmt.Printf("apitest: cannot apply migrations: %v\n", err)
		os.Exit(1)
	}

	code := m.Run()

	db.Close()
	os.Exit(code)
}

func rebuildSchema(db *sqlx.DB) error {
	drops := []string{
		"SET FOREIGN_KEY_CHECKS = 0",
		"DROP TABLE IF EXISTS loan_requests",
		"DROP TABLE IF EXISTS equipments",
		"DROP TABLE IF EXISTS users",
		"SET FOREIGN_KEY_CHECKS = 1",
	}
	for _, statement := range drops {
		if _, err := db.Exec(statement); err != nil {
			return fmt.Errorf("%s: %w", statement, err)
		}
	}

	paths, err := filepath.Glob("../database/migrations/00[123]_*.sql")
	if err != nil {
		return fmt.Errorf("find migrations: %w", err)
	}
	if len(paths) != 3 {
		return fmt.Errorf("expected 3 schema migrations, found %d", len(paths))
	}

	for _, path := range paths {
		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read %s: %w", path, err)
		}
		if _, err := db.Exec(string(content)); err != nil {
			return fmt.Errorf("apply %s: %w", path, err)
		}
	}

	return nil
}

package postgres

import (
	"database/sql"
	_ "embed"
	"fmt"
)

//go:embed sqlc/schema.sql
var schema string

// InitSchema executes the embedded SQL schema against the database.
// It is safe to run on every startup because of "IF NOT EXISTS".
func InitSchema(db *sql.DB) error {
	_, err := db.Exec(schema)
	if err != nil {
		return fmt.Errorf("failed to initialize database schema: %w", err)
	}
	return nil
}

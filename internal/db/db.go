package db

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
)

var CurrentDatabase *sql.DB = nil

func SetDatabase(db *sql.DB) {
	CurrentDatabase = db
}

// Attempt SQL Server Connection through Windows Authentication
func AttemptConnection(server string, database string) (*sql.DB, error) {
	connString := fmt.Sprintf(
		"server=%s;database=%s;TrustServerCertificate=true;authentication=windows",
		server, database,
	)

	fmt.Println("Attempting Database Connection...")
	db, err := sql.Open("sqlserver", connString)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	fmt.Println("Conected Successfully")
	return db, nil
}

func GetTable(table string) (*sql.Rows, error) {
	query := fmt.Sprintf("SELECT * FROM %s", table)
	return RunQuery(query)
}

func RunQuery(query string) (*sql.Rows, error) {
	return CurrentDatabase.Query(query)
}

// loadQueryFromFile reads a SQL file under the queries/ directory.
func loadQueryFromFile(relativePath string) (string, error) {
	if CurrentDatabase == nil {
		return "", fmt.Errorf("database not initialized")
	}

	cleanPath := filepath.Clean(relativePath)
	path := filepath.Join("queries", cleanPath)

	content, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("reading query file %s: %w", path, err)
	}

	return string(content), nil
}

// ExecFromFile executes a non-SELECT statement located in queries/.
// Use SQL named parameters (e.g., sql.Named("id", 1)) to bind inputs.
func ExecFromFile(ctx context.Context, relativePath string, args ...any) (sql.Result, error) {
	query, err := loadQueryFromFile(relativePath)
	if err != nil {
		return nil, err
	}
	return CurrentDatabase.ExecContext(ctx, query, args...)
}

// QueryRowsFromFile runs a SELECT statement located in queries/.
func QueryRowsFromFile(ctx context.Context, relativePath string, args ...any) (*sql.Rows, error) {
	query, err := loadQueryFromFile(relativePath)
	if err != nil {
		return nil, err
	}
	return CurrentDatabase.QueryContext(ctx, query, args...)
}

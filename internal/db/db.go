package db

import (
	"database/sql"
	"fmt"
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
package models

import (
	"database/sql"
	"fmt"
	"strings"
)

func ensureDB(conn *sql.DB) error {
	if conn == nil {
		return fmt.Errorf("no hay conexion a base de datos")
	}
	return nil
}

func requirePositive(name string, value int) error {
	if value <= 0 {
		return fmt.Errorf("%s debe ser mayor a cero", name)
	}
	return nil
}

func requireNonEmpty(name, value string) (string, error) {
	trim := strings.TrimSpace(value)
	if trim == "" {
		return "", fmt.Errorf("%s no puede ser vacio", name)
	}
	return trim, nil
}

func optionalString(value string) sql.NullString {
	trim := strings.TrimSpace(value)
	if trim == "" {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: trim, Valid: true}
}

func optionalInt(value int) sql.NullInt32 {
	if value <= 0 {
		return sql.NullInt32{Valid: false}
	}
	return sql.NullInt32{Int32: int32(value), Valid: true}
}

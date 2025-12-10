package internal

import (
	"database/sql"
	"fmt"
)

type Stringable interface {
	String() string
}

func ListItems[T Stringable](items []T) {
	for i, item := range items {
		fmt.Printf("#%d: %s\n", i+1, item.String())
	}
}

func NullString(str sql.NullString) string {
	if str.Valid {
		return str.String
	} else {
		return "| null |"
	}
}

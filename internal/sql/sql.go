package sql

import (
	"database/sql"
	"fmt"
	"reflect"
)

func ParseRow[T any](rows *sql.Rows) ([]T, error) {
	// Lista inicial de items
	items := []T{}

	// Mientras haya una columna aun
	for rows.Next() {
		// Inicializamos item
		var item T

		// Usamos reflection para obtener los elementos y campos del tipo
		v := reflect.ValueOf(&item).Elem()

		if v.Kind() != reflect.Struct {
			return nil, fmt.Errorf("ParseRow: T debe ser una estructura, se dió %s", v.Kind())
		}

		// Usamos reflection para obtener los elementos y campos del tipo
		numCols := v.NumField()

		// Obtenemos el numero de columnas
		columns := make([]any, numCols)

		// Vamos introduciendo el valor de la columna a cada campo
		for i := 0; i < numCols; i++ {
			field := v.Field(i)
			if !field.CanSet() {
				return nil, fmt.Errorf("ParseRow: El campo %s no esta exportado o no es asignable", v.Field(i))
			}
			columns[i] = field.Addr().Interface()
		}

		// Vemos si la escritura genera un error
		err := rows.Scan(columns...)
		if err != nil {
			return nil, fmt.Errorf("ParseRow: Error escaneando: %v", err)
		}

		// Añadimos a la lista de objetos
		items = append(items, item)
	}

	// Vemos si hubo algun tipo de error en las filas
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

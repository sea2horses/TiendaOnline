package main

import (
	"context"
	"flag"
	"fmt"
	"tienda-online/internal"
	"tienda-online/internal/db"
	"tienda-online/models"

	_ "github.com/microsoft/go-mssqldb"
)

func main() {
	var (
		debug   *bool   = flag.Bool("debug", false, "Shows debug messages")
		server  *string = flag.String("server", "localhost", "The Database Server URL")
		db_name *string = flag.String("database", "master", "The Database Name")
	)
	flag.Parse()

	if *debug {
		fmt.Printf("Debug Mode.\nServer: %s\nDatabase: %s\n", *server, *db_name)
	}

	conn, err := db.AttemptConnection(*server, *db_name)
	internal.Check(err, "No se pudo conectar a la base de datos")
	// Else, make sure we close the connection when done
	defer conn.Close()

	// Set the current connection as the database
	db.SetDatabase(conn)

	// Ejemplo: listar clientes existentes
	manager := models.NewClienteManager(conn)
	clientes, err := manager.List(context.Background())
	internal.Check(err, "No se pudieron leer los clientes")
	fmt.Println("Clientes Encontrados:")

	internal.ListItems(clientes)
}

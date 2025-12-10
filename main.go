package main

import (
	"flag"
	"fmt"
	"tienda-online/internal"
	"tienda-online/internal/db"

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
	internal.Check(err, "Couldn't connect to database")
	// Else, make sure we close the connection when done
	defer conn.Close()

	// Set the current connection as the database
	db.SetDatabase(conn)

	// Now let's list our students
	e := models.EstudianteManager{}
	e.List()
}

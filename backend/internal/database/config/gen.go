package main

import (
	"log"
	"strconv"

	"github.com/darylhjd/oams/backend/internal/env"
	"github.com/go-jet/jet/v2/generator/postgres"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
)

func main() {
	host, err := strconv.Atoi(env.GetDatabasePort())
	if err != nil {
		log.Fatal(err)
	}

	genPath := "../gen"
	dbConn := postgres.DBConnection{
		Host:       env.GetDatabaseHost(),
		Port:       host,
		User:       env.GetDatabaseUser(),
		Password:   env.GetDatabasePassword(),
		SslMode:    env.GetDatabaseSSLMode(),
		DBName:     env.GetDatabaseName(),
		SchemaName: "public",
	}

	err = postgres.Generate(genPath, dbConn)
	if err != nil {
		log.Fatal(err)
	}
}

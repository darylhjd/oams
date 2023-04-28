package main

import (
	"log"

	"github.com/golang-migrate/migrate/v4"
)

func main() {
	// TODO: Test script for migrate library.
	m, err := migrate.New("", "")
	if err != nil {
		log.Fatal(err)
	}

	if err = m.Up(); err != nil {
		log.Fatal(err)
	}
}

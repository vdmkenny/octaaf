package main

import (
	"log"

	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/pop"
)

// DB Global Connection
var DB *pop.Connection

func connectDB() {
	var err error
	env := envy.Get("GO_ENV", "development")
	DB, err = pop.Connect(env)
	if err != nil {
		log.Fatal(err)
	}
	pop.Debug = env == "development"
}

func migrateDB() {
	fileMigrator, err := pop.NewFileMigrator("./migrations", DB)

	if err != nil {
		log.Panic(err)
	}

	err = fileMigrator.Up()

	if err != nil {
		log.Panic(err)
	}
}

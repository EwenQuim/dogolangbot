package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "modernc.org/sqlite"
)

func saveToDatabase(animal int) {

	db, err := sql.Open("sqlite", "./compteur.db")
	if err != nil {
		log.Fatal("cant open db", err)
	}
	defer db.Close()

	var animalSays string
	switch animal {
	case DOG:
		animalSays = "woof"
	case CAT:
		animalSays = "meow"
	case GUINEA_PIG:
		animalSays = "pouic"
	}

	sqlStmt := `CREATE TABLE IF NOT EXISTS commands (

		date text PRIMARY KEY,
	
		command text 
	);
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}

	println("bof")

	_, err = db.Exec(fmt.Sprintf("INSERT INTO commands (date, command) VALUES (\"%v\", \"%v\")", time.Now().Format("2006-01-02 15:04:05"), animalSays))
	if err != nil {
		log.Fatal(err)
	}

}

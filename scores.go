package main

import (
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"sort"
	"time"

	_ "modernc.org/sqlite"
)

func createDb() *sql.DB {
	db, err := sql.Open("sqlite", "./data/compteur.db")
	if err != nil {
		log.Fatal("cant open db", err)
	}

	sqlStmt := `CREATE TABLE IF NOT EXISTS commands (

		date text PRIMARY KEY,
	
		command text 
	);
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return nil
	}

	return db
}

func (dogobot Dogobot) saveToDatabase(animalSays string) error {
	_, err := dogobot.db.Exec(fmt.Sprintf("INSERT INTO commands (date, command) VALUES (\"%v\", \"%v\");", time.Now().Format("2006-01-02 15:04:05.000000000"), animalSays))
	if err != nil {
		return fmt.Errorf("error while inserting: %v", err)
	}

	return nil
}

func (dogobot Dogobot) getScores() string {
	rows, err := dogobot.db.Query("SELECT command, count(*) FROM commands GROUP BY command;")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	total_count := 0

	for rows.Next() {
		var shout string
		var animal_count int
		err = rows.Scan(&shout, &animal_count)
		if err != nil {
			log.Fatal(err)
		}
		dogobot.animals[shout].count = animal_count
		total_count += animal_count
	}
	if rows.Err() != nil {
		slog.Error("error while reading rows", "error", err)
	}

	dogobot.total_calls = total_count

	return dogobot.formatScoresResponse()

}

// Most asked animal (15319 requests) :
// 🐱 55% - Winner ! 🏆
// 🐶 45%
func (dogobot Dogobot) formatScoresResponse() string {
	text := fmt.Sprintf("Most asked (%v requests):", dogobot.total_calls)

	s := make(animalSlice, 0, len(dogobot.animals))
	for _, d := range dogobot.animals {
		s = append(s, d)
	}

	sort.Sort(sort.Reverse(s))

	for _, animal := range s {
		text += fmt.Sprintf("\n%v %.0f%% (%v)", animal.emoji, 100*float64(animal.count)/float64(dogobot.total_calls+1), animal.count)
	}
	return text
}

type animalSlice []*Animal

// Len is part of sort.Interface.
func (d animalSlice) Len() int {
	return len(d)
}

// Swap is part of sort.Interface.
func (d animalSlice) Swap(i, j int) {
	d[i], d[j] = d[j], d[i]
}

// Less is part of sort.Interface. We use count as the value to sort by
func (d animalSlice) Less(i, j int) bool {
	return d[i].count < d[j].count
}

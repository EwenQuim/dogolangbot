package main

import (
	"database/sql"
	"fmt"
	"log"
	"sort"
	"time"

	_ "modernc.org/sqlite"
)

func saveToDatabase(animalSays string) {

	db, err := sql.Open("sqlite", "./compteur.db")
	if err != nil {
		log.Fatal("cant open db", err)
	}
	defer db.Close()

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

	_, err = db.Exec(fmt.Sprintf("INSERT INTO commands (date, command) VALUES (\"%v\", \"%v\");", time.Now().Format("2006-01-02 15:04:05.000000000"), animalSays))
	if err != nil {
		log.Fatal(err)
	}

}

func (dogobot Dogobot) getScores() string {
	db, err := sql.Open("sqlite", "./compteur.db")
	if err != nil {
		log.Fatal("cant open db", err)
	}
	defer db.Close()

	total_count := 0
	rows, err := db.Query("SELECT command, count(*) FROM commands GROUP BY command;")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
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
	dogobot.total_calls = total_count

	err = rows.Err()
	if err != nil {
		log.Fatal("error while updating ", err)
	}

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

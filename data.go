package main

import (
	"database/sql"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

type Data struct {
	db *sql.DB
}

var data Data

const (
	dbFilename = "./s3.sqlite"
)

func init() {
	if _, err := os.Stat(dbFilename); os.IsNotExist(err) {
		os.Create(dbFilename)
	}

	db, err := sql.Open("sqlite3", dbFilename)
	if err != error(nil) {
		exitErrorf("Could not open SQLite database: %v", err)
	}

	if _, err := db.Exec("CREATE TABLE IF NOT EXISTS objects (id INTEGER PRIMARY KEY AUTOINCREMENT, objectPath TEXT, processed TINYINT(1), error TINYINT(1))"); err != error(nil) {
		exitErrorf("Could not create table 'objects': %v", err)
	}

	data.db = db
}

func (d Data) addPath(objectPath string) {
	if _, err := d.db.Exec("INSERT INTO objects (objectPath, processed, error) VALUES (?, ?, ?)", objectPath, 0, 0); err != error(nil) {
		exitErrorf("Failed to insert path in DB: %v", err)
	}
}

type PathRecord struct {
	Id   int
	Path string
}

func (d Data) getUnprocessedPaths(ch chan<- PathRecord) {
	defer close(ch)

	rows, err := d.db.Query("SELECT id, objectPath FROM objects WHERE processed = 0")
	if err != error(nil) {
		exitErrorf("Failed to query objects from DB: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		record := PathRecord{}

		err := rows.Scan(&record.Id, &record.Path)
		if err != error(nil) {
			exitErrorf("Failed to read object record: %v", err)
		}

		ch <- record
	}
}

func (d Data) setPathAsProcessed(id int, processingError error) {
	if _, err := d.db.Exec("UPDATE objects SET processed = ?, error = ? WHERE id = ?", 1, processingError != error(nil), id); err != error(nil) {
		exitErrorf("Failed to set object as processed: %v", err)
	}
}

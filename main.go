package main

import (
	"database/sql"
	"log"
	"os"
	"strconv"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

const (
	ACTION_ADD    = "add"
	ACTION_UPDATE = "updated"
	ACTION_REMOVE = "remove"
	ACTION_LIST   = "list"
)

var ALL_ACTIONS string = strings.Join([]string{ACTION_ADD, ACTION_UPDATE, ACTION_REMOVE, ACTION_LIST}, ", ")

func main() {
	db, err := sql.Open("sqlite3", "player.db")
	check(err, "on DB open")
	defer db.Close()

	createTable := `
		CREATE TABLE IF NOT EXISTS player (id INTEGER PRIMARY KEY, name TEXT, score INTEGER)`
	_, err = db.Exec(createTable)
	check(err, "create table player")

	args := os.Args
	if len(args) < 2 {
		log.Fatalf("Expected an action argument: %s", ALL_ACTIONS)
	}

	action := args[1]
	switch action {
	case ACTION_ADD:
		addPlayer(db, args[2:])
	case ACTION_UPDATE:
		updatePlayer(db, args[2:])
	case ACTION_REMOVE:
		removePlayer(db, args[2:])
	case ACTION_LIST:
		listPlayers(db)
	default:
		log.Fatalf("Unknown action, expected one of: %s", ALL_ACTIONS)
	}
}

func addPlayer(db *sql.DB, args []string) {
	if len(args) < 1 {
		log.Fatalf("Expected name argument")
	}

	score := 0
	if len(args) >= 2 {
		var err error // Declare error first to assign score below.
		score, err = strconv.Atoi(args[1])
		check(err, "convert score argument to int")
	}

	_, err := db.Exec("INSERT INTO player (name, score) VALUES (?, ?)", args[0], score)
	check(err, "insert into player")
}

func removePlayer(db *sql.DB, args []string) {
	if len(args) < 1 {
		log.Fatalf("Expected id argument")
	}

	id, err := strconv.Atoi(args[0])
	check(err, "convert id argument to int")

	result, err := db.Exec("DELETE FROM player WHERE id = ?", id)
	check(err, "delete from player")
	checkRowFound(result, id)
}

func checkRowFound(result sql.Result, id int) {
	rowsAffected, err := result.RowsAffected()
	check(err, "rows affected player")

	if rowsAffected < 1 {
		log.Printf("Could not find any row with id %d\n", id)
	}
}

func updatePlayer(db *sql.DB, args []string) {
	if len(args) < 2 {
		log.Fatalf("Expected id and score argument")
	}

	id, err := strconv.Atoi(args[0])
	check(err, "convert id argument to int")

	score, err := strconv.Atoi(args[1])
	check(err, "convert score argument to int")

	result, err := db.Exec("UPDATE player SET score = ? WHERE id = ?", score, id)
	check(err, "update player")
	checkRowFound(result, id)
}

func listPlayers(db *sql.DB) {
	result, err := db.Query("SELECT id, name, score FROM player ORDER BY id ASC")
	check(err, "select from player")
	defer result.Close()

	log.Println("id | name | score")
	for result.Next() {
		var id int
		var name string
		var score int
		err = result.Scan(&id, &name, &score)
		check(err, "scan player select")

		log.Printf("%d | %s | %d\n", id, name, score)
	}

	err = result.Err()
	check(err, "iterate player select result")
}

func check(err error, context string) {
	if err != nil {
		log.Fatalf("Error on %s: %v", context, err)
	}
}

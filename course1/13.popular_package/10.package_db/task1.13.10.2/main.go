package main

import (
	"database/sql"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	ID       int
	Username string
	Email    string
}

func main() {
	CreateUserTable()
	InsertUser(User{Username: "Alex", Email: "al@mail.ru"})
	fmt.Println(SelectUser(1))
	UpdateUser(User{ID: 1, Username: "Mariya", Email: "ma@mailru"})
}

func CreateUserTable() error {
	db, err := sql.Open("sqlite3", "users.db")
	if err != nil {
		fmt.Println("error creating the table", err)
		return err
	}
	defer db.Close()

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS users(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT,
		email TEXT
	)`)
	if err != nil {
		return err
	}

	return nil
}

func InsertUser(user User) error {
	db, err := sql.Open("sqlite3", "users.db")
	if err != nil {
		return err
	}
	defer db.Close()

	query, args, err := PrepareQuery("insert", "users", user)
	if err != nil {
		return err
	}

	_, err = db.Exec(query, args...)
	if err != nil {
		return err
	}

	return nil
}

func SelectUser(userID int) (User, error) {
	db, err := sql.Open("sqlite3", "users.db")
	if err != nil {
		return User{}, err
	}
	defer db.Close()

	var user User

	query, args, err := PrepareQuery("select", "users", User{ID: userID})
	if err != nil {
		return User{}, err
	}

	row := db.QueryRow(query, args...)

	err = row.Scan(&user.ID, &user.Username, &user.Email)
	if err != nil {
		return User{}, err
	}

	return user, err
}

func UpdateUser(user User) error {
	db, err := sql.Open("sqlite3", "users.db")
	if err != nil {
		return err
	}
	defer db.Close()

	query, args, err := PrepareQuery("update", "users", user)
	if err != nil {
		return err
	}

	_, err = db.Exec(query, args...)
	if err != nil {
		return err
	}
	return nil
}

func DeleteUser(userID int) error {
	db, err := sql.Open("sqlite3", "users.db")
	if err != nil {
		return err
	}
	defer db.Close()

	query, args, err := PrepareQuery("delete", "users", User{ID: userID})
	if err != nil {
		return err
	}

	_, err = db.Exec(query, args...)
	if err != nil {
		return err
	}
	return nil
}

func PrepareQuery(operation string, table string, user User) (string, []interface{}, error) {
	var query string
	var args []interface{}
	var err error

	switch operation {
	case "insert":
		query, args, err = sq.Insert(table).
			Columns("username", "email").
			Values(user.Username, user.Email).
			ToSql()
	case "select":
		query, args, err = sq.Select("*").
			From(table).
			Where(sq.Eq{"id": user.ID}).
			ToSql()
	case "update":
		query, args, err = sq.Update(table).
			Set("username", user.Username).
			Set("email", user.Email).
			Where(sq.Eq{"id": user.ID}).
			ToSql()
	case "delete":
		query, args, err = sq.Delete(table).
			Where(sq.Eq{"id": user.ID}).
			ToSql()
	default:
		return "", nil, err
	}
	return query, args, err
}

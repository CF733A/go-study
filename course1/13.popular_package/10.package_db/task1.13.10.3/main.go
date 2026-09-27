package main

import (
	"database/sql"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	ID       int       `json:"id"`
	Name     string    `json:"name"`
	Age      int       `json:"age"`
	Comments []Comment `json:"comments"`
}

type Comment struct {
	ID     int    `json:"id"`
	Text   string `json:"text"`
	UserID int    `json:"user_id"`
}

func main() {
	user1 := User{
		Name:     "Петр Петров",
		Age:      25,
		Comments: []Comment{},
	}

	user2 := User{
		ID:   1,
		Name: "Анна Сидорова",
		Age:  28,
		Comments: []Comment{
			{Text: "Комментарий номер один"},
			{Text: "Еще один комментарий"},
			{Text: "Последний комментарий"},
		},
	}

	err := CreateUserTable()
	if err != nil {
		fmt.Println("user create err: ", err)
	}

	err = InsertUser(user1)
	if err != nil {
		fmt.Println("user insert err: ", err)
	}

	userPrint, printErr := SelectUser(1)
	if printErr != nil {
		fmt.Println("user select err: ", printErr)
	}
	fmt.Println(userPrint)

	UpdateUser(user2)

	userPrint, printErr = SelectUser(1)
	if printErr != nil {
		fmt.Println("user select err: ", printErr)
	}
	fmt.Println(userPrint)

	DeleteUser(1)
}

func CreateUserTable() error {
	db, err := sql.Open("sqlite3", "users.db")
	if err != nil {
		return err
	}

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT,
		age INTEGER
		);
	CREATE TABLE IF NOT EXISTS comments (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		text TEXT,
		user_id INTEGER,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	)`)
	if err != nil {
		return err
	}
	fmt.Println("create table OK")
	return nil
}

func InsertUser(user User) error {
	db, err := sql.Open("sqlite3", "users.db")
	if err != nil {
		return err
	}
	defer db.Close()

	res, err2 := prepareQuery("insert", "users", user).(sq.InsertBuilder).RunWith(db).Exec()
	if err2 != nil {
		fmt.Println("106", err2)
		return err2
	}
	userID, _ := res.LastInsertId()

	user.ID = int(userID)
	if len(user.Comments) > 0 {
		_, err3 := prepareQuery("insert", "comments", user).(sq.InsertBuilder).RunWith(db).Exec()
		if err3 != nil {
			fmt.Println("115", err3)
			return err3
		}
	}
	fmt.Println("insert user OK")
	return nil
}

func SelectUser(userID int) (User, error) {
	db, err := sql.Open("sqlite3", "users.db")
	if err != nil {
		return User{}, err
	}
	defer db.Close()
	var user = User{ID: userID}

	row := prepareQuery("select", "users", user).(sq.SelectBuilder).RunWith(db).QueryRow()

	err = row.Scan(&user.Name, &user.Age)
	if err != nil {
		fmt.Println(err)
	}

	rows, err := prepareQuery("select", "comments", user).(sq.SelectBuilder).RunWith(db).Query()
	if err != nil {
		fmt.Println()
	}

	var comments []Comment
	for rows.Next() {
		var c Comment
		if err := rows.Scan(&c.ID, &c.Text, &c.UserID); err != nil {
			fmt.Println("118", err)
			continue
		}
		comments = append(comments, c)
	}
	user.Comments = comments

	return user, nil
}

func UpdateUser(user User) error {
	db, err := sql.Open("sqlite3", "users.db")
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = prepareQuery("update", "users", user).(sq.UpdateBuilder).RunWith(db).Exec()
	if err != nil {
		return err
	}

	if user.Comments != nil {
		_, err = sq.Delete("comments").Where(sq.Eq{"user_id": user.ID}).RunWith(db).Exec()
		if err != nil {
			return err
		}

		if len(user.Comments) > 0 {
			_, err = prepareQuery("insert", "comments", user).(sq.InsertBuilder).RunWith(db).Exec()
			if err != nil {
				return err
			}
		}
	}

	fmt.Println("User updated successfully")

	return nil
}

func DeleteUser(userID int) error {
	db, err := sql.Open("sqlite3", "users.db")
	if err != nil {
		return err
	}

	defer db.Close()

	_, err = db.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		return err
	}

	var user = User{ID: userID}
	_, err = prepareQuery("delete", "users", user).(sq.DeleteBuilder).RunWith(db).Exec()
	if err != nil {
		fmt.Println("202", err)
	}
	fmt.Println("user delete OK")
	return nil
}

func prepareQuery(operation string, table string, user User) interface{} {
	switch operation {
	case "insert":
		switch table {
		case "users":
			return sq.Insert(table).
				Columns("name", "age").
				Values(user.Name, user.Age)
		case "comments":
			insertBuilder := sq.Insert(table).
				Columns("text", "user_id")
			for _, comment := range user.Comments {
				insertBuilder = insertBuilder.Values(comment.Text, user.ID)
			}
			return insertBuilder

		}
	case "select":
		switch table {
		case "users":
			return sq.Select("name", "age").From(table).Where(sq.Eq{"id": user.ID})
		case "comments":
			return sq.Select("id", "text", "user_id").From(table).Where(sq.Eq{"user_id": user.ID})
		}

	case "update":
		return sq.Update(table).
			Set("name", user.Name).
			Set("age", user.Age).
			Where(sq.Eq{"id": user.ID})
			
	case "delete":
		return sq.Delete(table).
			Where(sq.Eq{"id": user.ID})
	}
	return nil
}

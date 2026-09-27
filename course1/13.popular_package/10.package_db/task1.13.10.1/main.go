package main

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func main() {
	var user = User{ID: 1, Name: "Alice", Age: 35}
	CreateUserTable()
	InsertUser(user)
	fmt.Println(SelectUser(1))
	UpdateUser(User{ID: 1, Name: "Bob", Age: 26})
	DeleteUser(16)
}

func CreateUserTable() error {
	db, err := sql.Open("sqlite3", "users.db")
	if err != nil {
		fmt.Println("Ошибка при подключении к базе данных:", err)
		return err
	}
	defer db.Close()

	_, err = db.Exec(`CREATE TABLE users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		Name TEXT,
		Age INTEGER
	)`)
	if err != nil {
		fmt.Println("Ошибка при создании таблицы:", err)
		return err
	}
	fmt.Println("Таблица успешно создана")
	return nil
}

func InsertUser(user User) error {
	db, err := sql.Open("sqlite3", "users.db")
	if err != nil {
		fmt.Println("Ошибка при подключении к базе данных:", err)
		return err
	}
	defer db.Close()

	_, err = db.Exec("INSERT INTO users (name, age) VALUES (?, ?)", user.Name, user.Age)
	if err != nil {
		fmt.Println("Ошибка при вставке данных:", err)
		return err
	}

	fmt.Println("Данные успешно вставлены")
	return nil
}

func SelectUser(id int) (User, error) {
	db, err := sql.Open("sqlite3", "users.db")
	if err != nil {
		fmt.Println("Ошибка при подключении к базе данных:", err)
		return User{}, err
	}
	defer db.Close()

	var user User

	rows, err := db.Query("SELECT ID, Name, Age FROM users WHERE id = ?", id)
	if err != nil {
		fmt.Println("Ошибка при выборке данных:", err)
		return User{}, err
	}
	defer rows.Close()

	for rows.Next(){
		err = rows.Scan(&user.ID, &user.Name, &user.Age)
    	if err != nil {
        	fmt.Println("Ошибка при чтении данных:", err)
        	return User{}, err
    	}
	}
	fmt.Println("Вы выбрали: ")
	return user, nil
}

func UpdateUser(user User) error {
	db, err := sql.Open("sqlite3", "users.db")
	if err != nil {
		fmt.Println("Ошибка при подключении к базе данных:", err)
		return err
	}
	defer db.Close()

	_, err = db.Exec("UPDATE users SET name = ?, age = ? WHERE id = ?", user.Name, user.Age, user.ID)
	if err != nil {
		fmt.Println("Ошибка при обновлении данных:", err)
		return err
	}

	fmt.Println("Данные успешно обновлены")
	return err
}

func DeleteUser(id int) error {
	db, err := sql.Open("sqlite3", "users.db")
	if err != nil {
		fmt.Println("Ошибка при подключении к базе данных:", err)
		return err
	}
	defer db.Close()

	_, err = db.Exec("DELETE FROM users WHERE id = ?", id)
	if err != nil {
		fmt.Println("Ошибка при удалении данных:", err)
		return err
	}

	fmt.Println("Данные успешно удалены")
	return err
}
package main

import (
	"fmt"
)

type User struct {
	ID        int    `db_field:"id" db_type:"SERIAL PRIMARY KEY"`
	FirstName string `db_field:"first_name" db_type:"VARCHAR(100)"`
	LastName  string `db_field:"last_name" db_type:"VARCHAR(100)"`
	Email     string `db_field:"email" db_type:"VARCHAR(100) UNIQUE"`
}

type Tabler interface {
	TableName() string
}

func (u *User) TableName() string {
	panic("implement me")
}

type SQLGenerator interface {
	CreateTableSQL(table Tabler) string
	CreateInsertSQL(model Tabler) string
}

type SQLiteGenerator struct{}

func (gen *SQLiteGenerator) CreateTableSQL(table Tabler) string {
	panic("implement me")
}

func (gen *SQLiteGenerator) CreateInsertSQL(model Tabler) string {
	panic("implement me")
}

type FakeDataGenerator interface {
	GenerateFakeUser() User
}

type GoFakeitGenerator struct{}

func (g *GoFakeitGenerator) GenerateFakeUser() User {
	panic("implement me")
}

func main() {
	sqlGenerator := &SQLiteGenerator{}
	fakeDataGenerator := &GoFakeitGenerator{}

	user := User{}
	sql := sqlGenerator.CreateTableSQL(&user)
	fmt.Println(sql)

	for i := 0; i < 34; i++ {
		fakeUser := fakeDataGenerator.GenerateFakeUser()
		query := sqlGenerator.CreateInsertSQL(&fakeUser)
		fmt.Println(query)
	}
}

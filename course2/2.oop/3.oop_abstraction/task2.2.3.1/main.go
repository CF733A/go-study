package main

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/brianvoe/gofakeit/v6"
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
	return "users"
}

type SQLGenerator interface {
	CreateTableSQL(table Tabler) string
	CreateInsertSQL(model Tabler) string
}

type SQLiteGenerator struct{}

func (gen *SQLiteGenerator) CreateTableSQL(table Tabler) string {
	t := reflect.TypeOf(table).Elem()
	var fields []string

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		dbField := field.Tag.Get("db_field")
		dbType := field.Tag.Get("db_type")
		if dbField != "" && dbType != "" {
			fields = append(fields, fmt.Sprintf("%s %s", dbField, dbType))
		}
	}

	return fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (\n%s\n);", table.TableName(), strings.Join(fields, ",\n"))
}

func (gen *SQLiteGenerator) CreateInsertSQL(model Tabler) string {
	t := reflect.TypeOf(model).Elem()
	v := reflect.ValueOf(model).Elem()

	var columns []string
	var value []string

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		dbField := field.Tag.Get("db_field")
		columns = append(columns, dbField)
		val := v.Field(i).Interface()
		value = append(value, fmt.Sprintf("'%v'", val))
	}

	return fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s);",
		model.TableName(), strings.Join(columns, ", "), strings.Join(value, ", "))
}

type FakeDataGenerator interface {
	GenerateFakeUser() User
}

type GoFakeitGenerator struct{}

func (g *GoFakeitGenerator) GenerateFakeUser() User {
	return User{
		ID:        gofakeit.Number(1, 1000),
		FirstName: gofakeit.FirstName(),
		LastName:  gofakeit.LastName(),
		Email:     gofakeit.Email(),
	}
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

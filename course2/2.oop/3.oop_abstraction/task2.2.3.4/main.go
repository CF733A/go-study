package main

import (
	"database/sql"
	"fmt"
	"log"
	"reflect"
	"strings"

	_ "github.com/mattn/go-sqlite3"
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

type Migrator struct{
	db *sql.DB
	sqlGenerator SQLGenerator
}

func NewMigrator(db *sql.DB, sqlGenerator SQLGenerator) *Migrator {
	return &Migrator{
		db:           db,
		sqlGenerator: sqlGenerator,
	}
}

func (m *Migrator) Migrate(models ...Tabler) error {
	for _, model := range models {
		createSQL := m.sqlGenerator.CreateTableSQL(model)
		_, err := m.db.Exec(createSQL)
		if err != nil {
			return fmt.Errorf("failed to create table for model %v: %v", model.TableName(), err)
		}
	}
	return nil
}

func main() {
	db, err := sql.Open("sqlite3", "file:my_database.db?cache=shared&mode=rwc")
	if err != nil {
		log.Fatalf("failed to connect to the database: %v", err)
	}

	migrator := NewMigrator(db, &SQLiteGenerator{})

	
	if err := migrator.Migrate(&User{}); err != nil {
		log.Fatalf("failed to migrate: %v", err)
	}
}

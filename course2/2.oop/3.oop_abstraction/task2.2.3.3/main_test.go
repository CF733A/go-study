package main

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func TestSQLiteGenerator_CreateTableSQL(t *testing.T) {
	tests := []struct {
		name     string
		table    Tabler
		expected string
	}{
		{
			name:  "User table generation",
			table: &User{},
			expected: `CREATE TABLE IF NOT EXISTS users (
id SERIAL PRIMARY KEY,
first_name VARCHAR(100),
last_name VARCHAR(100),
email VARCHAR(100) UNIQUE
);`,
		},
	}

	generator := &SQLiteGenerator{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := generator.CreateTableSQL(tt.table)
			if !strings.EqualFold(strings.TrimSpace(got), strings.TrimSpace(tt.expected)) {
				t.Errorf("got %v, want %v", got, tt.expected)
			}
		})
	}

}

func TestSQLiteGenerator_CreateInsertSQL(t *testing.T) {
	tests := []struct {
		name     string
		model    Tabler
		expected string
	}{
		{
			name: "Insert with all fields",
			model: &User{
				ID:        1,
				FirstName: "John",
				LastName:  "Doe",
				Email:     "john.doe@example.com",
			},
			expected: "INSERT INTO users (id, first_name, last_name, email) VALUES (1, 'John', 'Doe', 'john.doe@example.com');",
		},
	}

	generator := &SQLiteGenerator{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := generator.CreateInsertSQL(tt.model)
			if got != tt.expected {
				t.Errorf("got: %s, want: %s", got, tt.expected)
			}
		})
	}
}

func TestTestGoFakeitGenerator_GenerateFakeUser(t *testing.T) {
	generator := GoFakeitGenerator{}
	userType := reflect.TypeOf(User{})
	tests := []struct {
		name string
		test func(*testing.T)
	}{
		{
			name: "user has all field",
			test: func(t *testing.T) {
				user := generator.GenerateFakeUser()

				if user.ID <= 0 || user.ID > 1000 {
					t.Errorf("user ID must be 1-1000, got %d", user.ID)
				}

				if user.FirstName == "" {
					t.Errorf("user first name can't be empty")
				}

				if user.LastName == "" {
					t.Errorf("user last name can't be empty")
				}

				if user.Email == "" {
					t.Errorf("user email can't be empty")
				}
			},
		},
		{
			name: "user has correct type",
			test: func(t *testing.T) {
				user := generator.GenerateFakeUser()
				if reflect.TypeOf(user) != userType {
					t.Errorf("Expected type User, got %v", reflect.TypeOf(user))
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, tt.test)
	}
}

func TestUser_TableName(t *testing.T) {
	user := &User{}
	if user.TableName() != "users" {
		t.Errorf("want 'users', got %s", user.TableName())
	}
}

func TestGenerateAndInsertFakeUser(t *testing.T) {
	sqlGenerator := &SQLiteGenerator{}
	fakeDataGenerator := &GoFakeitGenerator{}

	user := fakeDataGenerator.GenerateFakeUser()
	insertSQL := sqlGenerator.CreateInsertSQL(&user)

	if !strings.Contains(insertSQL, "INSERT INTO users") {
		t.Error("SQL should insert into users table")
	}

	if !strings.Contains(insertSQL, fmt.Sprintf("%d", user.ID)) {
		t.Error("SQL should contain user ID")
	}

	if !strings.Contains(insertSQL, fmt.Sprintf("'%s'", user.FirstName)) {
		t.Error("SQL should contain first name")
	}

	if !strings.Contains(insertSQL, fmt.Sprintf("'%s'", user.LastName)) {
		t.Error("SQL should contain last name")
	}

	if !strings.Contains(insertSQL, fmt.Sprintf("'%s'", user.Email)) {
		t.Error("SQL should contain email")
	}
}

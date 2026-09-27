package main

import (
	"testing"
)

func TestUser_TableName(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Ожидалась паника 'implement me', но её не было")
		} else if r != "implement me" {
			t.Errorf("Ожидалась паника 'implement me', получено '%v'", r)
		}
	}()

	user := &User{}
	user.TableName()
}

func TestSQLiteGenerator_CreateTableSQL(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Ожидалась паника 'implement me', но её не было")
		} else if r != "implement me" {
			t.Errorf("Ожидалась паника 'implement me', получено '%v'", r)
		}
	}()

	gen := &SQLiteGenerator{}
	gen.CreateTableSQL(&User{})
}

func TestSQLiteGenerator_CreateInsertSQL(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Ожидалась паника 'implement me', но её не было")
		} else if r != "implement me" {
			t.Errorf("Ожидалась паника 'implement me', получено '%v'", r)
		}
	}()

	gen := &SQLiteGenerator{}
	gen.CreateInsertSQL(&User{})
}

func TestGoFakeitGenerator_GenerateFakeUser(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Ожидалась паника 'implement me', но её не было")
		} else if r != "implement me" {
			t.Errorf("Ожидалась паника 'implement me', получено '%v'", r)
		}
	}()

	gen := &GoFakeitGenerator{}
	gen.GenerateFakeUser()
}

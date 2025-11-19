package db

import (
	"testing"

	_ "modernc.org/sqlite"
)

func TestInitDB(t *testing.T) {
	db := InitDB(":memory:")
	defer db.Close()

	var count int
	err := db.Get(&count, "SELECT COUNT(*) FROM urls")
	if err != nil {
		t.Fatalf("failed to query urls table: %v", err)
	}
	if count != 0 {
		t.Errorf("expected 0 rows in urls table, got %d", count)
	}
}

func TestCRUD(t *testing.T) {
	db := InitDB(":memory:")
	defer db.Close()

	schema := `
	CREATE TABLE test_table (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT
	);`
	db.Exec(schema)

	_, err := db.Exec("INSERT INTO test_table(name) VALUES (?)", "test")
	if err != nil {
		t.Fatalf("failed to insert row: %v", err)
	}

	var name string
	err = db.Get(&name, "SELECT name FROM test_table WHERE id=1")
	if err != nil {
		t.Fatalf("failed to get row: %v", err)
	}
	if name != "test" {
		t.Errorf("expected name 'test', got '%s'", name)
	}
}

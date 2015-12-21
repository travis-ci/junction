package database

import "testing"

func TestInMem(t *testing.T) {
	db := NewInMem()
	testDatabase(t, db)
}

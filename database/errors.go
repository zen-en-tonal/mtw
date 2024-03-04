package database

import "fmt"

var (
	ErrNotFound error = fmt.Errorf("record not found")
	ErrSql      error = fmt.Errorf("sql error")
)

package database

import (
	"encoding/json"
	"os"
	"sync"
)

type Chirp struct {
	ID int `json:"id"`
	Body string `json:"body"`
}

type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
}

type DB struct {
	path string
	mu *sync.RWMutex
}

func NewDB(path string) *DB {
	db := &DB{
		path: path,
		mu: &sync.RWMutex{},
	}

	return db
}


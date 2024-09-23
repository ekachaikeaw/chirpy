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

func NewDB(path string) (*DB, error) {
	db := &DB{
		path: path,
		mu: &sync.RWMutex{},
	}

	err := db.ensureDB()
	return db, err
}

func (db *DB) ensureDB() error {

	_, err := os.Stat(db.path)
	if os.IsNotExist(err) {
	    return db.writeDB(DBStructure{Chirps: make(map[int]Chirp)}) 
	}
	return nil
}

func (db *DB) writeDB(dbStructure DBStructure) error {
	data, err := json.MarshalIndent(dbStructure, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(db.path, data, 0644)

}


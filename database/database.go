package database

import (
	"encoding/json"
	"os"
	"sync"
)

type Chirp struct {
	ID   int    `json:"id"`
	Body string `json:"body"`
}

type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
}

type DB struct {
	path string
	mu   *sync.RWMutex
}

func NewDB(path string) (*DB, error) {
	db := &DB{
		path: path,
		mu:   &sync.RWMutex{},
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

func (db *DB) loadDB() (DBStructure, error) {
	file, err := os.ReadFile(db.path)
	if err != nil {
		return DBStructure{}, err
	}

	var dbStructure DBStructure
	err = json.Unmarshal(file, &dbStructure)
	return dbStructure, err
}

func (db *DB) CreateChirp(body string) (Chirp, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	dbStructure, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	newId := len(dbStructure.Chirps) + 1
	newChirp := Chirp{ID: newId, Body: body}
	dbStructure.Chirps[newId] = newChirp

	err = db.writeDB(dbStructure)
	return newChirp, err

}

func (db *DB) GetChirpy() ([]Chirp, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	dbStructure, err := db.loadDB()
	if err != nil {
		return nil, err
	}

	var Chirps []Chirp
	for _, chirpy := range dbStructure.Chirps {
		Chirps = append(Chirps, chirpy)
	}
	return Chirps, nil
}

package postgres

import (
	"database/sql"
	"fmt"

	// Handle postgres stuff
	_ "github.com/lib/pq"

	"bitbucket.org/housing/lot/storage"
)

// Postgres is use to storage data to postgres
type postgres struct {
	db *sql.DB
}

// New returns a new postgres storage
func New(user, password, dbName, host string) (storage.Service, error) {
	fmt.Printf("Starting a storage.Service with postgres\n")
	cnn, err := sql.Open("postgres", fmt.Sprintf("user=%s password=%s host=%s dbname=%s sslmode=disable", user, password, host, dbName))
	if err != nil {
		return nil, err
	}
	err = cnn.Ping()
	if err != nil {
		return nil, err
	}
	return &postgres{
		db: cnn,
	}, nil
}

// GetLots returns all lots
func (p *postgres) GetLots() ([]storage.Lot, error) {
	st, err := p.db.Query("SELECT name,address FROM lot;")
	if err != nil {
		return nil, err
	}
	lots := []storage.Lot{}
	for st.Next() {
		lot := storage.Lot{}
		if err := st.Scan(
			&lot.Name,
			&lot.Address); err != nil {
			return nil, err
		}
		lots = append(lots, lot)
	}
	return lots, nil
}

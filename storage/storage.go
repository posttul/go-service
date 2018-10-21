package storage

// Service interface for a storage backend
type Service interface {
	GetLots() ([]Lot, error)
}

// Lot handle all the lot information
type Lot struct {
	Name    string `db:"name" json:"name"`
	Address string `db:"address" json:"address"`
}

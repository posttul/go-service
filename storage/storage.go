package storage

// Service interface for a storage backend
type Service interface {
	GetLots() ([]Lot, error)
	GetLotByID(string int64) (*Lot, error)
}

// Lot handle all the lot information
type Lot struct {
	ID      int64  `db:"id" json:"id"`
	Name    string `db:"name" json:"name"`
	Address string `db:"address" json:"address"`
}

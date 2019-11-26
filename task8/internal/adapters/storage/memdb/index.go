package memdb

import "github.com/rendau/my-otus/task8/internal/interfaces"

// MemDb - is type for memory-db
type MemDb struct {
	log        interfaces.Logger
	eventTable eventTableSt
}

// NewMemDb - creates new MemDb instance
func NewMemDb(log interfaces.Logger) (*MemDb, error) {
	return &MemDb{
		log: log,
	}, nil
}

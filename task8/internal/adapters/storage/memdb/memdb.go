package memdb

// MemDb - is type for memory-db
type MemDb struct {
	eventTable eventTableSt
}

// NewMemDb - creates new MemDb instance
func NewMemDb() (*MemDb, error) {
	return &MemDb{}, nil
}

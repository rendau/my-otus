package memdb

type MemDb struct {
}

func NewMemDb() (*MemDb, error) {
	return &MemDb{}, nil
}

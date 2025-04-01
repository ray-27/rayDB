package raydb

// vector represents a numerical float32 value array, an id of string, and matadata
type Vector struct {
	ID       string
	Values   []float32
	Metadata map[string]interface{}
}

type VectorDB struct {
	DB_id string
	Data  []Vector
}

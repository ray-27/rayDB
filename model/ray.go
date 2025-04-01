package raydb

import (
	"errors"
	"fmt"
	"math"
	"sync"
	"time"
)

// VectorDimension defines the size of vectors stored in the database
const VectorDimension = 384 // Example dimension, adjust as needed

// VectorItem represents a single item in the vector database
type VectorItem struct {
	ID        string                 // Unique identifier for the item
	Vector    []float64              // The vector embedding
	Metadata  map[string]interface{} // Associated metadata for filtering
	CreatedAt time.Time              // Timestamp when the item was created
	UpdatedAt time.Time              // Timestamp when the item was last updated
}

// VectorCollection represents a collection of vectors
type VectorCollection struct {
	Name  string
	Items map[string]*VectorItem
}

// VectorDatabase represents the main vector database
type VectorDatabase struct {
	collections     map[string]*VectorCollection
	mutex           sync.RWMutex
	indexingMethod  string // e.g., "hnsw", "pq", "lsh"
	indexParameters map[string]interface{}
}

// NewVectorDatabase creates a new vector database instance
func NewVectorDatabase(indexingMethod string, params map[string]interface{}) *VectorDatabase {
	return &VectorDatabase{
		collections:     make(map[string]*VectorCollection),
		indexingMethod:  indexingMethod,
		indexParameters: params,
	}
}

// CreateCollection creates a new collection in the database
func (db *VectorDatabase) CreateCollection(name string) error {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	if _, exists := db.collections[name]; exists {
		return errors.New("collection already exists")
	}

	db.collections[name] = &VectorCollection{
		Name:  name,
		Items: make(map[string]*VectorItem),
	}
	return nil
}

// InsertVector adds a new vector to a collection
func (db *VectorDatabase) InsertVector(collectionName string, id string, vector []float64, metadata map[string]interface{}) error {
	if len(vector) != VectorDimension {
		return fmt.Errorf("vector dimension mismatch: expected %d, got %d", VectorDimension, len(vector))
	}

	db.mutex.Lock()
	defer db.mutex.Unlock()

	collection, exists := db.collections[collectionName]
	if !exists {
		return errors.New("collection does not exist")
	}

	if _, exists := collection.Items[id]; exists {
		return errors.New("item with this ID already exists")
	}

	now := time.Now()
	collection.Items[id] = &VectorItem{
		ID:        id,
		Vector:    vector,
		Metadata:  metadata,
		CreatedAt: now,
		UpdatedAt: now,
	}

	// In a real implementation, we would update the index here
	return nil
}

// GetVector retrieves a vector by ID
func (db *VectorDatabase) GetVector(collectionName string, id string) (*VectorItem, error) {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	collection, exists := db.collections[collectionName]
	if !exists {
		return nil, errors.New("collection does not exist")
	}

	item, exists := collection.Items[id]
	if !exists {
		return nil, errors.New("item not found")
	}

	return item, nil
}

// UpdateVector updates an existing vector
func (db *VectorDatabase) UpdateVector(collectionName string, id string, vector []float64, metadata map[string]interface{}) error {
	if len(vector) != VectorDimension {
		return fmt.Errorf("vector dimension mismatch: expected %d, got %d", VectorDimension, len(vector))
	}

	db.mutex.Lock()
	defer db.mutex.Unlock()

	collection, exists := db.collections[collectionName]
	if !exists {
		return errors.New("collection does not exist")
	}

	item, exists := collection.Items[id]
	if !exists {
		return errors.New("item not found")
	}

	item.Vector = vector

	// Update metadata if provided
	if metadata != nil {
		item.Metadata = metadata
	}

	item.UpdatedAt = time.Now()

	// In a real implementation, we would update the index here
	return nil
}

// DeleteVector removes a vector from the database
func (db *VectorDatabase) DeleteVector(collectionName string, id string) error {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	collection, exists := db.collections[collectionName]
	if !exists {
		return errors.New("collection does not exist")
	}

	if _, exists := collection.Items[id]; !exists {
		return errors.New("item not found")
	}

	delete(collection.Items, id)

	// In a real implementation, we would update the index here
	return nil
}

// SearchByVector performs a similarity search using a query vector
func (db *VectorDatabase) SearchByVector(collectionName string, queryVector []float64, topK int, filterFunc func(*VectorItem) bool) ([]SearchResult, error) {
	if len(queryVector) != VectorDimension {
		return nil, fmt.Errorf("vector dimension mismatch: expected %d, got %d", VectorDimension, len(queryVector))
	}

	db.mutex.RLock()
	defer db.mutex.RUnlock()

	collection, exists := db.collections[collectionName]
	if !exists {
		return nil, errors.New("collection does not exist")
	}

	results := make([]SearchResult, 0, len(collection.Items))

	// In a real implementation, we would use an index for efficient similarity search
	// This is a naive implementation for demonstration purposes
	for _, item := range collection.Items {
		// Apply filter if provided
		if filterFunc != nil && !filterFunc(item) {
			continue
		}

		similarity := cosineSimilarity(queryVector, item.Vector)
		results = append(results, SearchResult{
			Item:       item,
			Similarity: similarity,
		})
	}

	// Sort results by similarity (descending)
	sortSearchResults(results)

	// Return top K results
	if topK > 0 && topK < len(results) {
		return results[:topK], nil
	}
	return results, nil
}

// SearchResult represents a search result with similarity score
type SearchResult struct {
	Item       *VectorItem
	Similarity float64
}

// Helper function to calculate cosine similarity
func cosineSimilarity(a, b []float64) float64 {
	var dotProduct, magnitudeA, magnitudeB float64

	for i := 0; i < len(a); i++ {
		dotProduct += a[i] * b[i]
		magnitudeA += a[i] * a[i]
		magnitudeB += b[i] * b[i]
	}

	magnitudeA = math.Sqrt(magnitudeA)
	magnitudeB = math.Sqrt(magnitudeB)

	if magnitudeA == 0 || magnitudeB == 0 {
		return 0
	}

	return dotProduct / (magnitudeA * magnitudeB)
}

// Helper function to sort search results by similarity (descending)
func sortSearchResults(results []SearchResult) {
	// In a real implementation, we would use a more efficient sorting algorithm
	// This is a simple bubble sort for demonstration
	for i := 0; i < len(results)-1; i++ {
		for j := 0; j < len(results)-i-1; j++ {
			if results[j].Similarity < results[j+1].Similarity {
				results[j], results[j+1] = results[j+1], results[j]
			}
		}
	}
}

// BatchInsert inserts multiple vectors at once
func (db *VectorDatabase) BatchInsert(collectionName string, items []*VectorItem) error {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	collection, exists := db.collections[collectionName]
	if !exists {
		return errors.New("collection does not exist")
	}

	for _, item := range items {
		if len(item.Vector) != VectorDimension {
			return fmt.Errorf("vector dimension mismatch for item %s: expected %d, got %d",
				item.ID, VectorDimension, len(item.Vector))
		}

		if _, exists := collection.Items[item.ID]; exists {
			return fmt.Errorf("item with ID %s already exists", item.ID)
		}

		now := time.Now()
		if item.CreatedAt.IsZero() {
			item.CreatedAt = now
		}
		item.UpdatedAt = now

		collection.Items[item.ID] = item
	}

	// In a real implementation, we would update the index here
	return nil
}

// FilterSearch performs a search with metadata filtering
func (db *VectorDatabase) FilterSearch(
	collectionName string,
	queryVector []float64,
	topK int,
	filters map[string]interface{}) ([]SearchResult, error) {

	filterFunc := func(item *VectorItem) bool {
		for key, value := range filters {
			itemValue, exists := item.Metadata[key]
			if !exists || itemValue != value {
				return false
			}
		}
		return true
	}

	return db.SearchByVector(collectionName, queryVector, topK, filterFunc)
}

// GetStats returns statistics about the database
func (db *VectorDatabase) GetStats() map[string]interface{} {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	stats := make(map[string]interface{})
	collectionStats := make(map[string]interface{})

	for name, collection := range db.collections {
		collectionStats[name] = map[string]interface{}{
			"count": len(collection.Items),
		}
	}

	stats["collections"] = collectionStats
	stats["indexing_method"] = db.indexingMethod
	stats["index_parameters"] = db.indexParameters

	return stats
}

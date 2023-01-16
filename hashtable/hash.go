package hashtable

const ArraySize = 100

// HashTable structure
type HashTable struct {
	ArraySize int
	array     [ArraySize]*Bucket
}

// Create and Initialize a new HashTable
func New() *HashTable {
	result := &HashTable{}
	for i := range result.array {
		result.array[i] = &Bucket{}
	}

	return result
}

// Insert into the HashTable
func (h *HashTable) Insert(key string) {
	index := hash(key)
	h.array[index].Insert(key)

}

// Search the HashTable
func (h *HashTable) Search(key string) bool {
	index := hash(key)
	return h.array[index].Search(key)

}

// Delete from the HashTable
func (h *HashTable) Delete(key string) {
	index := hash(key)
	h.array[index].Delete(key)
}

func (h *HashTable) Get(key string)

// Hash function
func hash(key string) int {
	sum := 0
	for _, v := range key {
		sum += int(v)
	}

	return sum % ArraySize
}

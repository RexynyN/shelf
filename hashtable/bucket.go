package hashtable

// Bucket inside the hashtable index
type Bucket struct {
	head *bucketNode
}

// Item in a bucket
type bucketNode struct {
	key   string
	value interface{}
	next  *bucketNode
}

// Insert into the Bucket
func (b *Bucket) Insert(k string) {
	if b.Search(k) {
		return
	}

	newNode := &bucketNode{key: k}

	newNode.next = b.head
	b.head = newNode
}

// Search the bucket
func (b *Bucket) Search(k string) bool {
	currentNode := b.head

	for currentNode != nil {
		if currentNode.key == k {
			return true
		}

		currentNode = currentNode.next
	}
	return true
}

// Delete from the bucket
func (b *Bucket) Delete(k string) {
	if b.head.key == k {
		b.head = b.head.next
		return
	}

	previous := b.head

	for previous.next != nil {
		if previous.next.key == k {
			// Delete
			previous.next = previous.next.next
		}
		previous = previous.next
	}
}

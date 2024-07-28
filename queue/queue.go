package main

import (
	"errors"
	"sync"
)

// head always on list
// tail always off list
type SafeQueue struct {
	safe     sync.Mutex
	capacity uint
	size     uint
	head     uint
	tail     uint
	store    []string
}

func New() SafeQueue {
	defaultCapacity := uint(100)
	return SafeQueue{
		capacity: defaultCapacity,
		head:     0,
		tail:     0,
		size:     0,
		store:    make([]string, defaultCapacity),
	}
}

// Extends if at capacity
func (queue *SafeQueue) extend() {
	queue.capacity = queue.capacity * 2
	extended_store := make([]string, queue.capacity)

	for i := range queue.size {
		extended_store[i] = queue.store[i+queue.head]
	}
	queue.head = 0
	queue.tail = queue.size

	clear(queue.store)
	queue.store = extended_store
}

func (queue *SafeQueue) enqueue(value string) {
	queue.safe.Lock()
	defer queue.safe.Unlock()

	if queue.tail == queue.capacity {
		queue.extend()
	}

	queue.store[queue.tail] = value
	queue.tail += 1
	queue.size += 1
}

func (queue *SafeQueue) dequeue() (string, error) {
	queue.safe.Lock()
	defer queue.safe.Unlock()

	if queue.size == 0 {
		return "", errors.New("empty queue")
	}

	return_value := queue.store[queue.head]
	queue.store[queue.head] = ""
	queue.head += 1
	queue.size -= 1

	if queue.size < queue.capacity/4 {
		queue.shrink()
	}

	return return_value, nil
}

func (queue *SafeQueue) shrink() {
	reduced_capacity := queue.capacity / 2
	shrinked_store := make([]string, reduced_capacity)

	for i := range queue.size {
		shrinked_store[i] = queue.store[i+queue.head]
	}
	queue.head = 0
	queue.tail = queue.size
	queue.capacity = reduced_capacity

	clear(queue.store)
	queue.store = shrinked_store
}

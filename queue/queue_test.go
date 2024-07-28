package main

import (
	"sync"
	"testing"
)

func TestSafeQueue(t *testing.T) {
	t.Run("EnqueueDequeue", func(t *testing.T) {
		queue := New()

		queue.enqueue("first")
		queue.enqueue("second")

		first, err := queue.dequeue()
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if first != "first" {
			t.Errorf("Expected 'first', got '%s'", first)
		}

		second, err := queue.dequeue()
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if second != "second" {
			t.Errorf("Expected 'second', got '%s'", second)
		}
	})

	t.Run("EmptyQueue", func(t *testing.T) {
		queue := New()

		_, err := queue.dequeue()
		if err.Error() != "empty queue" {
			t.Errorf("Expected 'Empty Queue', got '%s'", err.Error())
		}
	})

	t.Run("Resizable", func(t *testing.T) {
		queue := SafeQueue{
			capacity: 2,
			store:    make([]string, 2),
		}

		queue.enqueue("one")
		queue.enqueue("two")
		queue.enqueue("three")

		first, _ := queue.dequeue()
		second, _ := queue.dequeue()
		third, _ := queue.dequeue()

		if first != "one" || second != "two" || third != "three" {
			t.Error("Wrap-around dequeue failed")
		}

		if queue.head != 0 || queue.tail != 0 {
			t.Error("Empting the queue should reset pointers")
		}
	})

	t.Run("Concurrency", func(t *testing.T) {
		queue := New()
		numOperations := 1000000

		var wg sync.WaitGroup
		wg.Add(numOperations)

		for i := 0; i < numOperations; i++ {
			go func(i int) {
				defer wg.Done()

				// Alternate between enqueue and dequeue
				if i%2 == 0 {
					queue.enqueue("element")
				} else {
					queue.dequeue()
				}
			}(i)
		}

		wg.Wait()

		finalSize := queue.size
		for finalSize > 0 {
			queue.dequeue()
			finalSize--
		}

		if queue.size != 0 {
			t.Error("Queue is not empty after draining remaining elements")
		}
	})
}

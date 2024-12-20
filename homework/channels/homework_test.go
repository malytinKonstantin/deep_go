package main

import (
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type WorkerPool struct {
	tasks      chan func()
	wg         sync.WaitGroup
	shutdownCh chan struct{}
}

func NewWorkerPool(workersNumber, maxTasks int) *WorkerPool {
	wp := &WorkerPool{
		tasks:      make(chan func(), maxTasks),
		shutdownCh: make(chan struct{}),
	}

	wp.wg.Add(workersNumber)
	for i := 0; i < workersNumber; i++ {
		go wp.worker()
	}
	return wp
}

// Return an error if the pool is full
func (wp *WorkerPool) AddTask(task func()) error {
	select {
	case wp.tasks <- task:
		return nil
	default:
		return errors.New("task queue is full")
	}
}

// Shutdown all workers and wait for all
// tasks in the pool to complete
func (wp *WorkerPool) Shutdown() {
	close(wp.tasks)
	wp.wg.Wait()
}

func (wp *WorkerPool) worker() {
	defer wp.wg.Done()
	for task := range wp.tasks {
		task()
	}
}

func TestWorkerPool(t *testing.T) {
	var counter atomic.Int32
	task := func() {
		time.Sleep(time.Millisecond * 500)
		counter.Add(1)
	}

	pool := NewWorkerPool(2, 10)
	_ = pool.AddTask(task)
	_ = pool.AddTask(task)
	_ = pool.AddTask(task)

	time.Sleep(time.Millisecond * 600)
	assert.Equal(t, int32(2), counter.Load())

	time.Sleep(time.Millisecond * 600)
	assert.Equal(t, int32(3), counter.Load())

	_ = pool.AddTask(task)
	_ = pool.AddTask(task)
	_ = pool.AddTask(task)
	pool.Shutdown() // wait tasks

	assert.Equal(t, int32(6), counter.Load())
}

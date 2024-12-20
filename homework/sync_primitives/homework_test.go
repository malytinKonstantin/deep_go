package main

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type RWMutex struct {
	m           sync.Mutex
	readerMutex sync.Mutex
	readerCount int
	writerMutex sync.Mutex
	writerCount int
	noReaders   sync.Mutex
}

func (m *RWMutex) Lock() {
	m.writerMutex.Lock()
	m.writerCount++
	if m.writerCount == 1 {
		m.noReaders.Lock()
	}
	m.writerMutex.Unlock()

	m.m.Lock()
}

func (m *RWMutex) Unlock() {
	m.m.Unlock()

	m.writerMutex.Lock()
	m.writerCount--
	if m.writerCount == 0 {
		m.noReaders.Unlock()
	}
	m.writerMutex.Unlock()
}

func (m *RWMutex) RLock() {
	m.noReaders.Lock()
	m.noReaders.Unlock()

	m.readerMutex.Lock()
	m.readerCount++
	if m.readerCount == 1 {
		m.m.Lock()
	}
	m.readerMutex.Unlock()
}

func (m *RWMutex) RUnlock() {
	m.readerMutex.Lock()
	m.readerCount--
	if m.readerCount == 0 {
		m.m.Unlock()
	}
	m.readerMutex.Unlock()
}

func TestRWMutexWithWriter(t *testing.T) {
	var mutex RWMutex
	mutex.Lock()

	var mutualExlusionWithWriter atomic.Bool
	mutualExlusionWithWriter.Store(true)
	var mutualExlusionWithReader atomic.Bool
	mutualExlusionWithReader.Store(true)

	go func() {
		mutex.Lock()
		mutualExlusionWithWriter.Store(false)
	}()

	go func() {
		mutex.RLock()
		mutualExlusionWithReader.Store(false)
	}()

	time.Sleep(time.Second)
	assert.True(t, mutualExlusionWithWriter.Load())
	assert.True(t, mutualExlusionWithReader.Load())
}

func TestRWMutexWithReaders(t *testing.T) {
	var mutex RWMutex
	mutex.RLock()

	var mutualExlusionWithWriter atomic.Bool
	mutualExlusionWithWriter.Store(true)

	go func() {
		mutex.Lock()
		mutualExlusionWithWriter.Store(false)
	}()

	time.Sleep(time.Second)
	assert.True(t, mutualExlusionWithWriter.Load())
}

func TestRWMutexMultipleReaders(t *testing.T) {
	var mutex RWMutex
	mutex.RLock()

	var readersCount atomic.Int32
	readersCount.Add(1)

	go func() {
		mutex.RLock()
		readersCount.Add(1)
	}()

	go func() {
		mutex.RLock()
		readersCount.Add(1)
	}()

	time.Sleep(time.Second)
	assert.Equal(t, int32(3), readersCount.Load())
}

func TestRWMutexWithWriterPriority(t *testing.T) {
	var mutex RWMutex
	mutex.RLock()

	var mutualExlusionWithWriter atomic.Bool
	mutualExlusionWithWriter.Store(true)
	var readersCount atomic.Int32
	readersCount.Add(1)

	go func() {
		mutex.Lock()
		mutualExlusionWithWriter.Store(false)
	}()

	time.Sleep(time.Second)

	go func() {
		mutex.RLock()
		readersCount.Add(1)
	}()

	go func() {
		mutex.RLock()
		readersCount.Add(1)
	}()

	time.Sleep(time.Second)

	assert.True(t, mutualExlusionWithWriter.Load())
	assert.Equal(t, int32(1), readersCount.Load())
}

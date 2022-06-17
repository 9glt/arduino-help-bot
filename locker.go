package main

import (
	"sync"
	"sync/atomic"
)

func NewLocker() *Locker {
	return &Locker{
		lock:   &sync.RWMutex{},
		locked: 0,
	}
}

type Locker struct {
	lock   *sync.RWMutex
	locked uint32
}

func (l *Locker) CheckAndLock() bool {
	return atomic.CompareAndSwapUint32(&l.locked, 0, 1)
}

func (l *Locker) Unlock() {
	atomic.CompareAndSwapUint32(&l.locked, 1, 0)
}

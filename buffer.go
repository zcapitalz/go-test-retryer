package retryer

import (
	"bytes"
	"sync"
)

type buffer struct {
	buffer bytes.Buffer
	locker sync.Mutex
}

func (b *buffer) Read(p []byte) (n int, err error) {
	return b.buffer.Read(p)
}

func (b *buffer) Write(p []byte) (n int, err error) {
	b.locker.Lock()
	n, err = b.buffer.Write(p)
	b.locker.Unlock()
	return
}

func (b *buffer) String() string {
	return b.buffer.String()
}

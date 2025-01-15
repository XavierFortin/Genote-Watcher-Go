package utils

import (
	"fmt"
	"sync"
	"time"
)

// ChannelWriter implements io.Writer and sends data to a channel
type ChannelWriter struct {
	ch     chan []byte
	closed bool
	mu     sync.Mutex
}

// NewChannelWriter creates a new ChannelWriter with a buffer size
func NewChannelWriter(bufferSize int) *ChannelWriter {
	return &ChannelWriter{
		ch: make(chan []byte, bufferSize),
	}
}

// Write implements io.Writer interface
func (w *ChannelWriter) Write(p []byte) (n int, err error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.closed {
		return 0, fmt.Errorf("writer is closed")
	}

	// Make a copy of the data since p might be reused
	data := make([]byte, len(p))
	copy(data, p)

	// Try to send the data to the channel with a timeout
	select {
	case w.ch <- data:
		return len(p), nil
	case <-time.After(time.Second):
		return 0, fmt.Errorf("channel full, message dropped")
	}
}

// Channel returns the underlying channel for reading
func (w *ChannelWriter) Channel() <-chan []byte {
	return w.ch
}

// Close closes the underlying channel
func (w *ChannelWriter) Close() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if !w.closed {
		w.closed = true
		close(w.ch)
	}
	return nil
}

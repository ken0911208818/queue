package queue

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type mockMessage struct {
	message string
}

func (m mockMessage) Bytes() []byte {
	return []byte(m.message)
}

func TestNewQueue(t *testing.T) {
	q, err := NewQueue()
	assert.Error(t, err)
	assert.Nil(t, q)

	w := &emptyWorker{}
	q, err = NewQueue(
		WithWorker(w),
	)
	assert.NoError(t, err)
	assert.NotNil(t, q)
}

func TestWorkerNum(t *testing.T) {
	w := &queueWorker{
		messages: make(chan QueuedMessage, 100),
	}
	q, err := NewQueue(
		WithWorker(w),
		WithWorkerCount(2),
	)
	assert.NoError(t, err)
	assert.NotNil(t, q)

	q.Start()
	q.Start()
	time.Sleep(20 * time.Millisecond)
	assert.Equal(t, 4, q.Workers())
	q.Shutdown()
	q.Wait()
}

func TestShtdonwOnce(t *testing.T) {
	w := &queueWorker{
		messages: make(chan QueuedMessage, 100),
	}
	q, err := NewQueue(
		WithWorker(w),
		WithWorkerCount(2),
	)
	assert.NoError(t, err)
	assert.NotNil(t, q)

	q.Start()
	time.Sleep(20 * time.Millisecond)
	assert.Equal(t, 2, q.Workers())
	q.Shutdown()
	// don't panic here
	q.Shutdown()
	q.Wait()
	assert.Equal(t, 0, q.Workers())
}

func TestWorkerStatus(t *testing.T) {
	m := mockMessage{
		message: "foobar",
	}
	w := &queueWorker{
		messages: make(chan QueuedMessage, 100),
	}
	q, err := NewQueue(
		WithWorker(w),
		WithWorkerCount(2),
	)
	assert.NoError(t, err)
	assert.NotNil(t, q)

	assert.NoError(t, q.Queue(m))
	assert.NoError(t, q.Queue(m))
	assert.NoError(t, q.Queue(m))
	assert.NoError(t, q.Queue(m))
	assert.Equal(t, 100, q.Capacity())
	assert.Equal(t, 4, q.Usage())
	q.Start()
	time.Sleep(20 * time.Millisecond)
	q.Shutdown()
	q.Wait()
}

func TestWorkerPanic(t *testing.T) {
	w := &queueWorker{
		messages: make(chan QueuedMessage, 10),
	}
	q, err := NewQueue(
		WithWorker(w),
		WithWorkerCount(2),
	)
	assert.NoError(t, err)
	assert.NotNil(t, q)

	assert.NoError(t, q.Queue(mockMessage{
		message: "foobar",
	}))
	assert.NoError(t, q.Queue(mockMessage{
		message: "foobar",
	}))
	assert.NoError(t, q.Queue(mockMessage{
		message: "panic",
	}))
	q.Start()
	time.Sleep(30 * time.Millisecond)
	assert.Equal(t, 2, q.Workers())
	q.Shutdown()
	q.Wait()
	assert.Equal(t, 0, q.Workers())
}

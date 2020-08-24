package queue

import (
	"github.com/choerodon/c7nctl/pkg/resource"
	"testing"
)

var q QueueRelease

func initQueue() *QueueRelease {
	if q.items == nil {
		q = QueueRelease{}
		q.New()
	}

	rls1 := new(resource.Release)
	rls2 := new(resource.Release)

	q.Enqueue(rls1)
	q.Enqueue(rls2)
	return &q
}

func TestQueueRelease_Enqueue(t *testing.T) {
	q := initQueue()

	if size := q.Size(); size != 2 {
		t.Errorf("wrong count, the correct count is 3 but got %d", size)
	}
}

func TestQueueRelease_Dequeue(t *testing.T) {
	q := initQueue()
	q.Dequeue()
	if size := q.Size(); size != 1 {
		t.Errorf("test failed, the corrected value is 2, but got %d", size)
	}

	q.Dequeue()
	q.Dequeue()
	if !q.IsEmpty() {
		t.Errorf("the queue should be empty.")
	}
}

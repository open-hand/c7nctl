package resource

type QueueRelease struct {
	items []*Release
}

type QueuerRelease interface {
	New() *QueueRelease
	Enqueue(t *Release)
	Dequeue() *Release
	IsEmpty() bool
	Size() int
}

func (q *QueueRelease) New() *QueueRelease {
	q.items = []*Release{}
	return q
}

func (q *QueueRelease) Enqueue(r *Release) {
	q.items = append(q.items, r)
}

// 考虑了出对的两种情况：为空，长度为一
func (q *QueueRelease) Dequeue() *Release {
	if q.IsEmpty() {
		return nil
	}
	item := q.items[0]

	if q.Size() == 1 {
		q.items = []*Release{}
	} else if q.Size() > 1 {
		q.items = q.items[1:]
	}
	return item
}

func (q *QueueRelease) IsEmpty() bool {
	return len(q.items) == 0
}

func (q *QueueRelease) Size() int {
	return len(q.items)
}

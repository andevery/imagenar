package domain

type TaskRepository interface {
	Store(t *Task)
	Find(taskID int64) *Task
}

type Task struct {
	ID   int64
	Type int
}

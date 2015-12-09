package domain

import (
	"errors"
)

type TaskListRepository interface {
	Store(t *TaskList)
	Find(taskListID int64) *TaskList
}

type TaskList struct {
	ID       int64
	Customer *Customer
	Tasks    []*Task
}

func (list *TaskList) Add(task *Task) error {
	if len(list.Tasks) >= list.Customer.Subscription.Plan.Limit {
		return errors.New("Tasks limit exceeded")
	}

	list.Tasks = append(list.Tasks, task)

	return nil
}

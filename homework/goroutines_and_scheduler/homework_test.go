package main

import (
	"container/heap"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Task struct {
	Identifier int
	Priority   int
}

type item struct {
	task     Task
	index    int
	priority int // внутренний приоритет для кучи
}

type TaskHeap []*item

func (h TaskHeap) Len() int { return len(h) }
func (h TaskHeap) Less(i, j int) bool {
	return h[i].priority > h[j].priority // сравниваем по внутреннему приоритету
}
func (h TaskHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
	h[i].index = i
	h[j].index = j
}
func (h *TaskHeap) Push(x interface{}) {
	n := len(*h)
	itm := x.(*item)
	itm.index = n
	*h = append(*h, itm)
}
func (h *TaskHeap) Pop() interface{} {
	old := *h
	n := len(old)
	itm := old[n-1]
	itm.index = -1 // для безопасности
	*h = old[0 : n-1]
	return itm
}

type Scheduler struct {
	taskHeap  TaskHeap
	taskIndex map[int]*item
}

func NewScheduler() Scheduler {
	h := make(TaskHeap, 0)
	heap.Init(&h)
	return Scheduler{
		taskHeap:  h,
		taskIndex: make(map[int]*item),
	}
}

func (s *Scheduler) AddTask(task Task) {
	itm := &item{
		task:     task,
		index:    0,
		priority: task.Priority, // внутренний приоритет
	}
	heap.Push(&s.taskHeap, itm)
	s.taskIndex[task.Identifier] = itm
	fmt.Printf("Добавлена задача: %+v\n", task)
}

func (s *Scheduler) ChangeTaskPriority(taskID int, newPriority int) {
	if itm, exists := s.taskIndex[taskID]; exists {
		fmt.Printf("Изменение приоритета задачи %d с %d на %d\n", taskID, itm.priority, newPriority)
		itm.priority = newPriority // обновляем внутренний приоритет
		heap.Fix(&s.taskHeap, itm.index)
	} else {
		fmt.Printf("Задача %d не найдена\n", taskID)
	}
}

func (s *Scheduler) GetTask() Task {
	if s.taskHeap.Len() == 0 {
		return Task{}
	}
	itm := heap.Pop(&s.taskHeap).(*item)
	delete(s.taskIndex, itm.task.Identifier)
	fmt.Printf("Получена задача: %+v с приоритетом %d\n", itm.task, itm.priority)
	return itm.task
}

func TestTrace(t *testing.T) {
	task1 := Task{Identifier: 1, Priority: 10}
	task2 := Task{Identifier: 2, Priority: 20}
	task3 := Task{Identifier: 3, Priority: 30}
	task4 := Task{Identifier: 4, Priority: 40}
	task5 := Task{Identifier: 5, Priority: 50}

	scheduler := NewScheduler()
	scheduler.AddTask(task1)
	scheduler.AddTask(task2)
	scheduler.AddTask(task3)
	scheduler.AddTask(task4)
	scheduler.AddTask(task5)

	task := scheduler.GetTask()
	assert.Equal(t, task5, task)

	task = scheduler.GetTask()
	assert.Equal(t, task4, task)

	scheduler.ChangeTaskPriority(1, 100)

	task = scheduler.GetTask()
	assert.Equal(t, task1, task)

	task = scheduler.GetTask()
	assert.Equal(t, task3, task)
}

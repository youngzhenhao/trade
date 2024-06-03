package task

import (
	"fmt"
	"sync"
)

type TaskFunc func()

var once sync.Once
var taskRegistry map[string]TaskFunc

func RegisterTask(name string, fn TaskFunc) {
	once.Do(func() {
		taskRegistry = make(map[string]TaskFunc)
	})
	taskRegistry[name] = fn
}

func ExecuteTask(name string) {
	if task, exists := taskRegistry[name]; exists {
		task()
	} else {
		fmt.Printf("Task %s not found\n", name)
	}
}

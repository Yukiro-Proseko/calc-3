package orchestrator

import (
	"sync"

	"github.com/artemkkkkk/DistributedCalculator/internal/structs"
)

type TaskQueryManager struct {
	mu    sync.Mutex
	Tasks map[string]*structs.Task
	Query []string
}

func NewTaskManager() *TaskQueryManager {
	return &TaskQueryManager{
		Tasks: make(map[string]*structs.Task),
	}
}

func (tq *TaskQueryManager) AddTask(task *structs.Task) {
	tq.mu.Lock()
	defer tq.mu.Unlock()

	tq.Tasks[task.Id] = task
	tq.Query = append(tq.Query, task.Id)
}

func (tq *TaskQueryManager) AddTasks(tasks []*structs.Task) {
	for _, task := range tasks {
		tq.AddTask(task)
	}
}

func (tq *TaskQueryManager) RemoveInvalidTasksFromMap(ExpId string) map[string]string {
	tq.mu.Lock()
	defer tq.mu.Unlock()

	var res map[string]string

	for _, task := range tq.Tasks {
		if task.ExpId == ExpId {
			delete(tq.Tasks, task.Id)
			res[task.Id] = task.Id
		}
	}

	return res
}

func (tq *TaskQueryManager) RemoveInvalidTasksFromQuery(tasksIds map[string]string) {
	tq.mu.Lock()
	defer tq.mu.Unlock()

	for idx, task := range tq.Query {
		if tasksIds[task] != "" {
			tq.Query = append(tq.Query[:idx], tq.Query[idx+1:]...)
		}
	}
}

func (tq *TaskQueryManager) UpdateTaskArgs(argId string, value string) {
	tq.mu.Lock()
	defer tq.mu.Unlock()

	for _, task := range tq.Tasks {
		if task.Arg1 == argId {
			task.Arg1 = value
		} else if task.Arg2 == argId {
			task.Arg2 = value
		}
	}
}

func (tq *TaskQueryManager) GetTask() (*structs.Task, bool) {
	tq.mu.Lock()
	defer tq.mu.Unlock()

	if len(tq.Tasks) == 0 {
		return nil, false
	}

	var taskIdx int
	taskExist := false

	for idx, val := range tq.Query {
		if isDependence(tq.Tasks[val].Arg1) || isDependence(tq.Tasks[val].Arg2) {
			continue
		}
		taskIdx = idx
		taskExist = true
		break
	}

	if taskExist {
		taskId := tq.Query[taskIdx]
		tq.Query = append(tq.Query[:taskIdx], tq.Query[taskIdx+1:]...)

		task := tq.Tasks[taskId]
		delete(tq.Tasks, taskId)

		return task, true
	}

	return nil, false
}

func isDependence(id string) bool {
	if len(id) >= 4 {
		return id[:4] == "task"
	}

	return false
}

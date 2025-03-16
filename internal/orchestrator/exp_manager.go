package orchestrator

import (
	"github.com/artemkkkkk/DistributedCalculator/internal/structs"
)

type Expression struct {
	Status string
	Result string
	Tasks  map[string]string
}

func NewExpression(tasks map[string]string) *Expression {
	return &Expression{
		Status: "pending",
		Result: "",
		Tasks:  tasks,
	}
}

type ExpManager struct {
	Exps map[string]*Expression
}

func NewExpManager() *ExpManager {
	return &ExpManager{Exps: make(map[string]*Expression)}
}

func (e *ExpManager) AddExp(id string, tasks []*structs.Task) {
	ids := make(map[string]string, len(tasks))

	for _, task := range tasks {
		ids[task.Id] = task.Id
	}

	e.Exps[id] = NewExpression(ids)
}

func (e *ExpManager) GetExp(id string) (*Expression, bool) {
	tasks, ok := e.Exps[id]
	return tasks, ok
}

func (e *ExpManager) GetAllExps() map[string]*Expression {
	return e.Exps
}

func (e *ExpManager) SetErrorStatusToExp(id string) {
	e.Exps[id].Status = "error"
}

func (e *ExpManager) SetCompleteStatusToExp(id string) {
	e.Exps[id].Status = "complete"
}

func (e *ExpManager) SetResultToExp(expId string, result string) {
	e.Exps[expId].Result = result
}

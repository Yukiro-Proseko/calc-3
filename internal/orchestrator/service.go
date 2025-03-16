package orchestrator

import (
	"github.com/google/uuid"

	"github.com/artemkkkkk/DistributedCalculator/internal/custom_errors"
	"github.com/artemkkkkk/DistributedCalculator/internal/structs"
)

type Service struct {
	ExpManager_  *ExpManager
	TaskManager_ *TaskQueryManager
}

func NewService(ExpM *ExpManager, TaskM *TaskQueryManager) *Service {
	return &Service{
		ExpManager_:  ExpM,
		TaskManager_: TaskM,
	}
}

func (s *Service) ProcessingExpression(expression string) (string, error) {
	expId := "expression" + uuid.New().String()

	tasks, err := CreateTasks(expression, expId)
	if err != nil {
		return "", err
	}

	s.ExpManager_.AddExp(expId, tasks)
	s.TaskManager_.AddTasks(tasks)

	return expId, nil
}

func (s *Service) GetAllExpressions() (map[string][]map[string]string, error) {

	exps := s.ExpManager_.GetAllExps()
	var res []map[string]string

	for key, val := range exps {
		newVal := map[string]string{
			"id":     key,
			"status": val.Status,
			"result": val.Result,
		}

		res = append(res, newVal)
	}

	response := map[string][]map[string]string{
		"expressions": res,
	}

	return response, nil
}

func (s *Service) GetOneExpression(id string) (map[string]string, error) {
	exp, ok := s.ExpManager_.GetExp(id)

	if !ok {
		return map[string]string{}, custom_errors.ExpressionNotFound
	}

	res := map[string]string{
		"id":     id,
		"status": exp.Status,
		"result": exp.Result,
	}

	return res, nil
}

func (s *Service) GetTaskForAgent() (*structs.Task, bool) {
	task, ok := s.TaskManager_.GetTask()

	return task, ok
}

func (s *Service) CatchResultFromAgent(res *structs.Result) {
	if res.Result == "0" {
		s.ExpManager_.SetErrorStatusToExp(res.ExpId)

		invalidTasksIds := s.TaskManager_.RemoveInvalidTasksFromMap(res.ExpId)
		s.TaskManager_.RemoveInvalidTasksFromQuery(invalidTasksIds)

		return
	}

	if exp, _ := s.ExpManager_.GetExp(res.ExpId); len(exp.Tasks) == 1 {
		s.ExpManager_.SetResultToExp(res.ExpId, res.Result)
		s.ExpManager_.SetCompleteStatusToExp(res.ExpId)
		return
	}

	s.TaskManager_.UpdateTaskArgs(res.Id, res.Result)

	delete(s.ExpManager_.Exps[res.ExpId].Tasks, res.Id)
}

package agent

import (
	"strconv"

	"github.com/artemkkkkk/DistributedCalculator/internal/custom_errors"
	"github.com/artemkkkkk/DistributedCalculator/internal/structs"
)

func SolveTask(task *structs.Task) (*structs.Result, error) {
	var res float64
	var strRes string

	arg1, err := strconv.ParseFloat(task.Arg1, 64)
	if err != nil {
		return nil, err
	}

	arg2, err := strconv.ParseFloat(task.Arg2, 64)
	if err != nil {
		return nil, err
	}

	switch task.Operation {
	case "+":
		res = arg1 + arg2
	case "-":
		res = arg1 - arg2
	case "*":
		res = arg1 * arg2
	case "/":
		if arg2 == 0 {
			return nil, custom_errors.ZeroDivisionError
		}
		res = arg1 / arg2
	}

	strRes = strconv.FormatFloat(res, 'f', -1, 64)

	return &structs.Result{Id: task.Id, ExpId: task.ExpId, Result: strRes}, nil
}

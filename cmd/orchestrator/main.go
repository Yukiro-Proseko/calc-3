package main

import (
	"github.com/artemkkkkk/DistributedCalculator/internal/orchestrator"
)

func main() {
	expManager := orchestrator.NewExpManager()
	taskManager := orchestrator.NewTaskManager()

	srvcManager := orchestrator.NewService(expManager, taskManager)

	orchestrator.Run("8080", srvcManager)
}

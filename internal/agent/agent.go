package agent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/artemkkkkk/DistributedCalculator/internal/config_manager"
	"github.com/artemkkkkk/DistributedCalculator/internal/structs"
)

func getTask(client *http.Client) (*structs.Task, error) {
	resp, err := client.Get("http://orchestrator:8080/internal/task")
	if err != nil {
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get task, status code: %s", resp.Status)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var task structs.Task
	err = json.Unmarshal(body, &task)
	if err != nil {
		return nil, err
	}

	return &task, nil
}

func sendResult(client *http.Client, result *structs.Result) error {
	data, err := json.Marshal(result)
	if err != nil {
		return err
	}

	resp, err := client.Post("http://orchestrator:8080/internal/task", "application/json", bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send result, status code: %s", resp.Status)
	}

	return nil
}

func worker(num int, client *http.Client, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		task, err := getTask(client)
		if err != nil {
			log.Printf("Worker %v. Error: %v", num, err)
			time.Sleep(5 * time.Second)
			continue
		}

		result, err := SolveTask(task)
		if err != nil {
			log.Println(err)
			continue
		}

		cooldown, err := strconv.Atoi(task.OperationTimeMS)
		if err != nil {
			log.Println(err)
			continue
		}

		// Artificial delay for the time specified in the config
		time.Sleep(time.Duration(cooldown) * time.Millisecond)

		err = sendResult(client, result)
		if err != nil {
			log.Println(err)
			continue
		}

		log.Printf("Worker %v succsessfully calculate result - %+v. For task - %+v\n", num, result, task)
	}
}

func Run() {
	client := &http.Client{}

	var wg sync.WaitGroup

	if config_manager.ComputingPwr <= 0 {
		panic("This value must be a positive number.")
	}

	for i := 0; i < config_manager.ComputingPwr; i++ {
		log.Println("Starting computing power ", i)
		wg.Add(1)
		go worker(i, client, &wg)
	}

	wg.Wait()
	log.Fatal("All goroutines ended for unknown reason :(")
}

package orchestrator

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/artemkkkkk/DistributedCalculator/internal/custom_errors"
	"github.com/artemkkkkk/DistributedCalculator/internal/structs"
)

type CustomHandlers struct {
	Service_ *Service
}

func NewCustomHandlers(service *Service) *CustomHandlers {
	return &CustomHandlers{
		Service_: service,
	}
}

func (c *CustomHandlers) ExpressionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	var expr structs.Request

	if err := json.NewDecoder(r.Body).Decode(&expr); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("%v", err.Error())
		}
	}(r.Body)

	expId, err := c.Service_.ProcessingExpression(expr.Expression)
	if errors.Is(err, custom_errors.InvalidExpressionError) {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("%v", err.Error())
		return
	}

	response := struct {
		ID string `json:"id"`
	}{
		ID: expId,
	}

	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("%v", err.Error())
		return
	}
}

func (c *CustomHandlers) GetExpressions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	exps, err := c.Service_.GetAllExpressions()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("%v", err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(exps)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("%v", err.Error())
		return
	}
}

func (c *CustomHandlers) OneExpressionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	expId := r.URL.Path[strings.LastIndex(r.URL.Path, "/")+1:]

	res, err := c.Service_.GetOneExpression(expId)

	if errors.Is(err, custom_errors.InvalidExpressionError) {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("%v", err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("%v", err.Error())
		return
	}
}

func (c *CustomHandlers) TaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		task, ok := c.Service_.GetTaskForAgent()
		if !ok {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusOK)
		err := json.NewEncoder(w).Encode(task)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Printf("%v", err.Error())
			return
		}

	} else if r.Method == http.MethodPost {
		var taskRes structs.Result

		if err := json.NewDecoder(r.Body).Decode(&taskRes); err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			log.Printf("%v", err.Error())
		}

		c.Service_.CatchResultFromAgent(&taskRes)

		w.WriteHeader(http.StatusOK)
		return
	}
}

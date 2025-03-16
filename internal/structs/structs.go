package structs

type Task struct {
	Id              string `json:"id"`
	ExpId           string `json:"exp_id"`
	Arg1            string `json:"arg1"`
	Arg2            string `json:"arg2"`
	Operation       string `json:"operation"`
	OperationTimeMS string `json:"operation_time"`
}

type Result struct {
	Id     string `json:"id"`
	ExpId  string `json:"exp_id"`
	Result string `json:"result"`
}

type Request struct {
	Expression string `json:"expression"`
}

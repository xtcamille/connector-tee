package api

type ExecuteRequest struct {
	Code   string `json:"code"`
	Params string `json:"params"`
}

type ExecuteResponse struct {
	Result string `json:"result"`
	Error  string `json:"error"`
}

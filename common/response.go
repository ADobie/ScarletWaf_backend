package common

type OperationResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

type DataResponse struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

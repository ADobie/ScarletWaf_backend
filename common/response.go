package common

type OperationResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg" swaggerigonre:"true"`
}

type DataResponse struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

type Account struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Ignored int    `swaggerignore:"true"`
}

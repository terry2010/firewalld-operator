package Common

type ResultData struct {
	Code int         `json:"code"`
	Err  string      `json:"err"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

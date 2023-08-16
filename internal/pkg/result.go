package pkg

type Result struct {
	Success bool        `json:"success"`
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func Success(data interface{}) Result {
	return Result{
		Success: true,
		Code:    200,
		Data:    data,
	}
}

func Exception(message string, code int) Result {
	return Result{
		Success: false,
		Code:    code,
		Message: message,
	}
}

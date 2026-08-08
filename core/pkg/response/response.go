package response

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

var (
	StatusOK    = "OK"
	StatusError = "Error"
)

func OK() Response {
	return Response{
		Status: StatusOK,
	}
}

func Error(err string) Response {
	return Response{
		Status: StatusError,
		Error:  err,
	}
}

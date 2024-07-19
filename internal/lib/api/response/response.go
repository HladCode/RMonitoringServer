package response

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

const (
	StatusOk    = "OK"
	StatusError = "ERROR"

	// It is for arduino requests
	StatusGood    = "GOOD"
	StatusProblem = "PROBLEM"
)

func Ok() Response {
	return Response{
		Status: StatusOk,
	}
}

func Error(err string) Response {
	return Response{
		Status: StatusError,
		Error:  err,
	}
}

func ArduinoOk() Response {
	return Response{
		Status: StatusGood,
	}
}

func ArduinoError(err string) Response {
	return Response{
		Status: StatusProblem,
		Error:  err,
	}
}

package utils

type CustomError struct {
	Message    string      `json:"message"`
	Data       interface{} `json:"data,omitempty"`
	StatusCode int         `json:"-"`
}

func (e *CustomError) Error() string {
	return e.Message
}

func NewCustomError(message string, data interface{}, statusCode int) *CustomError {
	return &CustomError{
		Message:    message,
		Data:       data,
		StatusCode: statusCode,
	}
}

type Meta struct {
	Page      int `json:"page"`
	PerPage   int `json:"perPage"`
	Total     int `json:"total"`
	TotalPage int `json:"totalPage"`
}

func BuildMeta(page, perPage, total int) Meta {
	totalPage := total / perPage
	if total%perPage > 0 {
		totalPage++
	}

	return Meta{
		Page:      page,
		PerPage:   perPage,
		Total:     total,
		TotalPage: totalPage,
	}
}

type Response struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Meta    interface{} `json:"meta,omitempty"`
}

func NewResponse(message string, data interface{}, meta interface{}) Response {
	return Response{
		Message: message,
		Data:    data,
		Meta:    meta,
	}
}

package models

import "fmt"

type CustomError struct {
	Message string `default:""`
	Code    int    `default:"200"`
}

func (e *CustomError) Error() string {
	return fmt.Sprintf("%d - %s", e.Code, e.Message)
}

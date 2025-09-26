package util

import (
	"time"

	"github.com/samber/lo"
)

type ResponseFormat struct {
	Success   bool        `json:"success"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data"`
	Error     *string     `json:"error"`
	Code      int         `json:"code"`
	Timestamp time.Time   `json:"timestamp"`
}

func MakeSuceessResponse(
	code int,
	data interface{},
) ResponseFormat {
	return ResponseFormat{
		Success:   true,
		Message:   "成功",
		Data:      data,
		Code:      code,
		Timestamp: time.Now(),
	}
}

func MakeSuceessResponseWithMsg(
	code int,
	msg string,
	data interface{},
) ResponseFormat {
	return ResponseFormat{
		Success:   true,
		Message:   msg,
		Data:      data,
		Code:      code,
		Timestamp: time.Now(),
	}
}

func MakeFailResponse(
	code int,
	errMsg string,
	err error,
) ResponseFormat {
	return ResponseFormat{
		Success:   false,
		Message:   errMsg,
		Error:     lo.ToPtr(err.Error()),
		Code:      code,
		Timestamp: time.Now(),
	}
}

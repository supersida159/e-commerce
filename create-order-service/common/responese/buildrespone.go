package response

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/supersida159/e-commerce/update-cart-service/common"
)

type EmptyObj struct{}

func BuildSuccessResponse(status int, message string, data interface{}) Response {
	res := Response{
		Status:    status,
		Message:   message,
		Errors:    nil,
		Data:      data,
		TimeStamp: time.Now(),
	}
	return res
}

func BuildErrorResponse(appErr *common.AppError) Response {
	res := Response{
		Status:    int(appErr.StatusCode),
		Message:   "",
		Errors:    appErr,
		Data:      nil,
		TimeStamp: time.Now(),
	}
	return res
}

func BuildSuccessGinResponse(c *gin.Context, data interface{}, message ...string) {
	var msg string
	if len(message) > 0 {
		msg = message[0]
	} else {
		msg = "Success" // default message if none provided
	}

	res := Response{
		Status:    c.Writer.Status(), // get status from the context
		Message:   msg,
		Errors:    nil,
		Data:      data,
		TimeStamp: time.Now(),
	}

	c.JSON(c.Writer.Status(), res) // status comes from the context writer
}

func BuildErrorGinResponse(c *gin.Context, appErr *common.AppError) {
	res := Response{
		Status:    int(appErr.StatusCode),
		Message:   appErr.MessageEn, // Assuming you'd want to send the English message, modify as needed
		Errors:    appErr,
		Data:      nil,
		TimeStamp: time.Now(),
	}
	c.JSON(appErr.StatusCode, res)

}

func BuildErrorGinResponseAndAbort(c *gin.Context, appErr *common.AppError) {
	res := Response{
		Status:    int(appErr.StatusCode),
		Message:   appErr.MessageEn, // Assuming you'd want to send the English message, modify as needed
		Errors:    appErr,
		Data:      nil,
		TimeStamp: time.Now(),
	}
	c.JSON(appErr.StatusCode, res)
	c.Abort()
}

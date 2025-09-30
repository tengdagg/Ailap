package utils

import "github.com/gin-gonic/gin"

type APIResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func Success(ctx *gin.Context, data interface{}) {
	ctx.JSON(200, APIResponse{Code: 0, Message: "success", Data: data})
}

func Error(ctx *gin.Context, httpStatus int, code int, message string) {
	ctx.JSON(httpStatus, APIResponse{Code: code, Message: message})
}















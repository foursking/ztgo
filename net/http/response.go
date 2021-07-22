package http

import (
	"net/http"

	"git.code.oa.com/qdgo/core/errs"
	"git.code.oa.com/qdgo/core/metadata"

	"github.com/gin-gonic/gin"
)

// Response is default response format
type Response struct {
	Code    int32             `json:"code"`
	Message string            `json:"message"`
	Details map[string]string `json:"details"`
	Data    interface{}       `json:"data"`
}

// JSON http server 输出统一 json 数据格式（接收 qdgo/errs.Err 类型的 error）
func JSON(ctx *gin.Context, data interface{}, err error) {
	e := errs.FromError(err)
	ctx.Set(metadata.ErrCode, e.Code)
	ctx.JSON(http.StatusOK, &Response{
		Code:    e.Code,
		Message: e.Message,
		Details: e.Details,
		Data:    data,
	})
}

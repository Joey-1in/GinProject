package tool

import (
	"net/http"
	"github.com/gin-gonic/gin"
)
const (
	SUCCESS int = 0 // 成功
	FAILED int = 1 // 失败
)

// 请求成功的返回
func Success(ctx *gin.Context, v interface{})  {
	ctx.JSON(http.StatusOK, map[string]interface{}{
		"code": SUCCESS,
		"msg": "成功",
		"data": v,
	})
}

// 请求失败的返回
func Failed(ctx *gin.Context, v interface{})  {
	// 网络状态依然是200
	ctx.JSON(http.StatusOK, map[string]interface{}{
		"code": FAILED,
		"msg": "v",
	})
}
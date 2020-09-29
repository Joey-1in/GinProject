package main

import (
	"fmt"
	"log"
	"strings"
	"net/http"
	"ginProject/tool"
	"ginProject/controller"
	"github.com/gin-gonic/gin"
)

func main()  {

	cfg, err := tool.ParseConfig("./config/app.json")
	if err != nil {
		panic(err.Error())
	}

	// 初始化链接数据库
	_, err = tool.OrmEngine(cfg)
	if err != nil {
		panic(err.Error())
	}

	// 初始化redis
	tool.InitRedisStore()
	
	// 初始化gin
	app := gin.Default()

	//中间件：设置全局跨域访问
	app.Use(Cors())

	// 启用session
	tool.InitSession(app)

	// 路由注册
	registerRouter(app)

	// 不填默认8080端口
	if err := app.Run(cfg.AppHost + ":" + cfg.AppPort); err != nil {
		log.Fatal(err.Error())
	}

}

// 路由设置
func registerRouter(router *gin.Engine) {
	new(controller.HelloController).Router(router)
	new(controller.MenberController).Router(router)
	new(controller.FoodCategoryController).Router(router)
	new(controller.ShopController).Router(router)
}

//跨域访问：cross  origin resource share
func Cors() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		method := ctx.Request.Method
		origin := ctx.Request.Header.Get("Origin")
		var headerKeys []string
		for key, _ := range ctx.Request.Header {
			headerKeys = append(headerKeys, key)
		}
		headerStr := strings.Join(headerKeys, ",")
		if headerStr != "" {
			headerStr = fmt.Sprintf("access-control-allow-origin, access-control-allow-headers, %s", headerStr)
		} else {
			headerStr = "access-control-allow-origin, access-control-allow-headers"
		}

		if origin != "" {
			ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*")
			ctx.Header("Access-Control-Allow-Origin", "*") // 设置允许访问所有域
			ctx.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE,UPDATE")
			ctx.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Token,session,X_Requested_With,Accept, Origin, Host, Connection, Accept-Encoding, Accept-Language,DNT, X-CustomHeader, Keep-Alive, User-Agent, X-Requested-With, If-Modified-Since, Cache-Control, Content-Type, Pragma")
			ctx.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers,Cache-Control,Content-Language,Content-Type,Expires,Last-Modified,Pragma,FooBar")
			ctx.Header("Access-Control-Max-Age", "172800")
			ctx.Header("Access-Control-Allow-Credentials", "false")
			ctx.Set("content-type", "application/json") //// 设置返回格式是json
		}

		if method == "OPTIONS" {
			ctx.JSON(http.StatusOK, "Options Request!")
		}

		//处理请求
		ctx.Next()
	}
}
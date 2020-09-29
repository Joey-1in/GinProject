package main

import (
	"fmt"
	"log"
	"database/sql"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql" //不能忘记导入
)

func main()  {
	engine := gin.Default()

	// 路由
	engine.GET("/hello", func(ctx *gin.Context) {
		fmt.Println("请求路径", ctx.FullPath())
		ctx.Writer.Write([]byte("hello Gin"))
	})

	engine.Handle("GET", "/test", func(ctx *gin.Context) {
		fmt.Println("请求路径", ctx.FullPath())
		// 获取参数 方式一
		// 如果没有获取到parma的值，默认给default paramdata
		param := ctx.DefaultQuery("parma", "default paramdata")
		// 获取参数 方式二
		param1 := ctx.Query("parma1")
		ctx.Writer.Write([]byte(param))
		ctx.Writer.Write([]byte(param1))
	})

	engine.Handle("POST", "/login", func(ctx *gin.Context) {
		fmt.Println("请求路径", ctx.FullPath())
		// 获取参数  方式一
		name := ctx.PostForm("name")
		pwd := ctx.PostForm("pwd")
		fmt.Println(name, pwd)
		// 获取参数  方式二
		name2, name2exist := ctx.GetPostForm("name")
		pwd2, pwd2exist := ctx.GetPostForm("pwd")
		if name2exist {
			fmt.Println(name2)
		}
		if pwd2exist {
			fmt.Println(pwd2)
		}
		
		ctx.Writer.Write([]byte("login"))
	})

	// GET form表单数据绑定结构体
	engine.Handle("GET", "/regist", func(ctx *gin.Context) {
		fmt.Println("请求路径", ctx.FullPath())
		// form表单数据绑定结构体
		var student Student
		
		if err := ctx.ShouldBindQuery(&student); err != nil {
			fmt.Println("出错了")
		}
		fmt.Println(student)
		ctx.Writer.Write([]byte("regist"))
	})

	// POST form表单数据绑定结构体
	engine.Handle("POST", "/newRegister", func(ctx *gin.Context) {
		fmt.Println("请求路径", ctx.FullPath())
		// form表单数据绑定结构体
		var regist Regist
		if err := ctx.ShouldBind(&regist); err != nil {
			fmt.Println("出错了")
		}
		// 返回数据格式一
		// ctx.Writer.Write([]byte("newRegister"))
		// 返回数据格式二
		// ctx.Writer.WriteString("newRegister")
		// 返回数据格式三
		ctx.JSON(200, map[string]interface{}{
			"code": 0,
			"message": "ok",
			"data": regist,
		})
	})
	
	// POST JSON数据格式
	engine.Handle("POST", "/addStudent", func(ctx *gin.Context) {
		fmt.Println("请求路径", ctx.FullPath())
		// form表单数据绑定结构体
		var student Student
		if err := ctx.BindJSON(&student); err != nil {
			fmt.Println("出错了")
		}

		// 返回数据格式一
		// ctx.Writer.Write([]byte("newRegister"))
		// 返回数据格式二
		// ctx.Writer.WriteString("newRegister")
		// 返回数据格式三
		resp := Response{
			Code: 1,
			Message: "OK",
			Data: student,
		}
		ctx.JSON(200, &resp)
	})

	// 路由组
	routerGroup := engine.Group("/user")
	// 把匿名函数封装到外面去：registHandle()
	routerGroup.POST("/register", registHandle)

	// 中间件（全局使用）
	engine.Use(RequestInof())
	engine.GET("/query", func(ctx *gin.Context) {
		fmt.Println("中间件解析执行........")
		ctx.JSON(202, map[string]interface{}{
			"code": 0,
			"message": ctx.FullPath(),
			"data": "中间件测试",
		})
	})

	// 中间件（单独使用）
	// 只针对某一个方法使用
	// engine.GET("/query", RequestInof(), func(ctx *gin.Context) {
	// 	ctx.JSON(200, map[string]interface{}{
	// 		"code": 0,
	// 		"message": ctx.FullPath(),
	// 		"data": "中间件测试",
	// 	})
	// })

	// msql的使用
	// 连接数据库
	connStr := "root:linyifan@tcp(127.0.0.1:3306)/gindemo"
	db, err := sql.Open("mysql", connStr)
	if err != nil {
		fmt.Println(err)
		return
	}
	// 创建数据库表  
	// 创建以后得注释掉
	// _, err = db.Exec("create table person(" +
	// 	"id int auto_increment primary key," +
	// 	"name varchar(12) not null," +
	// 	"age int default 1" +
	// ");")
	// if err != nil {
	// 	log.Fatal(err.Error())
	// 	return
	// }

	// 插入数据
	// _, err = db.Exec("insert into person(name, age) values(?, ?);", "joson.1in", 11)
	// if err == nil {
	// 	fmt.Println("数据插入成功")
	// }
	
	// 查询
	result, err := db.Query("select id, name, age from person")
	if err != nil {
		fmt.Println("查询失败")
	}
	scan:
		if result.Next() {
			var person Person
			err = result.Scan(&person.Id, &person.Name, &person.Age)
			if err != nil {
				log.Fatal(err.Error())
				return
			}
			fmt.Println(person)
			goto scan
		} 

	// 不填默认8080端口
	if err := engine.Run(":8080"); err != nil {
		log.Fatal(err.Error())
	}
}

type Person struct {
	Id int
	Name string
	Age int
}

// 打印请求信息的中间件
func RequestInof() gin.HandlerFunc {
	return func(context *gin.Context) {
		path := context.FullPath()
		method := context.Request.Method
		fmt.Println("请求路由：", path, method)

		// 当添加context.Next() 会先去执行中间件的解析(也就是调用中间件的地方)
		// 当执行解析中间件的地方完成后再回来接着往下执行
		context.Next()

		// 获取处理完后的状态信息
		fmt.Println("处理完成后的状态信息：", context.Writer.Status())
	}
}

func registHandle(ctx *gin.Context) {
	fmt.Println("请求路径", ctx.FullPath())
	// form表单数据绑定结构体
	var regist Regist
	if err := ctx.ShouldBind(&regist); err != nil {
		fmt.Println("出错了")
	}
	// 返回数据格式一
	// ctx.Writer.Write([]byte("newRegister"))
	// 返回数据格式二
	// ctx.Writer.WriteString("newRegister")
	// 返回数据格式三
	ctx.JSON(200, map[string]interface{}{
		"code": 0,
		"message": "ok1111111",
		"data": regist,
	})
}

type Response struct {
	Code int
	Message string
	Data interface{}
}

// form表单数据绑定结构体
// tag标签的form表示是form表单的name属性名字
type Student struct {
	Username string	`form:"username"`
	Classes string `form:"classes"`
}

// form表单数据绑定结构体
// tag标签的form表示是form表单的name属性名字
type Regist struct {
	Username string	`form:"username"`
	Classes string `form:"classes"`
}
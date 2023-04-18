package main

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"io"
	"ohmygin/dbconfig"
	"ohmygin/middlewares"
	"ohmygin/nacosconfig"
	"ohmygin/pojo"
	"ohmygin/router"
	"os"

	"github.com/gin-gonic/gin"
)

func setupLogging() {
	f, _ := os.Create("gin.log")
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)
}

func main() {
	//日志输出文件,必须放在最上面
	//强制日志颜色化
	gin.ForceConsoleColor()
	setupLogging()

	//默认logger和recover中间件
	r := gin.Default()
	//r.Use(middlewares.Logger(), gin.BasicAuth(gin.Accounts{"admin": "admin"}))

	//注册自定义validator
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("validPw", middlewares.ValidatePw)
		//注册自定义的构造体validator
		v.RegisterStructValidation(middlewares.ValidateSize, pojo.Users{})
	}

	// router.GET("/test", func(ctx *gin.Context) {
	// 	fmt.Println("ping")
	// 	ctx.JSON(200, gin.H{
	// 		"msg": "pong!",
	// 	})
	// })
	v1 := r.Group("/v1")
	router.UserRouterV1(v1)

	ch := make(chan int)
	//连接nacos
	go nacosconfig.Init(ch)

	<-ch
	//连接数据库
	go dbconfig.MysqlConnect()
	//连接redis
	go dbconfig.RedisConnect()

	r.Run(":1234")
}

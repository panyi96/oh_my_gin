package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"io"
	dbconfig2 "ohmygin/src/dbconfig"
	"ohmygin/src/middlewares"
	"ohmygin/src/nacosconfig"
	"ohmygin/src/pojo"
	"ohmygin/src/router"
	"os"
	"strconv"
)

func setupLogging() {
	f, _ := os.Create("gin.log")
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)
}

func main() {
	fmt.Println("version 1.0.0")
	ipAddr, port, grpcPort := getNacosParamFromEnv()
	fmt.Printf("=== USE NACOS CONFIG=== \nNACOS_IP=%s\nPORT=%d\nGRPC_PORT=%d\n",
		ipAddr, port, grpcPort)
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
	go nacosconfig.Init(ipAddr, port, grpcPort, ch)

	<-ch
	//连接数据库
	go dbconfig2.MysqlConnect()
	//连接redis
	go dbconfig2.RedisConnect()

	r.Run(":1234")
}

func getNacosParamFromEnv() (string, uint64, uint64) {
	ipAddr := os.Getenv("NACOS_IP")
	if "" == ipAddr {
		fmt.Println("nacos_ip env is blank, use default")
		ipAddr = "asd"
	}
	portStr := os.Getenv("PORT")
	if "" == portStr {
		portStr = "8848"
	}
	grpcPortStr := os.Getenv("GRPC_PORT")
	if "" == grpcPortStr {
		grpcPortStr = "9848"
	}
	port, err := strconv.ParseUint(portStr, 0, 64)
	if nil != err {
		panic(err)
	}
	grpcPort, err := strconv.ParseUint(grpcPortStr, 0, 64)
	if nil != err {
		panic(err)
	}
	return ipAddr, port, grpcPort
}

package router

import (
	"ohmygin/src/middlewares"
	"ohmygin/src/pojo"
	service2 "ohmygin/src/service"

	"github.com/gin-gonic/gin"
)

func UserRouterV1(r *gin.RouterGroup) {
	//todo swagger

	routerUser := r.Group("/user", middlewares.SetSession())

	//routerUser.GET("/find", service.FindUserById)
	routerUser.GET("/find", service2.CacheUser(service2.FindUserById, "id", "userId:%s", pojo.User{}))
	routerUser.POST("/add", service2.AddUser)
	routerUser.POST("/addList", service2.AddUsers)
	//routerUser.DELETE("/delete", service.DelUser)
	routerUser.PUT("/update", service2.UpdateUser)
	routerUser.POST("/upload", service2.Upload)
	routerUser.POST("/uploads", service2.Uploads)
	routerUser.POST("/login", service2.Login)

	routerUser.Use(middlewares.AuthSession)
	routerUser.DELETE("/delete", service2.DelUser)
	routerUser.GET("/logout", service2.LoginOut)

}

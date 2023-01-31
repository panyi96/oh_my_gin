package router

import (
	"ohmygin/middlewares"
	"ohmygin/pojo"
	"ohmygin/service"

	"github.com/gin-gonic/gin"
)

func UserRouterV1(r *gin.RouterGroup) {
	routerUser := r.Group("/user", middlewares.SetSession())

	//routerUser.GET("/find", service.FindUserById)
	routerUser.GET("/find", service.CacheUser(service.FindUserById, "id", "userId:%s", pojo.User{}))
	routerUser.POST("/add", service.AddUser)
	routerUser.POST("/addList", service.AddUsers)
	//routerUser.DELETE("/delete", service.DelUser)
	routerUser.PUT("/update", service.UpdateUser)
	routerUser.POST("/upload", service.Upload)
	routerUser.POST("/uploads", service.Uploads)
	routerUser.POST("/login", service.Login)

	routerUser.Use(middlewares.AuthSession)
	routerUser.DELETE("/delete", service.DelUser)
	routerUser.GET("/logout", service.LoginOut)

}

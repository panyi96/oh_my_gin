package middlewares

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"net/http"
)

const userKey = "session_id"

// 使用cookie存储sessionId
func SetSession() gin.HandlerFunc {
	store := cookie.NewStore([]byte(userKey))
	return sessions.Sessions("oh_my_gin_session", store)
}

// 用户鉴权session中间件
func AuthSession(ctx *gin.Context) {
	session := sessions.Default(ctx)
	sessionId := session.Get(userKey)
	if sessionId == nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"msg": "need login",
		})
		return
	}
	ctx.Next()
}

// 保存用户session
func SaveUserSession(ctx *gin.Context, userId uint) {
	session := sessions.Default(ctx)
	session.Set(userKey, userId)
	session.Save()
}

// 清除用户session
func ClearUserSession(ctx *gin.Context) {
	session := sessions.Default(ctx)
	session.Clear()
	session.Save()
}

// 获取用户session
func GetUserSession(ctx *gin.Context) uint {
	session := sessions.Default(ctx)
	sessionId := session.Get(userKey)
	if sessionId == nil {
		return 0
	}
	return sessionId.(uint)
}

// 检查用户session
func CheckUserSession(ctx *gin.Context) bool {
	sessionId := GetUserSession(ctx)
	if sessionId == 0 {
		return false
	}
	return true
}

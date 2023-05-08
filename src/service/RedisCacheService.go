package service

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"net/http"
	"ohmygin/src/dbconfig"
)

func CacheUser(handler gin.HandlerFunc, paramField string, keyPattern string, _struct interface{}) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userId := ctx.Query(paramField)
		redisKey := fmt.Sprintf(keyPattern, userId)
		val, err := dbconfig.RedisClient.Get(goCtx, redisKey).Result()
		//has data
		if nil == err {
			err := json.Unmarshal([]byte(val), &_struct)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, err)
				return
			}
			ctx.JSON(http.StatusOK, gin.H{
				"msg":  "from redis",
				"data": _struct,
			})
			return
		}
		//no data
		if err == redis.Nil {
			//not exist
			handler(ctx)
			//get from request threadLocal
			dbResult, _ := ctx.Get("dbResult")
			dbResultJSONStr, _ := json.Marshal(dbResult)
			_, err := dbconfig.RedisClient.Set(goCtx, "userId:"+ctx.Query("id"), string(dbResultJSONStr), 0).Result()
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, err)
				return
			}
		}
	}
}

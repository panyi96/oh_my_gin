package service

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"ohmygin/middlewares"
	"ohmygin/pojo"
	"strconv"

	"github.com/gin-gonic/gin"
)

var goCtx = context.Background()

const dbResultKey = "dbResult"

// FindUserById Get user
func FindUserById(ctx *gin.Context) {
	//业务逻辑
	userId, err := strconv.ParseUint(ctx.Query("id"), 0, 64)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, "获取参数发生错误！"+err.Error())
		fmt.Println("获取参数发生错误！" + err.Error())
		return
	}
	user, err := pojo.FindAllUserById(userId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "from db",
		"data": user,
	})
	ctx.Set(dbResultKey, user)
}

//func FindAllUserByIdFromRedis(ctx *gin.Context) {
//	var userList []pojo.User
//	idParam := ctx.Query("id")
//	if "" == idParam {
//		users, err := pojo.FindAllUser()
//		if err != nil {
//			ctx.JSON(http.StatusInternalServerError, err.Error())
//		} else {
//			userList = users
//		}
//	} else {
//		//get from redis
//		val, err := dbconfig.RedisClient.Get(goCtx, "userId:"+ctx.Query("id")).Result()
//		if err != nil {
//			fmt.Println("Get failed", err)
//		}
//		if "" == val || err == redis.Nil {
//			//get from db
//			userId, err := strconv.ParseUint(ctx.Query("id"), 0, 64)
//			if err != nil {
//				ctx.JSON(http.StatusInternalServerError, "获取参数发生错误！"+err.Error())
//				fmt.Println("获取参数发生错误！" + err.Error())
//				return
//			}
//			user, err := pojo.FindAllUserById(userId)
//			if err != nil {
//				ctx.JSON(http.StatusInternalServerError, err.Error())
//			} else {
//				userList = []pojo.User{user}
//				data, _ := json.Marshal(userList)
//				_, err := dbconfig.RedisClient.Set(goCtx, "userId:"+ctx.Query("id"), string(data), 0).Result()
//				if err != nil {
//					ctx.JSON(http.StatusInternalServerError, err)
//					return
//				}
//			}
//		} else {
//			err := json.Unmarshal([]byte(val), &userList)
//			if err != nil {
//				ctx.JSON(http.StatusInternalServerError, err)
//				return
//			}
//		}
//	}
//	ctx.JSON(http.StatusOK, userList)
//}

// AddUser Post user
func AddUser(ctx *gin.Context) {
	var user pojo.User
	err := ctx.ShouldBindJSON(&user)
	if err != nil {
		fmt.Println("获取参数发生错误！" + err.Error())
		ctx.JSON(http.StatusInternalServerError, "获取参数发生错误！"+err.Error())
		return
	}
	addUser, err := pojo.AddUser(user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
	} else {
		ctx.JSON(http.StatusOK, addUser)
	}
}

// AddUsers Post user
func AddUsers(ctx *gin.Context) {
	var users pojo.Users
	err := ctx.BindJSON(&users)
	if err != nil {
		fmt.Println("获取参数发生错误！" + err.Error())
		ctx.JSON(http.StatusInternalServerError, "获取参数发生错误！"+err.Error())
		return
	}
	addUsers, err := pojo.AddUsers(users)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
	} else {
		ctx.JSON(http.StatusOK, addUsers)
	}
}

// DelUser DELETE user
func DelUser(ctx *gin.Context) {
	idStr := ctx.Query("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		fmt.Println("强转int发生错误！" + err.Error())
	}
	row, err := pojo.DelUser(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
	} else {
		ctx.JSON(http.StatusOK, row)
	}
}

// UpdateUser PUT user
func UpdateUser(ctx *gin.Context) {
	newUser := pojo.User{}
	err := ctx.BindJSON(&newUser)
	if err != nil {
		fmt.Println("强转int发生错误！" + err.Error())
	}
	user, err := pojo.UpdateUser(newUser)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
	} else {
		ctx.JSON(http.StatusOK, user)
	}
}

func Upload(ctx *gin.Context) {
	// form-data 单文件
	file, _ := ctx.FormFile("file")
	log.Println(file.Filename)

	dst := "./" + file.Filename
	// 上传文件至指定的完整文件路径
	ctx.SaveUploadedFile(file, dst)

	ctx.String(http.StatusOK, fmt.Sprintf("'%s' uploaded!", file.Filename))
}

func Uploads(ctx *gin.Context) {
	// Multipart form
	form, _ := ctx.MultipartForm()
	//form-data多文件
	files := form.File["files"]

	for _, file := range files {
		log.Println(file.Filename)
		dst := "./" + file.Filename
		// 上传文件至指定目录
		ctx.SaveUploadedFile(file, dst)
	}
	ctx.String(http.StatusOK, fmt.Sprintf("%d files uploaded!", len(files)))
}

func Login(ctx *gin.Context) {
	userName := ctx.PostForm("userName")
	password := ctx.PostForm("password")
	user, err := pojo.GetUserByPw(userName, password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
	} else if user.Id == 0 {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": "login fail",
		})
	} else {
		middlewares.SaveUserSession(ctx, user.Id)
		ctx.JSON(http.StatusOK, gin.H{
			"message": "login successfully",
			"user":    user,
			"session": middlewares.GetUserSession(ctx),
		})
	}
}

func LoginOut(ctx *gin.Context) {
	middlewares.ClearUserSession(ctx)
	ctx.JSON(http.StatusOK, gin.H{
		"message": "logout successfully",
	})
}

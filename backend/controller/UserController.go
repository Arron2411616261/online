package controller

import (
	"gin/common"
	"gin/model"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"log"
	"math/rand"
	"net/http"
	"time"
)

// 注册接口函数
func Register(ctx *gin.Context) {
	//获取数据
	DB := common.GetDB()
	var receiveUser model.User
	if err := ctx.BindJSON(&receiveUser); err != nil {
		ctx.JSON(422, gin.H{"code": 422, "msg": "获取失败"})
		return
	}

	if len(receiveUser.Telephone) != 11 {
		ctx.JSON(422, gin.H{"code": 422, "msg": "手机号必须为11位"})
		return
	}
	if len(receiveUser.Password) < 6 || len(receiveUser.Password) > 14 {
		ctx.JSON(422, gin.H{"code": 422, "msg": "密码长度需要设置为6-14个字符"})
		return
	}

	//验证手机号是否被注册过
	if isTelephoneExist(DB, receiveUser.Telephone) {
		ctx.JSON(422, gin.H{"code": 422, "msg": "该手机号已被注册"})
		return
	}
	//对密码进行加密处理
	HashPassword, err := bcrypt.GenerateFromPassword([]byte(receiveUser.Password), bcrypt.DefaultCost)
	if err != nil {
		ctx.JSON(500, gin.H{"code": 500, "msg": "加密错误"})
		return
	}
	// 如果名称为空，随机创建一个名称
	if len(receiveUser.Name) == 0 {
		receiveUser.Name = RandomName(10)
	}

	newUser := model.User{
		Gender:    receiveUser.Gender,
		Name:      receiveUser.Name,
		Telephone: receiveUser.Telephone,
		Password:  string(HashPassword),
		Avatar:    receiveUser.Avatar,
	}
	DB.Create(&newUser)
	ctx.JSON(200, gin.H{
		"msg": "注册成功",
	})
}

// 登录接口函数
func Login(ctx *gin.Context) {
	//获取参数
	DB := common.GetDB()
	var receiveUser model.User
	if err := ctx.BindJSON(&receiveUser); err != nil {
		ctx.JSON(422, gin.H{"code": 422, "msg": "获取失败"})
		return
	}

	//数据验证
	if len(receiveUser.Telephone) != 11 {
		ctx.JSON(422, gin.H{"code": 422, "msg": "手机号为11位"})
		return
	}
	if len(receiveUser.Password) < 6 || len(receiveUser.Password) > 14 {
		ctx.JSON(422, gin.H{"code": 422, "msg": "密码长度为6-14个字符"})
		return
	}

	//验证手机号对应的用户是否存在
	var user model.User
	DB.Where("telephone = ?", receiveUser.Telephone).First(&user)
	if user.ID == 0 {
		ctx.JSON(422, gin.H{"code": 422, "msg": "用户不存在"})
		return
	}

	//验证密码是否正确
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(receiveUser.Password)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "密码错误",
		})
		return
	}

	//发放token
	var err error
	user.Token, err = common.ReleaseToken(user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "系统异常"})
		log.Println("token generate error:%v", err)
	}

	ctx.JSON(200, gin.H{
		"code": 200,
		"msg":  "登录成功",
		"result": gin.H{
			"id":       user.ID,
			"account":  user.Telephone,
			"token":    user.Token,
			"avatar":   user.Avatar,
			"nickname": user.Name,
			"gendar":   user.Gender,
		},
	})
}

// 验证手机号是否已被注册
func isTelephoneExist(db *gorm.DB, telephone string) bool {
	var user model.User
	db.Where("telephone = ?", telephone).First(&user)
	if user.ID != 0 {
		return true
	}
	return false
}

// 随机创建用户名称
func RandomName(n int) string {
	var letters = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	result := make([]byte, n)
	rand.Seed(time.Now().Unix())
	for i := range result {
		result[i] = letters[rand.Intn(len(letters))]
	}
	return string(result)
}

func Info(ctx *gin.Context) {
	user, _ := ctx.Get("user")
	ctx.JSON(http.StatusOK, gin.H{"code": 200, "data": gin.H{"user": user}})
}

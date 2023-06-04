package controller

import (
	"fmt"
	"gin/common"
	"gin/model"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AllIdle struct { // "_2" 区分于commoditycontroller的AllIdle

	Id      string
	Name    string `gorm:"type:varchar(20);not null"`
	Picture string `gorm:"type:varchar(1024);not null"`
	Goods   []model.Goods
}

type SingleIdle struct {
	Id      string
	Name    string `gorm:"type:varchar(20);not null"`
	Picture string `gorm:"type:varchar(1024);not null"`
	Goods   []model.Goods
}

func GetGoods(ctx *gin.Context) {

	DB := common.GetDB()
	var result [4]AllIdle

	for i := 0; i < 4; i++ {

		var category model.Category
		DB.Table("categories").Where("id = ?", i+1).Find(&category)
		result[i].Id = category.Id
		result[i].Name = category.Name
		result[i].Picture = category.Picture

		var goods []model.Goods
		DB.Table("goods").Where("Cate_Id = ? AND is_sold=?", i+1, false).Find(&goods)
		result[i].Goods = append(result[i].Goods, goods...)
	}
	ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*")

	ctx.JSON(200, gin.H{
		"code":   "1",
		"msg":    "获取全部商品成功",
		"result": result,
	})
}

//暂且不考虑id转换错误

func GetOneGood(c *gin.Context) {
	db := common.GetDB()
	idStr := c.Query("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	//i := 1
	fmt.Println(1)
	if err != nil {
		fmt.Errorf("invalid id fomrat %v", err)
	}
	var target model.Goods
	db.Table("goods").Where("id = ?", id).First(&target)
	//if err != nil {
	//	fmt.Println("断点位于查表")
	//}
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.JSON(200, gin.H{
		"result": target,
	})
}

func RecentIdle(ctx *gin.Context) {
	DB := common.GetDB()

	NUM := ctx.DefaultQuery("limit", "4")
	IntNum, err := strconv.Atoi(NUM) // 函数原型 ：func Atoi(s string) (int, error)
	if err != nil {
		print(err)
		//do not thing
	}
	var count int64
	DB.Table("goods").Where("is_sold=?", false).Count(&count)
	if IntNum > int(count) {
		IntNum = int(count) //让返回的数目不大于库存
	}
	var recentGoods = make([]model.Goods, IntNum)

	for i := int(count); i > int(count)-IntNum; i-- {
		print(int(count) - i)
		DB.Table("goods").Where("id = ? AND is_sold=?", i, false).Find(&recentGoods[int(count)-i])
	}
	ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*")

	ctx.JSON(200, gin.H{
		"code":   "1",
		"msg":    "获取最近发布成功",
		"result": recentGoods,
	})

}

func ChooseCategory(ctx *gin.Context) {

	DB := common.GetDB()
	var result SingleIdle

	Cate_id := ctx.DefaultQuery("id", "3")

	var category model.Category
	DB.Table("categories").Where("id = ?", Cate_id).Find(&category)
	result.Id = category.Id
	result.Name = category.Name
	result.Picture = category.Picture

	var goods []model.Goods
	DB.Table("goods").Where("cate_Id = ? AND is_sold=?", Cate_id, false).Find(&goods)
	result.Goods = append(result.Goods, goods...)

	ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*")

	ctx.JSON(200, gin.H{
		"code":   "1",
		"msg":    "获取分类下属物品成功",
		"result": result,
	})

}

type OrderInfo struct {
	Id        string `json:"goodId"` //goods_Id
	AddressId string `json:"addressId"`
}

func CreateOrder(ctx *gin.Context) {

	user, exist := ctx.Get("user")
	if exist == false {
		ctx.JSON(http.StatusUnauthorized, gin.H{"msg": "user not exist"})
		return
	}
	user.ID
	var orderInfo OrderInfo
	//绑定结构体
	if err := ctx.BindJSON(&orderInfo); err != nil {
		// 返回错误信息
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//生成订单表
	//买家id（user）、商品id（body）、地址id（body）？、订单id（gorm）、

	//返回id
	ctx.JSON(200, gin.H{
		//订单id
	})

	//地址表：接口一  用户id？
	//2  用token验证，body参数与user建表
	//3  ？
	//1  ？
}

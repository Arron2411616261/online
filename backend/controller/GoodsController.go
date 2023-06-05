package controller

import (
	"crypto/rand"
	"fmt"
	"gin/common"
	"gin/model"
	"github.com/gin-gonic/gin"
	"math/big"
	"net/http"
	"strconv"
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
	id, _ := strconv.ParseInt(idStr, 10, 64)
	//i := 1
	fmt.Println(1)
	// if err != nil {
	// 	fmt.Errorf("invalid id fomrat %v", err)
	// }
	var target model.Goods
	db.Table("goods").Where("id = ?", id).First(&target)
	//if err != nil {
	//	fmt.Println("断点位于查表")
	//}

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

	ctx.JSON(200, gin.H{
		"code":   "1",
		"msg":    "获取分类下属物品成功",
		"result": result,
	})

}

func RecommendGoods(ctx *gin.Context) {
	DB := common.GetDB()
	NUM := ctx.DefaultQuery("limit", "4")
	IntNum, err := strconv.Atoi(NUM) // 函数原型 ：func Atoi(s string) (int, error)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get limit"})
		return
	}

	result := make([]model.Goods, IntNum)
	var count int64
	if err := DB.Model(&model.Goods{}).Count(&count).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get count"})
		return
	}

	var idRecord = make([]uint, IntNum)
	for i := 0; i < IntNum; i++ {
		var ranGood model.Goods
		for {
			ranID, _ := rand.Int(rand.Reader, big.NewInt(count))
			id := int(ranID.Int64())
			//str_ranID := strconv.Itoa(id)
			DB.Table("goods").Where("id = ? AND is_sold=?", id+1, false).Find(&ranGood)
			if checkRanID(idRecord, i+1, ranGood.ID) {
				break
			}
		}
		idRecord[i] = ranGood.ID
		result[i] = ranGood
	}

	//user, _ := ctx.Get("user")
	//userinfo := user.(model.User)

	ctx.JSON(200, gin.H{
		"code":   200,
		"msg":    "操作成功",
		"result": result,
	})
}

//检查查询商品结果的ID号，如果重复或者没有对应的商品，则返回false
func checkRanID(idArray []uint, num int, checkID uint) bool {
	if checkID == 0 {
		return false
	}
	for i := 0; i < num; i++ {
		if idArray[i] == checkID {
			return false
		}
	}
	return true
}

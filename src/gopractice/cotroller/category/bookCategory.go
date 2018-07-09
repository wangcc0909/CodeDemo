package category

import (
	"github.com/gin-gonic/gin"
	"gopractice/cotroller/common"
	"gopractice/model"
	"fmt"
	"net/http"
	"github.com/microcosm-cc/bluemonday"
	"strings"
	"unicode/utf8"
	"strconv"
)

//图书分类列表
func BookCategoryList(c *gin.Context) {
	sendErrJson := common.SendErrJson
	var categories []model.Category

	if err := model.DB.Order("sequence asc").Find(&categories).Error; err != nil {
		fmt.Println(err.Error())
		sendErrJson("error", c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"errNo": model.ErrorCode.SUCCESS,
		"msg":   "success",
		"data": gin.H{
			"categories": categories,
		},
	})
}

//保存图书分类
func SaveBookCategory(c *gin.Context, isEdit bool) {
	sendErrJson := common.SendErrJson

	minOrder := model.MinOrder
	maxOrder := model.MaxOrder

	var category model.BookCategory
	if err := c.ShouldBindJSON(&category); err != nil {
		sendErrJson("参数无效", c)
		return
	}

	category.Name = bluemonday.UGCPolicy().Sanitize(category.Name)
	category.Name = strings.TrimSpace(category.Name)

	if category.Name == "" {
		sendErrJson("分类名称不能为空", c)
		return
	}

	if utf8.RuneCountInString(category.Name) > model.MaxNameLen {
		msg := "分类名称不能超过" + fmt.Sprintf("%d", model.MaxNameLen) + "个字符"
		sendErrJson(msg, c)
		return
	}

	if category.Sequence < minOrder || category.Sequence > maxOrder {
		msg := "分类排序要在" + strconv.Itoa(minOrder) + "到" + strconv.Itoa(maxOrder) + "之间"
		sendErrJson(msg, c)
		return
	}

	if category.ParentID != 0 {
		var parentCate model.BookCategory
		if err := model.DB.First(&parentCate, category.ParentID).Error; err != nil {
			sendErrJson("无效的父类ID", c)
			return
		}
	}

	var updateCategory model.BookCategory
	if !isEdit { //创建分类
		if err := model.DB.Create(&category).Error; err != nil {
			sendErrJson("error", c)
			return
		}
	} else {
		//更新分类
		if err := model.DB.First(&updateCategory, category.ID).Error; err == nil {
			updateMap := make(map[string]interface{})
			updateMap["name"] = category.Name
			updateMap["sequence"] = category.Sequence
			updateMap["parent_id"] = category.ParentID
			if err := model.DB.Model(&updateMap).Updates(updateMap).Error; err != nil {
				fmt.Println(err.Error())
				sendErrJson("error", c)
				return
			}
		} else {
			sendErrJson("无效的分类ID", c)
			return
		}
	}

	var categoryJson model.BookCategory
	if isEdit {
		categoryJson = updateCategory
	} else {
		categoryJson = category
	}

	c.JSON(http.StatusOK, gin.H{
		"errNo": model.ErrorCode.SUCCESS,
		"msg":   "success",
		"data": gin.H{
			"category": categoryJson,
		},
	})
}

//创建图书分类
func CreateBookCategory(c *gin.Context) {
	SaveBookCategory(c, false)
}

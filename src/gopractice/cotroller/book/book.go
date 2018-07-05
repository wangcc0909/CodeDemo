package book

import (
	"github.com/gin-gonic/gin"
	"gopractice/cotroller/common"
	"gopractice/model"
	"strconv"
	"fmt"
	"net/http"
)

func Create(c *gin.Context) {

}

//获取图书列表
func List(c *gin.Context) {
	sendErrJson := common.SendErrJson

	var books []model.Book
	var categoryId int

	cateStr := c.Query("cateId")

	if cateStr == "" {
		categoryId = 0
	} else {
		var err error
		if categoryId, err = strconv.Atoi(cateStr); err != nil {
			fmt.Println(err.Error())
			sendErrJson("分类ID不正确", c)
			return
		}
	}

	var queryCMD = model.DB.Model(&model.Book{}).Where("read_limits <> ?", model.BookReadLimitsPrivate).
		Where("status <> ?", model.BookVerifyFail).Where("status <> ?", model.BookUnpublish)

	if categoryId != 0 {
		queryCMD = queryCMD.Joins("JOIN book_category on book.id = book_category.book_id").
			Joins("JOIN book_categories ON book_category.book_category_id = book_categories.id AND book_categories.id = ?", categoryId)
	}

	if err := queryCMD.Order("created_at DESC").Find(&books).Error; err != nil {
		fmt.Println(err.Error())
		sendErrJson("error", c)
		return
	}

	for i := 0; i < len(books); i++ {
		if err := model.DB.Model(&books[i]).Related(&books[i].User, "users").Error; err != nil {
			fmt.Println(err.Error())
			sendErrJson("error", c)
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"errNo": model.ErrorCode.SUCCESS,
		"msg":   "success",
		"data": gin.H{
			"books": books,
		},
	})

}

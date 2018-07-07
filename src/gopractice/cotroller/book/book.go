package book

import (
	"github.com/gin-gonic/gin"
	"gopractice/cotroller/common"
	"gopractice/model"
	"strconv"
	"fmt"
	"net/http"
	"gopractice/util"
	"strings"
	"unicode/utf8"
)

func Create(c *gin.Context) {
	Save(c, false)
}

func Update(c *gin.Context) {
	Save(c, true)
}

func Save(c *gin.Context, isEdit bool) {
	sendErrJson := common.SendErrJson
	var book model.Book

	if err := c.ShouldBindJSON(&book); err != nil {
		sendErrJson("参数错误", c)
		return
	}

	if book.ContentType != model.ContentTypeMarkdown && book.ContentType != model.ContentTypeHTML {
		sendErrJson("无效的图书格式", c)
		return
	}

	book.Name = strings.TrimSpace(book.Name)
	book.Name = util.AvoidXss(book.Name)

	book.Content = strings.TrimSpace(book.Content)
	book.HTMLContent = strings.TrimSpace(book.HTMLContent)

	if book.HTMLContent != "" {
		book.HTMLContent = util.AvoidXss(book.HTMLContent)
	}

	if book.Name == "" {
		sendErrJson("图书的名字不能为空", c)
		return
	}

	if utf8.RuneCountInString(book.Name) > model.MaxNameLen {
		msg := "图书的名字不能超过" + fmt.Sprintf("%d", model.MaxNameLen) + "个字符"
		sendErrJson(msg, c)
		return
	}

	if book.ReadLimits != model.BookReadLimitsPublic && book.ReadLimits != model.BookReadLimitsPrivate &&
		book.ReadLimits != model.BookReadLimitsPay {
		sendErrJson("无效的阅读权限", c)
		return
	}

	var theContent string
	if book.ContentType == model.ContentTypeMarkdown {
		theContent = book.Content
	} else {
		theContent = book.HTMLContent
	}

	contentCount := utf8.RuneCountInString(theContent)
	if theContent == "" || contentCount <= 0 {
		sendErrJson("图书简介不能为空", c)
		return
	}

	if contentCount > model.MaxContentLen {
		msg := "图书简介不能超过" + fmt.Sprintf("%d", model.MaxContentLen) + "个字符"
		sendErrJson(msg, c)
		return
	}

	if book.Categories == nil || len(book.Categories) <= 0 {
		sendErrJson("请选择分类", c)
		return
	}

	if len(book.Categories) > model.MaxCategoriesLen {
		msg := "图书最多属于" + fmt.Sprintf("%d", model.MaxCategoriesLen) + "个分类"
		sendErrJson(msg, c)
		return
	}

	for i := 0; i < len(book.Categories); i++ {
		var category model.BookCategory
		if err := model.DB.First(&category, book.Categories[i].ID).Error; err != nil {
			sendErrJson("无效的分类ID", c)
			return
		}
		book.Categories[i] = category
	}

	iUser, _ := c.Get("user")
	user := iUser.(model.User)

	var updateBook model.Book

	if !isEdit {
		//创建图书
		book.Status = model.BookUnpublish
		book.UserID = user.ID
		//创建图书时可以选择格式,之后不能修改
		if err := model.DB.Create(&book).Error; err != nil {
			fmt.Println(err.Error())
			sendErrJson("error", c)
			return
		}
	} else {
		var sql = "DELETE FROM book_category WHERE book_id = ?"
		if err := model.DB.Exec(sql, book.ID).Error; err != nil {
			fmt.Println(err.Error())
			sendErrJson("error", c)
			return
		}
		//更新图书
		if err := model.DB.First(&updateBook, book.ID).Error; err == nil {
			if updateBook.UserID != user.ID {
				sendErrJson("您没有权限执行此操作", c)
				return
			}

			updateBook.ReadLimits = book.ReadLimits
			updateBook.Name = book.Name
			updateBook.Categories = book.Categories
			updateBook.CoverURL = book.CoverURL
			updateBook.Content = book.Content
			updateBook.HTMLContent = book.HTMLContent
			if err := model.DB.Save(&updateBook).Error; err != nil {
				fmt.Println(err.Error())
				sendErrJson("error", c)
				return
			}
		} else {
			sendErrJson("无效的图书ID", c)
			return
		}
	}

	var bookJson model.Book
	if isEdit {
		bookJson = updateBook
	} else {
		bookJson = book
	}

	c.JSON(http.StatusOK, gin.H{
		"errNo": model.ErrorCode.SUCCESS,
		"msg":   "success",
		"data": gin.H{
			"book": bookJson,
		},
	})

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

func MyBook(c *gin.Context) {
	sendErrJson := common.SendErrJson
	var books []model.Book
	iUser, _ := c.Get("user")
	user := iUser.(model.User)

	var pageNo int
	var pageNoErr error
	if pageNo, pageNoErr = strconv.Atoi(c.Query("pageNo")); pageNoErr != nil {
		pageNo = 1
	}

	if pageNo < 1 {
		pageNo = 1
	}

	pageSize := model.PageSize

	offset := (pageNo - 1) * pageSize
	if err := model.DB.Model(&model.Book{}).Where("user_id = ?", user.ID).Offset(offset).Limit(pageSize).
		Order("created_at DESC").Find(&books).Error; err != nil {
		fmt.Println(err.Error())
		sendErrJson("error", c)
		return
	}

	var totalCount int
	if err := model.DB.Model(&model.Book{}).Where("user_id = ?", user.ID).Count(&totalCount).Error; err != nil {
		fmt.Println(err.Error())
		sendErrJson("error", c)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"errNo": model.ErrorCode.SUCCESS,
		"msg":   "success",
		"data": gin.H{
			"books":      books,
			"pageNo":     pageNo,
			"pageSize":   pageSize,
			"totalCount": totalCount,
		},
	})
}

//用户公开的图书列表
func UserPublicBooks(c *gin.Context) {
	sendErrJson := common.SendErrJson

	userID, idErr := strconv.Atoi(c.Param("userID"))

	if idErr != nil {
		sendErrJson("无效的userID", c)
		return
	}

	var pageNo int
	var pageNoErr error

	if pageNo, pageNoErr = strconv.Atoi(c.Query("pageNo")); pageNoErr != nil {
		pageNo = 1
	}

	if pageNo < 1 {
		pageNo = 1
	}

	pageSize := model.PageSize

	offset := (pageNo - 1) * pageSize

	var books []model.Book
	if err := model.DB.Model(&model.Book{}).Where("read_limits <> ?", model.BookReadLimitsPrivate).
		Where("status <> ?", model.BookVerifyFail).Where("status <> ?", model.BookUnpublish).
		Where("user_id = ?", userID).Offset(offset).Limit(pageSize).Find(&books).Error; err != nil {
		fmt.Println(err.Error())
		sendErrJson("error", c)
		return
	}

	var totalCount int
	if err := model.DB.Model(&model.Book{}).Where("read_limits <> ?", model.BookReadLimitsPrivate).
		Where("status <> ?", model.BookVerifyFail).Where("status <> ?", model.BookUnpublish).
		Where("user_id = ?", userID).Offset(offset).Limit(pageSize).Count(&totalCount).Error; err != nil {
		fmt.Println(err.Error())
		sendErrJson("error", c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"errNo": model.ErrorCode.SUCCESS,
		"msg":   "success",
		"data": gin.H{
			"books":      books,
			"pageNo":     pageNo,
			"pageSize":   pageSize,
			"totalCount": totalCount,
		},
	})
}

func Info(c *gin.Context) {
	sendErrJson := common.SendErrJson

	bookId, idErr := strconv.Atoi(c.Param("id"))
	if idErr != nil {
		sendErrJson("无效的ID", c)
		return
	}

	var book model.Book

	if err := model.DB.Where("status != ?", model.BookVerifyFail).First(&book, bookId).Error; err != nil {
		sendErrJson("无效的图书ID", c)
		return
	}

	if book.ReadLimits == model.BookReadLimitsPrivate {
		iUser, _ := c.Get("user")

		if iUser == nil {
			sendErrJson("没有权限", c)
			return
		}

		user := iUser.(model.User)

		if user.ID != book.UserID {
			sendErrJson("没有权限", c)
			return
		}
	}

	if err := model.DB.Model(&book).Related(&book.Categories, "categories").Error; err != nil {
		fmt.Println(err.Error())
		sendErrJson("error", c)
		return
	}

	if c.Query("f") != "md" {
		if book.ContentType == model.ContentTypeMarkdown {
			book.HTMLContent = util.MarkdownToHTML(book.Content)
		} else if book.ContentType == model.ContentTypeHTML {
			book.HTMLContent = util.AvoidXss(book.HTMLContent)
		} else {
			book.HTMLContent = util.MarkdownToHTML(book.Content)
		}
		book.Content = ""
	}

	c.JSON(http.StatusOK, gin.H{
		"errNo": model.ErrorCode.SUCCESS,
		"msg":   "success",
		"data": gin.H{
			"book": book,
		},
	})

}

// Chapters 获取图书的所有章节, 若图书是私有的，那么只有作者本人才能查看
func Chapters(c *gin.Context) {
	sendErrJson := common.SendErrJson

	bookID, idErr := strconv.Atoi(c.Param("bookID"))
	if idErr != nil {
		sendErrJson("无效的bookID", c)
		return
	}

	var book model.Book

	if err := model.DB.Where("status != ?", model.BookVerifyFail).First(&book, bookID).Error; err != nil {
		fmt.Println(err.Error())
		sendErrJson("无效的bookID", c)
		return
	}

	if book.ReadLimits == model.BookReadLimitsPrivate {
		iUser, _ := c.Get("user")

		if iUser == nil {
			sendErrJson("没有权限", c)
			return
		}

		user := iUser.(model.User)

		if user.ID != book.UserID {
			sendErrJson("没有权限", c)
			return
		}
	}

	var chapters []model.BookChapter
	if err := model.DB.Model(&model.BookChapter{}).Where("book_id = ?", bookID).Select("id,name,parent_id").
		Order("created_at DESC").Find(chapters).Error; err != nil {
		fmt.Println(err.Error())
		sendErrJson("error", c)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"errNo": model.ErrorCode.SUCCESS,
		"msg":   "success",
		"data": gin.H{
			"chapters": chapters,
		},
	})
}

func Chapter(c *gin.Context) {
	sendErrJson := common.SendErrJson

	chapterID, idErr := strconv.Atoi(c.Param("chapterID"))
	if idErr != nil {
		sendErrJson("无效的ChapterID", c)
		return
	}

	var chapter model.BookChapter

	if err := model.DB.First(&chapter, chapterID).Error; err != nil {
		fmt.Println(err.Error())
		sendErrJson("无效的ChapterID", c)
		return
	}

	var book model.Book

	if err := model.DB.First(&book, chapter.BookID).Error; err != nil {
		//如果图书被删了 章节不让查询
		fmt.Println(err.Error())
		sendErrJson("未找到对应的图书", c)
		return
	}

	if book.ReadLimits == model.BookReadLimitsPrivate {
		iUser, _ := c.Get("user")

		if iUser == nil {
			sendErrJson("没有权限", c)
			return
		}

		user := iUser.(model.User)

		if user.ID != book.UserID {
			sendErrJson("没有权限", c)
			return
		}
	}

	if c.Query("f") != "md" {
		if chapter.ContentType == model.ContentTypeMarkdown {
			chapter.HTMLContent = util.MarkdownToHTML(book.Content)
		} else if book.ContentType == model.ContentTypeHTML {
			chapter.HTMLContent = util.AvoidXss(book.HTMLContent)
		} else {
			chapter.HTMLContent = util.MarkdownToHTML(book.Content)
		}
		chapter.Content = ""
	}

	c.JSON(http.StatusOK, gin.H{
		"errNo": model.ErrorCode.SUCCESS,
		"msg":   "success",
		"data": gin.H{
			"chapter": chapter,
		},
	})
}

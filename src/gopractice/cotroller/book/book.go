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
	"time"
)

//创建图书
func Create(c *gin.Context) {
	Save(c, false)
}

//更新图书
func Update(c *gin.Context) {
	Save(c, true)
}

//更新名字
func UpdateName(c *gin.Context) {
	sendErrJson := common.SendErrJson
	var bookData model.Book
	if err := c.ShouldBindJSON(&bookData); err != nil {
		sendErrJson("参数错误", c)
		return
	}

	bookData.Name = util.AvoidXss(bookData.Name)
	bookData.Name = strings.TrimSpace(bookData.Name)

	if bookData.Name == "" {
		sendErrJson("图书名称不能为空", c)
		return
	}

	if utf8.RuneCountInString(bookData.Name) > model.MaxNameLen {
		msg := "图书的名称不能超过" + fmt.Sprintf("%d", model.MaxNameLen) + "个字符"
		sendErrJson(msg, c)
		return
	}

	var book model.Book

	if err := model.DB.First(&book, bookData.ID).Error; err != nil {
		sendErrJson("无效的图书ID", c)
		return
	}

	iUer, _ := c.Get("user")
	user := iUer.(model.User)

	if book.UserID != user.ID {
		sendErrJson("没有权限", c)
		return
	}

	book.Name = bookData.Name

	if err := model.DB.Save(&book).Error; err != nil {
		fmt.Println(err.Error())
		sendErrJson("error", c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"errNo": model.ErrorCode.SUCCESS,
		"msg":   "success",
		"data": gin.H{
			"book": book,
		},
	})
}

//保存图书
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

//创建图书的章节
func CreateChapter(c *gin.Context) {
	sendErrJson := common.SendErrJson

	type ReqData struct {
		Name     string `json:"name" binding:"required,min=1,max=100"`
		ParentID uint   `json:"parentId"`
		BookID   uint   `json:"bookId"`
	}

	var reqData ReqData
	if err := c.ShouldBindJSON(&reqData); err != nil {
		fmt.Println(err.Error())
		sendErrJson("参数错误", c)
		return
	}

	reqData.Name = util.AvoidXss(reqData.Name)
	reqData.Name = strings.TrimSpace(reqData.Name)

	if reqData.Name == "" {
		sendErrJson("章节的名称不能为空", c)
		return
	}

	var chapter model.BookChapter
	chapter.Name = reqData.Name
	chapter.ParentID = reqData.ParentID
	chapter.BookID = reqData.BookID

	if chapter.ParentID != model.NoParent {
		var parentChapter model.BookChapter
		if err := model.DB.First(&parentChapter, chapter.ParentID).Error; err != nil {
			sendErrJson("无效的parentID", c)
			return
		}
	}

	var book model.Book
	if err := model.DB.First(&book, chapter.BookID).Error; err != nil {
		sendErrJson("无效的bookId", c)
		return
	}

	iUser, _ := c.Get("user")
	user := iUser.(model.User)

	if book.UserID != user.ID {
		sendErrJson("没有权限", c)
		return
	}

	chapter.ContentType = book.ContentType
	chapter.UserID = user.ID

	if err := model.DB.Create(&chapter).Error; err != nil {
		fmt.Println(err.Error())
		sendErrJson("error", c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"errNo": model.ErrorCode.SUCCESS,
		"msg":   "success",
		"data": gin.H{
			"chapter": chapter,
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

//获取某用户的所有书籍
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

//获取图书的信息
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

//获取章节
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

//发布图书
func Publish(c *gin.Context) {
	sendErrJson := common.SendErrJson
	id, idErr := strconv.Atoi(c.Param("bookID"))
	if idErr != nil {
		sendErrJson("无效的bookID", c)
		return
	}

	var book model.Book
	if err := model.DB.First(&book, id).Error; err != nil {
		sendErrJson("无效的bookID", c)
		return
	}

	iUser, _ := c.Get("user")
	user := iUser.(model.User)

	if book.UserID != user.ID {
		sendErrJson("您没有权限执行此操作", c)
		return
	}

	book.Status = model.BookVerifySuccess
	if err := model.DB.Save(&book).Error; err != nil {
		sendErrJson("error", c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"errNo": model.ErrorCode.SUCCESS,
		"msg":   "success",
		"data": gin.H{
			"book": book,
		},
	})

}

//更新图书的章节内容
func UpdateChapterContent(c *gin.Context) {
	sendErrJson := common.SendErrJson
	type ReqData struct {
		ID          uint   `json:"chapterID"`
		Content     string `json:"content"`
		HTMLContent string `json:"htmlContent"`
	}

	var reqData ReqData
	if err := c.ShouldBindJSON(&reqData); err != nil {
		fmt.Println(err.Error())
		sendErrJson("参数错误", c)
		return
	}

	reqData.Content = strings.TrimSpace(reqData.Content)
	reqData.HTMLContent = strings.TrimSpace(reqData.HTMLContent)

	if reqData.HTMLContent != "" {
		reqData.HTMLContent = util.AvoidXss(reqData.HTMLContent)
	}

	var chapter model.BookChapter
	if err := model.DB.First(&chapter, reqData.ID).Error; err != nil {
		sendErrJson("错误的章节ID", c)
		return
	}

	iUser, _ := c.Get("user")
	user := iUser.(model.User)

	if chapter.UserID != user.ID {
		sendErrJson("您没有权限执行此操作", c)
		return
	}

	chapter.Content = reqData.Content
	chapter.HTMLContent = reqData.HTMLContent

	if err := model.DB.Save(&chapter).Error; err != nil {
		fmt.Println(err.Error())
		sendErrJson("error", c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"errNo": model.ErrorCode.SUCCESS,
		"msg":   "success",
		"data": gin.H{
			"id": reqData.ID,
		},
	})
}

//更新章节的名字
func UpdateChapterName(c *gin.Context) {
	sendErrJson := common.SendErrJson
	type ReqData struct {
		ID   uint   `json:"id"`
		Name string `json:"name" binding:"required,min=1,max=100"`
	}
	var reqData ReqData
	if err := c.ShouldBindJSON(&reqData); err != nil {
		sendErrJson("参数错误", c)
		return
	}

	reqData.Name = util.AvoidXss(reqData.Name)
	reqData.Name = strings.TrimSpace(reqData.Name)

	if reqData.Name == "" {
		sendErrJson("章节的名字不能为空", c)
		return
	}

	if utf8.RuneCountInString(reqData.Name) > model.MaxNameLen {
		msg := "章节的名称不能超过" + fmt.Sprintf("%d", model.MaxNameLen) + "个字符"
		sendErrJson(msg, c)
		return
	}

	var chapter model.BookChapter

	if err := model.DB.First(&chapter, reqData.ID).Error; err != nil {
		fmt.Println(err.Error())
		sendErrJson("无效的章节ID", c)
		return
	}

	iUser, _ := c.Get("user")
	user := iUser.(model.User)

	if chapter.UserID != user.ID {
		sendErrJson("您没有权限执行此操作", c)
		return
	}

	chapter.Name = reqData.Name
	if err := model.DB.Save(&chapter).Error; err != nil {
		fmt.Println(err.Error())
		sendErrJson("error", c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"errNo": model.ErrorCode.SUCCESS,
		"msg":   "success",
		"data": gin.H{
			"chapter": chapter,
		},
	})
}

//删除图书
func Delete(c *gin.Context) {
	sendErrJson := common.SendErrJson
	bookID, idErr := strconv.Atoi(c.Param("id"))

	if idErr != nil {
		sendErrJson("无效的ID", c)
		return
	}

	var book model.Book
	if err := model.DB.First(&book, bookID).Error; err != nil {
		fmt.Println(err.Error())
		sendErrJson("无效的ID", c)
		return
	}

	iUser, _ := c.Get("user")
	user := iUser.(model.User)

	if book.UserID != user.ID {
		sendErrJson("您没有权限执行此操作", c)
		return
	}

	if err := model.DB.Delete(&book).Error; err != nil {
		fmt.Println(err.Error())
		sendErrJson("error", c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"errNo": model.ErrorCode.SUCCESS,
		"msg":   "success",
		"data": gin.H{
			"id": bookID,
		},
	})
}

func DeleteChapter(c *gin.Context) {
	sendErrJson := common.SendErrJson
	chapterID, idErr := strconv.Atoi(c.Param("chapterID"))
	if idErr != nil {
		sendErrJson("无效的chapterID", c)
		return
	}

	var chapter model.BookChapter
	if err := model.DB.First(&chapter, chapterID).Error; err != nil {
		sendErrJson("无效的chapterID", c)
		return
	}

	iUser, _ := c.Get("user")
	user := iUser.(model.User)

	if chapter.UserID != user.ID {
		sendErrJson("您没有权限执行此操作", c)
		return
	}

	var sql = "UPDATE book_chapters SET delete_at = ? WHERE id = ? OR parent_id = ?"
	if err := model.DB.Exec(sql, time.Now(), chapterID, chapterID).Error; err != nil {
		fmt.Println(err.Error())
		sendErrJson("error", c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"errNo": model.ErrorCode.SUCCESS,
		"msg":   "success",
		"data": gin.H{
			"id": chapterID,
		},
	})

}

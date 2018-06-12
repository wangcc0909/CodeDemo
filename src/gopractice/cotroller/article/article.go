package article

import (
	"github.com/gin-gonic/gin"
	"gopractice/cotroller/common"
	"gopractice/model"
	"strconv"
	"time"
	"fmt"
	"strings"
	"net/http"
	"math"
)

func queryList(c *gin.Context, isBackend bool) {
	sendErrJson := common.SendErrJson
	var articles []model.Article
	var categoryID int
	var pageNo int
	var err error
	var startTime string
	var endTime string
	if pageNo, err = strconv.Atoi(c.Query("pageNo")); err != nil {
		pageNo = 1
		err = nil
	}

	if pageNo < 1 {
		pageNo = 1
	}

	pageSize := 40
	offset := (pageNo - 1) * pageSize

	if startAt, err := strconv.Atoi(c.Query("startAt")); err != nil {
		startTime = time.Unix(0, 0).Format("2006-01-02 17:26:35")
	} else {
		startTime = time.Unix(int64(startAt/1000), 0).Format("2006-01-02 17:26:35")
	}

	if endAt, err := strconv.Atoi(c.Query("endAt")); err != nil {
		endTime = time.Now().Format("2006-01-02 17:26:35")
	} else {
		endTime = time.Unix(int64(endAt/1000), 0).Format("2006-01-02 17:26:35")
	}

	//默认按创建时间降序
	var orderField = "created_at"
	var orderASC = "DESC"
	if c.Query("asc") == "1" {
		orderASC = "ASC"
	} else {
		orderASC = "DESC"
	}

	cateIDStr := c.Query("cateId")
	if cateIDStr == "" {
		categoryID = 0
	} else if categoryID, err = strconv.Atoi(cateIDStr); err != nil {
		fmt.Println("categoryId err :", err)
		sendErrJson("分类ID不正确", c)
		return
	}

	var topArticles []model.TopArticle
	if err := model.DB.Find(&topArticles).Error; err != nil {
		sendErrJson("error", c)
		return
	}

	var topArr []string
	for i := 0; i < len(topArticles); i++ {
		topArr = append(topArr, strconv.Itoa(int(topArticles[i].ArticleID)))
	}

	topIDs := strconv.Itoa(model.NoParent)

	if len(topArr) > 0 {
		topIDs = strings.Join(topArr, ",")
	}

	type TotalCountResult struct {
		TotalCount int
	}

	var totalCountResult TotalCountResult

	if categoryID != 0 {
		var category model.Category
		if model.DB.First(&category, categoryID).Error != nil {
			sendErrJson("分类ID不正确", c)
			return
		}

		var sql = `SELECT distinct(articles.id),articles.name,articles.browse_count,articles.comment_count,articles.collect_count,
						articles.status,articles.created_at,articles.updated_at,articles.user_id,articles.last_user_id 
					FROM articles,article_category
					WHERE articles.id = article_category.article_id
					{statusSQL}
					AND article_category.category_id = {categoryID}
					AND articles.deleted_at IS NULL
					AND articles.id NOT IN ({topIDS})
					{timeSQL}
					ORDER BY {orderField} {orderASC}
					LIMIT {offset},{pageSize}`
		sql = strings.Replace(sql, "{categoryID}", strconv.Itoa(categoryID), -1)
		sql = strings.Replace(sql, "{orderField}", orderField, -1)
		sql = strings.Replace(sql, "{topIDS}", topIDs, -1)
		sql = strings.Replace(sql, "{timeSQL}", "AND created_at >= '"+startTime+"' AND created_at <'"+endTime+"'", -1)
		sql = strings.Replace(sql, "{orderASC}", orderASC, -1)
		sql = strings.Replace(sql, "{offset}", strconv.Itoa(offset), -1)
		sql = strings.Replace(sql, "{pageSize}", strconv.Itoa(pageSize), -1)

		if isBackend {
			sql = strings.Replace(sql, "{statusSQL}", "", -1)
		} else {
			sql = strings.Replace(sql, "{statusSQL}", " AND (status = 1 OR status = 2)", -1)
		}

		if err := model.DB.Raw(sql).Scan(&articles).Error; err != nil {
			sendErrJson("error", c)
			return
		}

		for i := 0; i < len(articles); i++ {
			articles[i].Categories = []model.Category{category}
		}

		countSQL := `SELECT COUNT(distinct(articles.id)) AS total_count
					FROM articles,article_category
					WHERE articles.id = article_category.article_id
					{statusSQL}
					AND article_category.category_id = {categoryID}
					AND articles.id NOT IN ({topIDS})
					{timeSQL}
					AND articles.deleted_at IS NULL`
		countSQL = strings.Replace(countSQL, "{categoryID}", strconv.Itoa(categoryID), -1)
		countSQL = strings.Replace(countSQL, "{topIDS}", topIDs, -1)
		countSQL = strings.Replace(countSQL, "{timeSQL}", "AND created_at >= '"+startTime+"' AND created_at <'"+endTime+"'", -1)

		if isBackend {
			//管理员查询话题列表时,会返回审核未通过的话题
			countSQL = strings.Replace(countSQL, "{statusSQL}", "", -1)
			if err := model.DB.Raw(countSQL).Scan(&totalCountResult).Error; err != nil {
				sendErrJson("error", c)
				return
			}
		} else {
			countSQL = strings.Replace(countSQL, "{statusSQL}", " AND (status = 1 OR status = 2)", -1)
			if err := model.DB.Raw(countSQL).Scan(&totalCountResult).Error; err != nil {
				sendErrJson("error", c)
				return
			}
		}
	} else {
		orderStr := orderField + " " + orderASC
		excludeIDs := "id NOT IN ({topIDs})"
		excludeIDs = strings.Replace(excludeIDs, "{topIDs}", topIDs, -1)

		if isBackend {
			//管理员查询话题列表时,会返回审核未通过的话题
			err = model.DB.Where(excludeIDs).
				Where("created_at >= ? AND created_at < ?", startTime, endTime).
				Offset(offset).Limit(pageSize).
				Order(orderStr).Find(&articles).Error
		} else {
			err = model.DB.Where(excludeIDs).Where("status = 1 OR status = 2").
				Where("created_at >= ? AND created_at < ?", startTime, endTime).
				Offset(offset).Limit(pageSize).
				Order(orderStr).Find(&articles).Error
		}

		if err != nil {
			sendErrJson("error", c)
			return
		}

		for i := 0; i < len(articles); i++ {
			if err = model.DB.Model(articles[i]).Related(&articles[i].Categories, "categories").Error; err != nil {
				fmt.Println(err)
				sendErrJson("error", c)
				return
			}
		}

		if isBackend {
			//管理员查询话题列表时,会返回审核未通过的话题
			err = model.DB.Model(&model.Article{}).Where(excludeIDs).
				Where("created_at >= ? AND created_at < ?", startTime, endTime).
				Count(&totalCountResult.TotalCount).Error

			if err != nil {
				sendErrJson("error", c)
				return
			}
		} else {
			err = model.DB.Model(&model.Article{}).Where(excludeIDs).Where("status = 1 OR status = 2").
				Where("created_at >= ? AND created_at < ?", startTime, endTime).
				Count(&totalCountResult.TotalCount).Error
			if err != nil {
				sendErrJson("error", c)
				return
			}
		}
	}

	for i := 0; i < len(articles); i++ {
		if err := model.DB.Model(&articles[i]).Related(&articles[i].User, "users").Error; err != nil {
			fmt.Println(err)
			sendErrJson("error", c)
			return
		}

		if articles[i].LastUserID != 0 {
			if err := model.DB.Model(&articles[i]).Related(&articles[i].LastUser, "users", "last_user_id").Error; err != nil {
				fmt.Println(err)
				sendErrJson("error", c)
				return
			}
		}

		if c.Query("noContent") == "true" {
			articles[i].Content = ""
			articles[i].HTMLContent = ""
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"errNo": model.ErrorCode.SUCCESS,
		"msg":   "success",
		"data": gin.H{
			"articles":   articles,
			"pageNo":     pageNo,
			"pageSize":   pageSize,
			"totalPage":  math.Ceil(float64(totalCountResult.TotalCount) / float64(pageSize)),
			"totalCount": totalCountResult.TotalCount,
		},
	})

}

//文章列表
func List(c *gin.Context) {
	queryList(c, false)
}

package crawler

import (
	"github.com/gin-gonic/gin"
	"gopractice/cotroller/common"
	"gopractice/model"
	"github.com/PuerkitoBio/goquery"
	"gopractice/util"
	"net/http"
	"fmt"
	"strings"
	"net/url"
	"os"
	"io"
	"gopractice/config"
	"sync"
)

type CrawlSelector struct {
	Form                  int
	ListItemSelector      string
	ListItemTitleSelector string
	TitleSelector         string
	ContentSelector       string
}

func CreateCrawlSelector(from int) CrawlSelector {
	selector := CrawlSelector{
		Form: from,
	}

	switch from {
	case model.ArticleFromNull:
		selector.ListItemSelector = ""
		selector.ListItemTitleSelector = ""
		selector.TitleSelector = ""
		selector.ContentSelector = ""
	case model.ArticleFromJianShu:
		selector.ListItemSelector = ".note-list li"
		selector.ListItemTitleSelector = ".title"
		selector.TitleSelector = ".article .title"
		selector.ContentSelector = ".show-content"
	case model.ArticleFromZhihu:
		selector.ListItemSelector = ".PostListItem"
		selector.ListItemTitleSelector = ".PostListInfo a"
		selector.TitleSelector = ".PostIndex-title"
		selector.ContentSelector = ".PostIndex-content"
	case model.ArticleFromHuXiu:
		selector.ListItemSelector = ".mod-art"
		selector.ListItemTitleSelector = ".mob-cct h2 a"
		selector.TitleSelector = ".t-h1"
		selector.ContentSelector = ".article-content-wrap"
	case model.ArticleFromCustom:
		selector.ListItemSelector = ""
		selector.ListItemTitleSelector = ""
		selector.TitleSelector = ""
		selector.ContentSelector = ""
	}

	return selector
}

func CreateSourceHTML(from int) string {
	var htmlArr []string
	switch from {
	case model.ArticleFromJianShu:
		htmlArr = []string{
			"<div id=\"golang123-content-outter-footer\">",
			"<blockquote>",
			"<p>来源: <a href=\"https://www.jianshu.com/\" target=\"_blank\">简书</a><br>",
			"原文: <a href=\"{articleURL}\" target=\"_blank\">{title}</a></p>",
			"</blockquote>",
			"</div>",
		}
	case model.ArticleFromZhihu:
		htmlArr = []string{
			"<div id=\"golang123-content-outter-footer\">",
			"<blockquote>",
			"<p>来源: <a href=\"https://www.zhihu.com\" target=\"_blank\">知乎</a><br>",
			"原文: <a href=\"{articleURL}\" target=\"_blank\">{title}</a></p>",
			"</blockquote>",
			"</div>",
		}
	case model.ArticleFromHuXiu:
		htmlArr = []string{
			"<div id=\"golang123-content-outter-footer\">",
			"<blockquote>",
			"<p>来源: <a href=\"https://www.huxiu.com\" target=\"_blank\">虎嗅</a><br>",
			"原文: <a href=\"{articleURL}\" target=\"_blank\">{title}</a></p>",
			"</blockquote>",
			"</div>",
		}
	case model.ArticleFromCustom:
		htmlArr = []string{
			"<div id=\"golang123-content-outter-footer\">",
			"<blockquote>",
			"<p>来源: <a href=\"{siteURL}\" target=\"_blank\">{siteName}</a><br>",
			"原文: <a href=\"{articleURL}\" target=\"_blank\">{title}</a></p>",
			"</blockquote>",
			"</div>",
		}
	case model.ArticleFromNull:
		htmlArr = []string{}
	}
	return strings.Join(htmlArr, "")
}

func CrawContent(pageUrl string, crawlSelector CrawlSelector,
	siteInfo map[string]string, isExist bool) map[string]string {
	var crawlArticle model.CrawlerArticle

	if err := model.DB.Where("url = ?", pageUrl).Find(&crawlArticle).Error; err == nil {
		if !isExist {
			return nil
		}
	}
	articleDoc, err := goquery.NewDocument(pageUrl)
	if err != nil {
		return nil
	}

	title := articleDoc.Find(crawlSelector.TitleSelector).Text()
	if title == "" && crawlSelector.Form != model.ArticleFromNull {
		return nil
	}

	contentDoc := articleDoc.Find(crawlSelector.ContentSelector)
	imgs := contentDoc.Find("img")

	if imgs.Length() > 0 {
		imgs.Each(func(i int, img *goquery.Selection) {
			imgUrl, isExist := img.Attr("src")
			var ext string
			if !isExist {
				if crawlSelector.Form != model.ArticleFromJianShu {
					return
				}

				originalSrc, originalExist := img.Attr("data-original-src")
				if originalExist && originalSrc != "" {
					tempImgUrl, err := util.RelativeURLConvertAbsoluteURL(originalSrc, pageUrl)
					if err != nil || tempImgUrl == "" {
						return
					}

					imgUrl = tempImgUrl

					resp, err := http.Head(imgUrl)
					if err != nil {
						fmt.Println(err.Error())
						return
					}

					defer resp.Body.Close()

					contentType := resp.Header.Get("content-type")

					if contentType == "image/jpeg" {
						ext = ".jpg"
					} else if contentType == "image/gif" {
						ext = ".gif"
					} else if contentType == "image/png" {
						ext = ".png"
					}
				}
			}

			if imgUrl == "" || crawlSelector.Form == model.ArticleFromJianShu &&
				strings.Index(imgUrl, "data:image/svg+xml;utf8,") == 0 {
				actualSrc, actualExist := img.Attr("data-actualsrc")

				if actualExist && actualSrc != "" {
					imgUrl = actualSrc
				}
			}

			imgUrl, err := util.RelativeURLConvertAbsoluteURL(imgUrl, pageUrl)

			if err != nil || imgUrl == "" {
				return
			}

			urlData, urlErr := url.Parse(imgUrl)
			if urlErr != nil {
				return
			}

			if ext == "" {
				index := strings.LastIndex(urlData.Path, ".")
				if index > 0 {
					ext = urlData.Path[index:]
				}
			}

			resp, err := http.Get(imgUrl)
			if err != nil {
				return
			}

			defer resp.Body.Close()

			imageUploadInfo := model.GenerateUploadedImgInfo(ext)

			if err := os.MkdirAll(imageUploadInfo.UploadDir, 0777); err != nil {
				fmt.Println(err.Error())
				return
			}

			file, err := os.OpenFile(imageUploadInfo.UploadFilePath, os.O_WRONLY|os.O_CREATE, 0666)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			defer file.Close()

			_, err = io.Copy(file, resp.Body)
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			img.SetAttr("src", imageUploadInfo.ImgUrl)
		})
	}

	contentDoc.Find("a").Each(func(i int, a *goquery.Selection) {
		oldHref, exist := a.Attr("href")
		if exist {
			href, err := util.RelativeURLConvertAbsoluteURL(oldHref, pageUrl)
			if err == nil {
				a.SetAttr("href", href)
			}
		}
	})

	articleHtml, htmlErr := contentDoc.Html()
	if htmlErr != nil {
		return nil
	}

	sourceHTML := CreateSourceHTML(crawlSelector.Form)

	if crawlSelector.Form == model.ArticleFromCustom {
		sourceHTML = strings.Replace(sourceHTML, "{siteRUL}", siteInfo["siteURL"], -1)
		sourceHTML = strings.Replace(sourceHTML, "{siteName}", siteInfo["siteName"], -1)
	}

	sourceHTML = strings.Replace(sourceHTML, "{title}", title, -1)
	sourceHTML = strings.Replace(sourceHTML, "{articleURL}", pageUrl, -1)

	articleHtml += sourceHTML
	articleHtml = "<div id=\"golang123-content-outter\">" + articleHtml + "</div>"
	return map[string]string{
		"Title":   title,
		"Content": articleHtml,
		"URL":     pageUrl,
	}
}

func CrawlNotSaveContent(c *gin.Context) {
	sendErrJson := common.SendErrJson

	type JsonData struct {
		URL             string `json:"url"`
		TitleSelector   string `json:"titleSelector"`
		ContentSelector string `json:"contentSelector"`
	}

	var jsonData JsonData

	//获取请求参数
	if err := c.ShouldBindJSON(&jsonData); err != nil {
		sendErrJson("参数无效", c)
		return
	}

	crawlSelector := CreateCrawlSelector(model.ArticleFromNull)
	crawlSelector.TitleSelector = jsonData.TitleSelector
	crawlSelector.ContentSelector = jsonData.ContentSelector

	data := CrawContent(jsonData.URL, crawlSelector, nil, true)

	//返回数据
	c.JSON(http.StatusOK, gin.H{
		"errNo": model.ErrorCode.SUCCESS,
		"msg":   "success",
		"data":  data,
	})
}

//获取爬虫账户
func CrawlAccount(c *gin.Context) {
	sendErrJson := common.SendErrJson
	var users []model.User
	if err := model.DB.Where("name = ?", config.ServerConfig.CrawlerName).Find(&users).Error; err != nil {
		fmt.Println(err.Error())
		sendErrJson("error", c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"errNo": model.ErrorCode.SUCCESS,
		"msg":   "success",
		"data": gin.H{
			"users": users,
		},
	})
}

func createArticle(user model.User, category model.Category, from int, data map[string]string) {
	var article model.Article
	article.Name = data["Title"]
	article.HTMLContent = data["Content"]
	article.ContentType = model.ContentTypeHTML
	article.UserID = user.ID
	article.Status = model.ArticleVerifying
	article.Categories = append(article.Categories, category)

	var crawlArticle model.CrawlerArticle
	crawlArticle.URL = data["URL"]
	crawlArticle.Title = article.Name
	crawlArticle.Content = article.HTMLContent
	crawlArticle.From = from

	tx := model.DB.Begin()
	if err := tx.Create(&article).Error; err != nil {
		tx.Rollback()
		return
	}

	crawlArticle.ID = article.ID
	if err := tx.Create(&crawlArticle).Error; err != nil {
		tx.Rollback()
		return
	}
	tx.Commit()
}

//爬取文章
func CrawlList(listUrl string, user model.User, category model.Category, crawlSelector CrawlSelector, siteInfo map[string]string, crawlExist bool, wg *sync.WaitGroup) {
	defer wg.Done()
	if _, err := url.Parse(listUrl); err != nil {
		return
	}

	doc, docErr := goquery.NewDocument(listUrl)
	if docErr != nil {
		fmt.Println(docErr.Error())
		return
	}

	var articleURLArr []string
	doc.Find(crawlSelector.ListItemSelector).Each(func(i int, selection *goquery.Selection) {
		articleLink := selection.Find(crawlSelector.ListItemTitleSelector)
		fmt.Println(selection.Html())
		fmt.Println(articleLink.Html())
		href, exists := articleLink.Attr("href")
		if exists {
			urlTemp, err := util.RelativeURLConvertAbsoluteURL(href, listUrl)
			if err == nil {
				articleURLArr = append(articleURLArr, urlTemp)
			}
		}
	})

	for i := 0; i < len(articleURLArr); i++ {
		articleMap := CrawContent(articleURLArr[i], crawlSelector, siteInfo, crawlExist)
		if articleMap != nil {
			createArticle(user, category, crawlSelector.Form, articleMap)
		}
	}
}

//爬取文章
func Crawl(c *gin.Context) {
	sendErrJson := common.SendErrJson
	type JSONData struct {
		URLS       []string `json:"urls"`
		From       int      `json:"from"`
		CategoryID int      `json:"categoryId"`
		Scope      string   `json:"scope"`
		CrawlExist bool     `json:"crawlExist"`
	}

	var jsonData JSONData
	if err := c.ShouldBindJSON(&jsonData); err != nil {
		sendErrJson("参数无效", c)
		return
	}

	if jsonData.From != model.ArticleFromJianShu && jsonData.From != model.ArticleFromZhihu && jsonData.From != model.ArticleFromHuXiu {
		sendErrJson("无效的from", c)
		return
	}

	if jsonData.Scope != model.CrawlerScopePage && jsonData.Scope != model.CrawlerScopeList {
		sendErrJson("无效的scope", c)
		return
	}

	iUser, _ := c.Get("user")
	user := iUser.(model.User)

	if user.Name != config.ServerConfig.CrawlerName {
		sendErrJson("您没有权限执行此操作,请使用爬虫账号", c)
		return
	}

	var category model.Category
	if err := model.DB.First(&category, jsonData.CategoryID).Error; err != nil {
		fmt.Println(err.Error())
		sendErrJson("错误的categoryID", c)
		return
	}

	crawlSelector := CreateCrawlSelector(jsonData.From)

	if jsonData.Scope == model.CrawlerScopeList {
		var wg sync.WaitGroup
		for i := 0; i < len(jsonData.URLS); i++ {
			wg.Add(1)
			go CrawlList(jsonData.URLS[i], user, category, crawlSelector, nil, jsonData.CrawlExist, &wg)
		}
		wg.Wait()
	} else {
		for i := 0; i < len(jsonData.URLS); i++ {
			data := CrawContent(jsonData.URLS[i], crawlSelector, nil, jsonData.CrawlExist)
			if data != nil {
				createArticle(user, category, jsonData.From, data)
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"errNo": model.ErrorCode.SUCCESS,
		"msg":   "抓取完成",
		"data":  gin.H{},
	})
}

func CustomCrawl(c *gin.Context) {
	sendErrJson := common.SendErrJson
	type JSONData struct {
		URLS                  []string `json:"urls"`
		From                  int      `json:"from"`
		CategoryID            int      `json:"categoryId"`
		Scope                 string   `json:"scope"`
		CrawlExist            bool     `json:"crawlExist"`
		ListItemSelector      string   `json:"listItemSelector"`
		ListItemTitleSelector string   `json:"listItemTitleSelector"`
		TitleSelector         string   `json:"titleSelector"`
		ContentSelector       string   `json:"contentSelector"`
		SiteURL               string   `json:"siteUrl" binding:"required,url"`
		SiteName              string   `json:"siteName" binding:"required"`
	}

	var jsonData JSONData
	if err := c.ShouldBindJSON(&jsonData); err != nil {
		sendErrJson("参数无效", c)
		return
	}

	if jsonData.From != model.ArticleFromCustom {
		sendErrJson("无效的from", c)
		return
	}

	if jsonData.Scope != model.CrawlerScopePage && jsonData.Scope != model.CrawlerScopeList {
		sendErrJson("无效的scope", c)
		return
	}

	iUser, _ := c.Get("user")
	user := iUser.(model.User)

	if user.Name != config.ServerConfig.CrawlerName {
		sendErrJson("您没有权限执行此操作,请使用爬虫账号", c)
		return
	}

	var category model.Category
	if err := model.DB.First(&category, jsonData.CategoryID).Error; err != nil {
		fmt.Println(err.Error())
		sendErrJson("无效的categoryID", c)
		return
	}

	crawlSelector := CreateCrawlSelector(jsonData.From)
	crawlSelector.TitleSelector = jsonData.TitleSelector
	crawlSelector.ListItemSelector = jsonData.ListItemSelector
	crawlSelector.ListItemTitleSelector = jsonData.ListItemTitleSelector
	crawlSelector.ContentSelector = jsonData.ContentSelector
	siteInfo := map[string]string{
		"siteURL":  jsonData.SiteURL,
		"siteName": jsonData.SiteName,
	}
	if jsonData.Scope == model.CrawlerScopeList {
		var wg sync.WaitGroup
		for i := 0; i < len(jsonData.URLS); i++ {
			wg.Add(1)
			go CrawlList(jsonData.URLS[i], user, category, crawlSelector, siteInfo, jsonData.CrawlExist, &wg)
		}
		wg.Wait()
	} else {
		for i := 0; i < len(jsonData.URLS); i++ {
			data := CrawContent(jsonData.URLS[i], crawlSelector, siteInfo, jsonData.CrawlExist)
			if data != nil {
				createArticle(user, category, jsonData.From, data)
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"errNo": model.ErrorCode.SUCCESS,
		"msg":   "抓取完成",
		"data":  gin.H{},
	})
}

//创建爬虫账号
func CreateAccount(c *gin.Context) {
	sendErrJson := common.SendErrJson
	var users []model.User
	if err := model.DB.Where("name = ?", config.ServerConfig.CrawlerName).Find(&users).Error; err != nil {
		fmt.Println(err.Error())
		sendErrJson("error", c)
		return
	}

	if len(users) <= 0 { //创建爬虫账号
		var user model.User
		user.Name = config.ServerConfig.CrawlerName
		user.Role = model.UserRoleCrawler
		user.AvatarURL = ""
		user.Status = model.UserStatusActived
		if err := model.DB.Save(&user).Error; err != nil {
			fmt.Println(err.Error())
			sendErrJson("error", c)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"errNo": model.ErrorCode.SUCCESS,
			"msg":   "success",
			"data":  []model.User{user},
		})
		return
	}

	sendErrJson("爬虫账号已存在", c)
}

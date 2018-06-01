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
	return ""
}

func CrawContent(pageUrl string, crawlSelector CrawlSelector,
	siteInfo map[string]string, isExsit bool) map[string]interface{} {
	var crawlArticle model.CrawlerArticle

	if err := model.DB.Where("url = ?", pageUrl).Find(&crawlArticle).Error; err == nil {
		if !isExsit {
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
	return map[string]interface{}{
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

	c.JSON(http.StatusOK, gin.H{
		"errNo": model.ErrorCode.SUCCESS,
		"msg":   "success",
		"data":  data,
	})

}

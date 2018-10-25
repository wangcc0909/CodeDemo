package crawler

import (
	"strings"
	"github.com/PuerkitoBio/goquery"
	"fmt"
	"net/url"
	"time"
	"net/http"
	"os"
	"io"
	"testAnswer/code9/model"
)


//爬虫爬取得文章
type CrawlerArticle struct {
	ID        uint       `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdateAt  time.Time  `json:"updateAt"`
	DeleteAt  *time.Time `sql:"index" json:"deleteAt"`
	URL       string     `json:"url"`
	Content   string     `json:"content"`
	Title     string     `json:"title"`
	From      int        `json:"from"`
	ArticleId uint       `json:"articleId"`
}

type crawlSelector struct {
	ListItemSelector      string
	ListItemTitleSelector string
	TitleSelector         string
	ContentSelector       string
}

func createCrawlSelector() crawlSelector {
	selector := crawlSelector{}

	selector.ListItemSelector = ".note-list li"
	selector.ListItemTitleSelector = ".title"
	selector.TitleSelector = ".article .title"
	selector.ContentSelector = ".show-content"

	return selector
}

func CrawlerContent(pageUrl string, crawlSelector crawlSelector) map[string]string {
	articleDOC, err := goquery.NewDocument(pageUrl)
	if err != nil {
		return nil
	}
	title := articleDOC.Find(crawlSelector.TitleSelector).Text()

	contentDOM := articleDOC.Find(crawlSelector.ContentSelector)

	imgs := contentDOM.Find("img")
	if imgs.Length() > 0 {
		imgs.Each(func(j int, img *goquery.Selection) {
			imgURL, exists := img.Attr("src")
			var ext string
			if !exists {

				originalSrc, originalExists := img.Attr("data-original-src")
				if originalExists && originalSrc != "" {
					tempImgURL, tempErr := RelativeURLConvertAbsoluteURL(originalSrc, pageUrl)
					if tempErr != nil || tempImgURL == "" {
						return
					}
					imgURL = tempImgURL
					resp, err := http.Head(imgURL)
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

			var imgURLErr error
			imgURL, imgURLErr = RelativeURLConvertAbsoluteURL(imgURL, pageUrl)
			if imgURLErr != nil || imgURL == "" {
				return
			}
			urlData, urlErr := url.Parse(imgURL)
			if urlErr != nil {
				return
			}

			if ext == "" {
				index := strings.LastIndex(urlData.Path, ".")
				if index >= 0 {
					ext = urlData.Path[index:]
				}
			}

			resp, err := http.Get(imgURL)

			if err != nil {
				return
			}

			defer resp.Body.Close()

			imgUploadedInfo := model.GenerateImgUploadedInfo(ext)
			if err := os.MkdirAll(imgUploadedInfo.UploadDir, 0777); err != nil {
				fmt.Println(err.Error())
				return
			}
			out, outErr := os.OpenFile(imgUploadedInfo.UploadFilePath, os.O_WRONLY|os.O_CREATE, 0666)
			if outErr != nil {
				fmt.Println(outErr.Error())
				return
			}

			defer out.Close()

			if _, err := io.Copy(out, resp.Body); err != nil {
				fmt.Println(err.Error())
				return
			}
			img.SetAttr("src", imgUploadedInfo.ImgURL)
		})
	}

	contentDOM.Find("a").Each(func(j int, a *goquery.Selection) {
		oldHref, exists := a.Attr("href")
		if exists {
			href, err := RelativeURLConvertAbsoluteURL(oldHref, pageUrl)
			if err == nil {
				a.SetAttr("href", href)
			}
		}
	})
	articleHTML, htmlErr := contentDOM.Html()
	if htmlErr != nil {
		return nil
	}

	sourceHTML := createSourceHTML()

	sourceHTML = strings.Replace(sourceHTML, "{title}", title, -1)
	sourceHTML = strings.Replace(sourceHTML, "{articleURL}", pageUrl, -1)
	articleHTML += sourceHTML
	articleHTML = "<div id=\"my-content-outter\">" + articleHTML + "</div>"
	return map[string]string{
		"Title":   title,
		"Content": articleHTML,
		"URL":     pageUrl,
	}
}

func createSourceHTML() string {
	var htmlArr []string

	htmlArr = []string{
		"<div id=\"my-content-outter-footer\">",
		"<blockquote>",
		"<p>来源: <a href=\"https://www.jianshu.com/\" target=\"_blank\">简书</a><br>",
		"原文: <a href=\"{articleURL}\" target=\"_blank\">{title}</a></p>",
		"</blockquote>",
		"</div>",
	}
	return strings.Join(htmlArr, "")
}

func RelativeURLConvertAbsoluteURL(curURL string, baseUrl string) (string, error) {

	curUrlData, err := url.Parse(curURL)
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}

	baseUrlData, err := url.Parse(baseUrl)
	if err != nil {
		return "", nil
	}

	curUrlData = baseUrlData.ResolveReference(curUrlData)

	return curUrlData.String(), nil
}

func CrawlerJianshu(pageUrl string) map[string]string {
	selector := createCrawlSelector()
	return CrawlerContent(pageUrl,selector)
}

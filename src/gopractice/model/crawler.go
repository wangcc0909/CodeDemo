package model

import (
	"time"
)

//爬虫爬取得文章
type CrawlerArticle struct {
	ID        uint       `gorm:"primary_key" json:"id"`
	CreateAt  time.Time  `json:"createAt"`
	UpdateAt  time.Time  `json:"updateAt"`
	DeleteAt  *time.Time `sql:"index" json:"deleteAt"`
	URL       string     `json:"url"`
	Content   string     `json:"content"`
	Title     string     `json:"title"`
	From      int        `json:"from"`
	ArticleId uint       `json:"articleId"`
}

const (
	//无来源
	ArticleFromNull = 0

	//简书
	ArticleFromJianShu = 1

	//知乎
	ArticleFromZhihu = 2

	//虎嗅
	ArticleFromHuXiu = 3

	//自定义
	ArticleFromCustom = 10
)

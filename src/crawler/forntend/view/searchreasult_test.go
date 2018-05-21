package view

import (
	"testing"
	"os"
	"crawler/forntend/model"
	model2 "crawler/model"
)

func TestSearchResultView(t *testing.T) {
	view := CreateSearchResultView("template.html")

	out,err := os.Create("template_test.html")
	if err != nil {
		panic(err)
	}

	//文件要先关闭
	defer out.Close()

	page := model.SearchResult{}

	page.Hint = 20

	item := model2.Item{
		Url:  "http://album.zhenai.com/u/108906739",
		Type: "zhenai",
		Id:   "108906739",
		Profile: model2.Profile{
			Age:        34,
			Height:     162,
			Weight:     57,
			InComing:     "3001-5000元",
			Gender:     "女",
			Name:       "安静的雪",
			XinZuo:     "牡羊座",
			Occupation: "人事/行政",
			Marriage:   "离异",
			Horse:      "已购房",
			HuKou:      "山东菏泽",
			Education:  "大学本科",
			Car:        "未购车",
		},
	}

	for i := 0;i < 10 ; i++ {
		page.Items = append(page.Items,item)
	}

	err = view.Render(out,page)
	if err != nil {
		t.Error(err)
	}
}

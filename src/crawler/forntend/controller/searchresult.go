package controller

import (
	"crawler/forntend/view"
	"gopkg.in/olivere/elastic.v3"
	"crawler/forntend/model"
	"golang.org/x/net/context"
	"reflect"
	model2 "crawler/model"
	"net/http"
	"regexp"
	"strings"
	"strconv"
	"crawler_distributed/config"
)

type SearchResultHandle struct {
	View   view.SearchResultView
	Client *elastic.Client
}

func (s SearchResultHandle) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	q := strings.TrimSpace(req.FormValue("q"))

	from, err := strconv.Atoi(req.FormValue("from"))
	if err != nil {
		from = 0
	}

	page, err := s.getSearchResult(q, from)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = s.View.Render(w, page)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func CreateSearchResultHandle(fileName string) SearchResultHandle {

	client, err := elastic.NewClient(
		elastic.SetSniff(false),
		elastic.SetURL("http://192.168.99.100:9200"),
	)

	if err != nil {
		panic(err)
	}

	return SearchResultHandle{
		View:   view.CreateSearchResultView(fileName),
		Client: client,
	}
}

const PageSize = 10

func (s SearchResultHandle) getSearchResult(q string, from int) (model.SearchResult, error) {

	var result model.SearchResult
	result.Query = q

	resp, err := s.Client.Search(config.RpcIndex).
		Query(elastic.NewQueryStringQuery(
			rewriteQueryString(q))).
		From(from).
		DoC(context.Background())

	if err != nil {
		return result, err
	}

	result.Hint = resp.TotalHits()
	result.Start = from
	result.Items = resp.Each(reflect.TypeOf(model2.Item{}))

	if result.Start == 0 {
		result.PreFrom = -1
	} else {
		result.PreFrom = (result.Start - 1) / PageSize * PageSize
	}

	result.NextFrom = result.Start + len(result.Items)

	return result, nil
}

func rewriteQueryString(q string) string {
	re := regexp.MustCompile(`([A-Z][a-z]*):`)
	return re.ReplaceAllString(q, "Profile.$1:")
}

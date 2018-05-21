package view

import (
	"html/template"
	"io"
	"crawler/forntend/model"
)

type SearchResultView struct {
	Template *template.Template
}

func CreateSearchResultView(fileName string) SearchResultView {

	return SearchResultView{
		Template:template.Must(template.ParseFiles(fileName)),
	}
}

func (s SearchResultView) Render(wr io.Writer, data model.SearchResult) error {
	return s.Template.Execute(wr,data)
}

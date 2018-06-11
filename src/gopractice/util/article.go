package util

import (
	"github.com/microcosm-cc/bluemonday"
	blackfriday "github.com/russross/blackfriday"
)

//markdown to html
func MarkdownToHTML(md string) string {
	myHTMLFlags := 0 |
		blackfriday.HTML_USE_XHTML |
		blackfriday.HTML_USE_SMARTYPANTS |
		blackfriday.HTML_SMARTYPANTS_FRACTIONS |
		blackfriday.HTML_SMARTYPANTS_DASHES |
		blackfriday.HTML_SMARTYPANTS_LATEX_DASHES
	myExtension := 0 |
		blackfriday.EXTENSION_NO_INTRA_EMPHASIS |
		blackfriday.EXTENSION_TABLES |
		blackfriday.EXTENSION_FENCED_CODE |
		blackfriday.EXTENSION_AUTOLINK |
		blackfriday.EXTENSION_STRIKETHROUGH |
		blackfriday.EXTENSION_SPACE_HEADERS |
		blackfriday.EXTENSION_HEADER_IDS |
		blackfriday.EXTENSION_BACKSLASH_LINE_BREAK |
		blackfriday.EXTENSION_DEFINITION_LISTS |
		blackfriday.EXTENSION_HARD_LINE_BREAK

	rendered := blackfriday.HtmlRenderer(myHTMLFlags,"","")
	bytes := blackfriday.MarkdownOptions([]byte(md),rendered,blackfriday.Options{
		Extensions:myExtension,
	})

	theHTML := string(bytes)
	return bluemonday.UGCPolicy().Sanitize(theHTML)

}

//避免 Xss
func AvoidXss(theHTML string) string {
	return bluemonday.UGCPolicy().Sanitize(theHTML)
}

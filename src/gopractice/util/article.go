package util

import "github.com/microcosm-cc/bluemonday"

//避免 Xss
func AvoidXss(theHTML string) string {
	return bluemonday.UGCPolicy().Sanitize(theHTML)
}

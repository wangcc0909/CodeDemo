package util

import (
	"net/url"
	"fmt"
)

func RelativeURLConvertAbsoluteURL(curURL string,baseUrl string) (string,error) {

	curUrlData,err := url.Parse(curURL)
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}

	baseUrlData,err := url.Parse(baseUrl)
	if err != nil {
		return "", nil
	}

	curUrlData = baseUrlData.ResolveReference(curUrlData)

	return curUrlData.String(),nil
}

package fetcher

import (
	"net/http"
	"bufio"
	"golang.org/x/text/transform"
	"io/ioutil"
	"golang.org/x/text/encoding"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding/unicode"
	"github.com/gpmgo/gopm/modules/log"
	"fmt"
	"time"
	"crawler_distributed/config"
)

var rateLimit = time.Tick(time.Second / config.RateQps)
func Fetch(url string) ([]byte,error) {
	<- rateLimit
	resp,err := http.Get(url)
	if err != nil {
		return nil,err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {

		return nil,fmt.Errorf("wrong status code : %d",
			resp.StatusCode)
	}

	bodyReader := bufio.NewReader(resp.Body)

	e := determineEncoding(bodyReader)

	utf8Reader := transform.NewReader(bodyReader,e.NewDecoder())

	return ioutil.ReadAll(utf8Reader)
}

func determineEncoding(r *bufio.Reader) encoding.Encoding {

	bytes,err := r.Peek(1024)
	if err != nil {
		log.Error("fetcher encoding error : %s",err)
		return unicode.UTF8
	}

	e,_,_ := charset.DetermineEncoding(bytes,"")
	return e
}
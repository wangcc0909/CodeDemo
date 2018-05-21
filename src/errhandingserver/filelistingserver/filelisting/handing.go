package filelisting

import (
	"net/http"
	"os"
	"io/ioutil"
	"strings"
	"fmt"
)

type userError string

func (r userError) Error() string {
	return r.Message()
}

func (r userError) Message() string {
	return string(r)
}

const prefix = "/list/"

func Hander(writer http.ResponseWriter, request *http.Request) error {

	//strings.Index 查找第二个参数在字符串中出现的位置
	if strings.Index(request.URL.Path,prefix) != 0 {
		return userError(fmt.Sprintf("path %s must be start with %s",request.URL.Path,prefix))
	}

	path := request.URL.Path[len(prefix):]

	file,err := os.Open(path)
	if err != nil {
		return err
	}

	defer file.Close()

	result,err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	writer.Write(result)
	return nil
}
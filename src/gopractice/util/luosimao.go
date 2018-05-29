package util

import (
	"errors"
	"net/url"
	"net/http"
	"fmt"
	"io/ioutil"
	"encoding/json"
)

func LuosimaoVerify(reqUrl, apiKey, response string) error {
	if apiKey == "" { //没有配置luosimao  就没有验证功能
		return nil
	}

	if response == "" {
		return errors.New("人机识别验证失败")
	}

	reqData := make(url.Values)
	reqData["api_key"] = []string{apiKey}
	reqData["response"] = []string{response}
	resp, err := http.PostForm(reqUrl, reqData)
	if err != nil {
		fmt.Println(err.Error())
		return errors.New("人机识别验证失败")
	}

	defer resp.Body.Close()

	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
		return errors.New("人机识别验证失败")
	}

	type LuosimaoResult struct {
		Error int    `json:"error"`
		Res   string `json:"res"`
		Msg   string `json:"msg"`
	}

	var luosimaoResult LuosimaoResult

	if err := json.Unmarshal(result,&luosimaoResult);err != nil{
		fmt.Println(err.Error())
		return errors.New("人机识别验证失败")
	}

	if luosimaoResult.Res != "Success" {
		fmt.Println("luosimaoResult.Res error ",luosimaoResult.Res)
		return errors.New("人机识别验证失败")
	}

	return nil
}

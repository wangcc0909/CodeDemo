package common

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"gopractice/model"
	"fmt"
	"errors"
	"strings"
	"mime"
	"os"
)

//上传文件
func upload(c *gin.Context) (map[string]interface{}, error) {
	file, err := c.FormFile("upFile")
	if err != nil {
		fmt.Println(err.Error())
		return nil, errors.New("参数无效")
	}

	var fileName = file.Filename

	lastIndex := strings.LastIndex(fileName, ".")
	if lastIndex < 0 {
		return nil, errors.New("无效的文件名")
	}

	var ext = fileName[lastIndex:]
	if len(ext) == 1 {
		return nil, errors.New("无效的扩展名")
	}

	mimeType := mime.TypeByExtension(ext)

	fmt.Printf("filename %s, index %d, ext %s, mimeType %s\n", fileName, lastIndex, ext, mimeType)
	if mimeType == "" && ext == ".jpge" {
		mimeType = "image/jpge"
	}

	if mimeType == "" {
		return nil, errors.New("无效的图片类型")
	}

	//创建一个上传文件的struct
	imgUpInfo := model.GenerateUploadedImgInfo(ext)
	fmt.Println(imgUpInfo.UploadDir)

	if err := os.MkdirAll(imgUpInfo.UploadDir, 0777); err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	if err := c.SaveUploadedFile(file, imgUpInfo.UploadFilePath); err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	image := model.Image{
		Title:        imgUpInfo.FileName,
		OrignalTitle: fileName,
		URL:          imgUpInfo.ImgUrl,
		Width:        0,
		Height:       0,
		Mime:         mimeType,
	}

	if err := model.DB.Create(&image).Error; err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	return map[string]interface{}{
		"id":       image.ID,
		"url":      image.URL,
		"title":    imgUpInfo.FileName,
		"original": fileName,
		"type":     image.Mime,
	}, nil

}

func UploadHandler(c *gin.Context) {

	data, err := upload(c)

	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"errNo": model.ErrorCode.ERROR,
			"msg":   err.Error(),
			"data":  gin.H{},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"errNo": model.ErrorCode.SUCCESS,
		"msg":   "success",
		"data":  data,
	})

}

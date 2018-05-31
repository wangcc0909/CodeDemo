package model

import (
	"os"
	"gopractice/config"
	"unicode/utf8"
	"gopractice/util"
	"github.com/satori/go.uuid"
	"strings"
)

type Image struct {
	ID           uint   `gorm:"primary_key" json:"id"`
	Title        string `json:"title"`
	OrignalTitle string `json:"orignalTitle"`
	URL          string `json:"url"`
	Width        int    `json:"width"`
	Height       int    `json:"height"`
	Mime         string `json:"mime"`
}

type ImageUpInfo struct {
	//目录  文件路径  文件名  UUIDname 请求URL
	UploadDir      string
	UploadFilePath string
	FileName       string
	UUIDName       string
	ImgUrl         string
}

func GenerateUploadedImgInfo(ext string) ImageUpInfo {
	sep := string(os.PathSeparator)
	uploadImgDir := config.ServerConfig.UploadImgDir
	length := utf8.RuneCountInString(uploadImgDir) //这里是对编码的判断
	lastChar := uploadImgDir[length-1:]

	ymStr := util.GetTodayYM(sep)

	var uploadDir string
	if lastChar != sep {
		uploadDir = uploadImgDir + sep + ymStr
	} else {
		uploadDir = uploadImgDir + ymStr
	}

	uuidName := uuid.Must(uuid.NewV4()).String()

	fileName := uuidName + ext

	uploadFilePath := uploadDir + sep + fileName

	imgURL := strings.Join([]string{
		"http://" + config.ServerConfig.ImgHost + config.ServerConfig.ImgPath,
		ymStr,
		fileName,
	}, "/")

	return ImageUpInfo{
		ImgUrl:         imgURL,
		UUIDName:       uuidName,
		UploadDir:      uploadDir,
		UploadFilePath: uploadFilePath,
		FileName:       fileName,
	}
}

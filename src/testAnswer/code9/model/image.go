package model

import (
	"os"
	"unicode/utf8"
	"github.com/satori/go.uuid"
	"strings"
	"readbook/utils"
)

// Image 图片
type Image struct {
	ID             uint       `gorm:"primary_key" json:"id"`
	Title          string     `json:"title"`
	OrignalTitle   string     `json:"orignalTitle"`
	URL            string     `json:"url"`
	Width          uint       `json:"width"`
	Height         uint       `json:"height"`
	Mime           string     `json:"mime"`
}

// ImageUploadedInfo 图片上传后的相关信息(目录、文件路径、文件名、UUIDName、请求URL)
type ImageUploadedInfo struct {
	UploadDir       string
	UploadFilePath  string
	Filename        string
	UUIDName        string
	ImgURL          string
}

// GenerateImgUploadedInfo 创建一个ImageUploadedInfo
func GenerateImgUploadedInfo(ext string) ImageUploadedInfo {
	sep          := string(os.PathSeparator)
	uploadImgDir := "/update/img"
	length       := utf8.RuneCountInString(uploadImgDir)
	lastChar     := uploadImgDir[length - 1:]
	ymStr        := utils.GetTodayYM(sep)

	var uploadDir string
	if lastChar != sep {
		uploadDir = uploadImgDir + sep	+ ymStr
	} else {
		uploadDir = uploadImgDir + ymStr
	}

	uuidName       := uuid.Must(uuid.NewV4()).String()
	filename       := uuidName + ext
	uploadFilePath := uploadDir + sep + filename
	imgURL         := strings.Join([]string{
		"https://127.0.0.1:8080",
		ymStr,
		filename,
	}, "/")
	return ImageUploadedInfo{
		ImgURL: imgURL,
		UUIDName: uuidName,
		Filename: filename,
		UploadDir: uploadDir,
		UploadFilePath: uploadFilePath,
	}
}

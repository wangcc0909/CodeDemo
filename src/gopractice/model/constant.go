package model

//通用的常量
const (
	//每页最大的请求数量
	MaxPageSize = 100

	//每页最小的请求数量
	MinPageSize = 5
)

const (
	//contentTypeMarkdown markdown
	ContentTypeMarkdown = 1

	//ContentTypeHTML html
	ContentTypeHTML = 2
)


//Redis的常量
const (
	//生成激活账号的链接
	ActiveTime = "activeTime"

	//用户信息
	LoginUser = "loginUser"

	//重置密码
	ResetTime = "resetTime"

)



package model

//通用的常量
const (
	//无父节点的parent_id
	NoParent = 0
	//每页最大的请求数量
	MaxPageSize = 100

	//每页最小的请求数量
	MinPageSize = 5

	//文章的名称最大字节数
	MaxNameLen = 100

	//文章的最大字节数
	MaxContentLen = 50000

)

//积分相关的常量
const (
	//ArticleScore  创建话题时增加的积分 5分
	ArticleScore = 5
)

const (
	//contentTypeMarkdown markdown
	ContentTypeMarkdown = 1

	//ContentTypeHTML html
	ContentTypeHTML = 2
)

const (
	//用户每分钟最多发表的文章数
	ArticleMinuteLimitCount = 30

	//用户每天最多发表的文章数
	ArticleDayLimitCount = 1000
)


//Redis的常量
const (
	//生成激活账号的链接
	ActiveTime = "activeTime"

	//用户信息
	LoginUser = "loginUser"

	//重置密码
	ResetTime = "resetTime"

	//用户每分钟最大发表的文章数
	ArticleMinuteLimit = "articleMinuteLimit"

	//用户每天最大发表的文章数
	ArticleDayLimit = "articleDayLimit"

)



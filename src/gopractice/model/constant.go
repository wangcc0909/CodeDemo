package model

//通用的常量
const (
	//无父节点的parent_id
	NoParent = 0
	//每页最大的请求数量
	MaxPageSize = 100

	PageSize = 20

	//每页最小的请求数量
	MinPageSize = 5

	//文章的名称最大字节数
	MaxNameLen = 100

	//最大的排序号
	MaxOrder = 1000

	//最小的排序号
	MinOrder = 0

	//文章的最大字节数
	MaxContentLen = 50000

	//最多的分类个数
	MaxCategoriesLen = 6
)

//积分相关的常量
const (
	//ArticleScore  创建话题时增加的积分 5分
	ArticleScore = 5

	//ByCollectScore 话题或投票被收藏时增加的积分
	ByCollectScore = 2

	//byCommentScore = 2 话题或投票被评论时增加的积分
	ByCommentScore = 2

	//评论话题或投票的积分
	CommentScore = 1
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

	//CommentMinuteLimitCount   用户每分钟最多发表的评论数
	CommentMinuteLimitCount = 30

	//CommentDayLimitCount  用户每天最多发表的评论数
	CommentDayLimitCount = 1000
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

	//CommentMinuteLimit 用户每分钟最多发表的评论数
	CommentMinuteLimit = "commentMinuteLimit"

	//CommentDayLimit 用户每天最多发表的评论数
	CommentDayLimit = "commentDayLimit"
)

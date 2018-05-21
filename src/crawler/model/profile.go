package model

type Profile struct {
	Name string
	Age int
	Height int
	Weight int
	InComing string
	Gender string
	Occupation string
	Education string
	HuKou string
	Car string
	Horse string
	Marriage string
	XinZuo string
}

type Item struct {
	Url string
	Id string
	Type string
	Profile Profile
}

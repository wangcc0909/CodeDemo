package parser

import (
	"crawler/engine"
	"regexp"
	"strconv"
	"crawler/model"
)

var AgeRe = regexp.MustCompile(`<td><span class="label">年龄：</span>(\d+)岁</td>`)
var HeightRe = regexp.MustCompile(`<td><span class="label">身高：</span>(\d+)CM</td>`)
var WeightRe = regexp.MustCompile(`<td><span class="label">体重：</span><span field="">(\d+)KG</span></td>`)
var InComingRe = regexp.MustCompile(`<td><span class="label">月收入：</span>([^<]+)</td>`)
var GenderRe = regexp.MustCompile(`<td><span class="label">性别：</span><span field="">([^<]+)</span></td>`)
var MarriageRe = regexp.MustCompile(`<td><span class="label">婚况：</span>([^<]+)</td>`)
var EducationRe = regexp.MustCompile(`<td><span class="label">学历：</span>([^<]+)</td>`)
var OccupationRe = regexp.MustCompile(`<td><span class="label">职业： </span>([^<]+)</td>`)
var HoKouRe = regexp.MustCompile(`<td><span class="label">籍贯：</span>([^<]+)</td>`)
var CarRe = regexp.MustCompile(`<td><span class="label">是否购车：</span><span field="">([^<]+)</span></td>`)
var HorseRe = regexp.MustCompile(`<td><span class="label">住房条件：</span><span field="">([^<]+)</span></td>`)

var GuessRe = regexp.MustCompile(`<a class="exp-user-name"[^>]+href="(http://album.zhenai.com/u/\d+)">([^<]+)</a>`)
var UrlRe = regexp.MustCompile(`http://album.zhenai.com/u/(\d+)`)

func parserProfile(contents []byte, url string, name string) engine.ParserResult {

	result := engine.ParserResult{}
	profile := model.Profile{}
	profile.Name = name

	age, err := strconv.Atoi(expectedString(contents, AgeRe))

	if err == nil {
		profile.Age = age
	}

	height, err := strconv.Atoi(expectedString(contents, HeightRe))

	if err == nil {
		profile.Height = height
	}

	weight, err := strconv.Atoi(expectedString(contents, WeightRe))
	if err == nil {
		profile.Weight = weight
	}

	profile.InComing = expectedString(contents, InComingRe)
	profile.Gender = expectedString(contents, GenderRe)
	profile.Marriage = expectedString(contents, MarriageRe)
	profile.Education = expectedString(contents, EducationRe)
	profile.Occupation = expectedString(contents, OccupationRe)
	profile.HuKou = expectedString(contents, HoKouRe)
	profile.Car = expectedString(contents, CarRe)
	profile.Horse = expectedString(contents, HorseRe)

	item := model.Item{
		Url:     url,
		Id:      expectedString(contents, UrlRe),
		Type:    "zhenai",
		Profile: profile,
	}

	result.Items = append(result.Items, item)

	maches := GuessRe.FindAllSubmatch(contents, -1)
	for _, m := range maches {
		result.Requests = append(result.Requests, engine.Request{
			Url:    string(m[1]),
			Parser: NewProfileParser(string(m[2])),
		})
	}

	return result
}

func expectedString(contents []byte, re *regexp.Regexp) string {
	result := re.FindSubmatch(contents)

	if len(result) >= 2 {
		return string(result[1])
	} else {
		return ""
	}
}

type ProfileParser struct {
	userName string
}

func (p *ProfileParser) Parse(content []byte, url string) engine.ParserResult {
	return parserProfile(content, url, p.userName)
}

func (p *ProfileParser) Name() (name string, args interface{}) {
	return "ProfileParser", p.userName
}

func NewProfileParser(userName string) *ProfileParser {
	return &ProfileParser{
		userName: userName,
	}
}

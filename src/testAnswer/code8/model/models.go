package model

type AutoGenerated struct {
	Status    int  `json:"status"`
	ErrorCode int  `json:"error_code"`
	Data      Data `json:"data"`
}

type Data struct {
	Page              int           `json:"page"`
	Tab               string        `json:"tab"`
	Lists             []Lists       `json:"lists"`
	Chinesecount      int           `json:"chinesecount"`
	Searchfull        int           `json:"searchfull"`
	Correctiontype    int           `json:"correctiontype"`
	Subjecttype       int           `json:"subjecttype"`
	Aggregation       []Aggregation `json:"aggregation"`
	Allowerr          int           `json:"allowerr"`
	Correctionsubject string        `json:"correctionsubject"`
	Correctionforce   int           `json:"correctionforce"`
	Total             int           `json:"total"`
	Istagresult       int           `json:"istagresult"`
	Istag             int           `json:"istag"`
	Correctiontip     string        `json:"correctiontip"`
	Pagesize          int           `json:"pagesize"`
}

type Lists struct {
	SongName       string     `json:"SongName"`
	OwnerCount     int        `json:"OwnerCount"`
	MvType         int        `json:"MvType"`
	TopicRemark    string     `json:"TopicRemark"`
	SQFailProcess  int        `json:"SQFailProcess"`
	Source         string     `json:"Source"`
	Bitrate        int        `json:"Bitrate"`
	HQExtName      string     `json:"HQExtName"`
	SQFileSize     int        `json:"SQFileSize"`
	ResFileSize    int        `json:"ResFileSize"`
	Duration       int        `json:"Duration"`
	MvTrac         int        `json:"MvTrac"`
	SQDuration     int        `json:"SQDuration"`
	ExtName        string     `json:"ExtName"`
	Auxiliary      string     `json:"Auxiliary"`
	SongLabel      string     `json:"SongLabel"`
	Scid           int        `json:"Scid"`
	OriSongName    string     `json:"OriSongName"`
	FailProcess    int        `json:"FailProcess"`
	SQBitrate      int        `json:"SQBitrate"`
	HQBitrate      int        `json:"HQBitrate"`
	Audioid        int        `json:"Audioid"`
	HiFiQuality    int        `json:"HiFiQuality"`
	//Grp            Grp        `json:"Grp"`
	OriOtherName   string     `json:"OriOtherName"`
	AlbumPrivilege int        `json:"AlbumPrivilege"`
	TopicURL       string     `json:"TopicUrl"`
	SuperFileHash  string     `json:"SuperFileHash"`
	ASQPrivilege   int        `json:"ASQPrivilege"`
	M4ASize        int        `json:"M4aSize"`
	AlbumName      string     `json:"AlbumName"`
	IsOriginal     int        `json:"IsOriginal"`
	Privilege      int        `json:"Privilege"`
	ResBitrate     int        `json:"ResBitrate"`
	FileHash       string     `json:"FileHash"`
	SQPayType      int        `json:"SQPayType"`
	HQPrice        int        `json:"HQPrice"`
	Type           string     `json:"Type"`
	TransParam     TransParam `json:"trans_param"`
	SourceID       int        `json:"SourceID"`
	FoldType       int        `json:"FoldType"`
	SingerID       []int      `json:"SingerId"`
	A320Privilege  int        `json:"A320Privilege"`
	ID             string     `json:"ID"`
	SuperFileSize  int        `json:"SuperFileSize"`
	SQPrivilege    int        `json:"SQPrivilege"`
	SQFileHash     string     `json:"SQFileHash"`
	AlbumID        string     `json:"AlbumID"`
	HQPrivilege    int        `json:"HQPrivilege"`
	SuperBitrate   int        `json:"SuperBitrate"`
	SuperDuration  int        `json:"SuperDuration"`
	MixSongID      string     `json:"MixSongID"`
	ResFileHash    string     `json:"ResFileHash"`
	FileSize       int        `json:"FileSize"`
	SuperExtName   string     `json:"SuperExtName"`
	HQFileHash     string     `json:"HQFileHash"`
	HQPkgPrice     int        `json:"HQPkgPrice"`
	AudioCdn       int        `json:"AudioCdn"`
	FileName       string     `json:"FileName"`
	OtherName      string     `json:"OtherName"`
	MvTotal        int        `json:"mvTotal"`
	PkgPrice       int        `json:"PkgPrice"`
	HQFileSize     int        `json:"HQFileSize"`
	HQFailProcess  int        `json:"HQFailProcess"`
	Publish        int        `json:"Publish"`
	QualityLevel   int        `json:"QualityLevel"`
	SQPrice        int        `json:"SQPrice"`
	ResDuration    int        `json:"ResDuration"`
	PublishAge     int        `json:"PublishAge"`
	Price          int        `json:"Price"`
	HQPayType      int        `json:"HQPayType"`
	SingerName     string     `json:"SingerName"`
	SQExtName      string     `json:"SQExtName"`
	MvHash         string     `json:"MvHash"`
	SQPkgPrice     int        `json:"SQPkgPrice"`
	HQDuration     int        `json:"HQDuration"`
	PayType        int        `json:"PayType"`
	HasAlbum       int        `json:"HasAlbum"`
	Accompany      int        `json:"Accompany"`
	OldCpy         int        `json:"OldCpy"`
}

type Grp struct {

}

type TransParam struct {
	RoamingAstrict   int `json:"roaming_astrict"`
	Display          int `json:"display"`
	PayBlockTpl      int `json:"pay_block_tpl"`
	MusicpackAdvance int `json:"musicpack_advance"`
}

type Aggregation struct {
	Key   string `json:"key"`
	Count int    `json:"count"`
}

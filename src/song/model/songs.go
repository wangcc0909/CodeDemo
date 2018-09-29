package model

type SongInfo struct {
	Status  int            `json:"status"`
	ErrCode int            `json:"err_code"`
	Data    SongDetailInfo `json:"data"`
}

type SongDetailInfo struct {
	Hash       string      `json:"hash"`
	Timelength int         `json:"timelength"`
	Filesize   int         `json:"filesize"`
	AudioName  string      `json:"audio_name"`
	HaveAlbum  int         `json:"have_album"`
	AlbumName  string      `json:"album_name"`
	AlbumID    int         `json:"album_id"`
	Img        string      `json:"img"`
	HaveMv     int         `json:"have_mv"`
	AuthorName string      `json:"author_name"`
	SongName   string      `json:"song_name"`
	Lyrics     string      `json:"lyrics"`
	AuthorID   interface{} `json:"author_id"`
	Privilege  int         `json:"privilege"`
	Privilege2 string      `json:"privilege2"`
	PlayURL    string      `json:"play_url"`
	Authors    []Author    `json:"authors"`
	Bitrate    int         `json:"bitrate"`
}

type Author struct {
	IsPublish  string `json:"is_publish"`
	AuthorID   interface{} `json:"author_id"`
	Avatar     string `json:"avatar"`
	AuthorName string `json:"author_name"`
}

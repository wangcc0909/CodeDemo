package model

type Song struct {
	Hash       string `json:"hash"`
	Timelength int    `json:"timelength"`
	Filesize   int    `json:"filesize"`
	AudioName  string `json:"audio_name"`
	HaveAlbum  int    `json:"have_album"`
	AlbumName  string `json:"album_name"`
	AlbumID    int    `json:"album_id"`
	Img        string `json:"img"`
	HaveMv     int    `json:"have_mv"`
	AuthorName string `json:"author_name"`
	SongName   string `json:"song_name"`
	Lyrics     string `json:"lyrics"`
	AuthorID   string `json:"author_id"`
	PlayURL    string `json:"play_url"`
	Bitrate    int    `json:"bitrate"`
}

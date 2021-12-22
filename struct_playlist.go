package gotube

type SPlaylist struct {
	ID       string      `json:"playlist_id"`
	Title    string      `json:"playlist_title"`
	Uploader *SChannel   `json:""`
	Videos   SVideoSlice `json:"entries"`
}

package gotube

// Holds information about a single playlist and it's videos.
type SPlaylist struct {
	ID       string      `json:"playlist_id"`    //
	Title    string      `json:"playlist_title"` //
	Uploader *SChannel   `json:""`               //
	Videos   SVideoSlice `json:"entries"`        //
}

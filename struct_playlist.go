package gotube

// This holds basic information about a playlist.
type SPlaylist struct {
	ID       string      `json:"playlist_id"`    // The playlist ID.
	Title    string      `json:"playlist_title"` // The playlist title.
	Uploader *SChannel   `json:""`               // The uploader of the playlist.
	Videos   SVideoSlice `json:"entries"`        // The videos in the playlist.
}

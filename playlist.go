package gotube

// This holds basic information about a playlist.
type Playlist struct {
	ID       string   `json:"playlist_id"`    // The playlist ID.
	Title    string   `json:"playlist_title"` // The playlist title.
	Uploader *Channel `json:""`               // The uploader of the playlist.
	Videos   Videos   `json:"entries"`        // The videos in the playlist.
}

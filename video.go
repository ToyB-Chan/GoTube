package gotube

// This holds all important information about a video.
type Video struct {
	ID            string     `json:"id"`             // The video ID.
	URL           string     `json:"webpage_url"`    // The video URL.
	Title         string     `json:"title"`          // The video title.
	AltTitle      string     `json:"alt_title"`      // The video alternative title.
	Description   string     `json:"description"`    // The video description.
	Tags          []string   `json:"tags"`           // The video tags.
	ViewCount     int        `json:"view_count"`     // The number of views.
	LikeCount     int        `json:"like_count"`     // The number of likes.
	Duration      float32    `json:"duration"`       // The video duration in seconds.
	UploadDate    string     `json:"upload_date"`    // The upload date in unix time.
	AgeLimit      int        `json:"age_limit"`      // The age limit of the video.
	IsLive        bool       `json:"is_live"`        // If the video is live.
	PlaylistIndex int        `json:"playlist_index"` // The index of the video in the playlist.
	Uploader      *Channel   `json:""`               // The uploader of the video.
	Thumbnails    Thumbnails `json:"thumbnails"`     // The thumbnails of the video.
	Streams       Streams    `json:"formats"`        // The streams of the video.
	//AutomaticCaptions      map[string][]*SCaption `json:"automatic_captions"`
}

type Videos []*Video

// Returns true if the video is age restricted.
func (v *Video) IsAgeRestricted() bool {
	return v.AgeLimit > 0
}

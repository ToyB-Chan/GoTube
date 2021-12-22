package gotube

// Holds information about a single video.
type SVideo struct {
	ID            string          `json:"id"`             //
	URL           string          `json:"webpage_url"`    //
	Title         string          `json:"title"`          //
	AltTitle      string          `json:"alt_title"`      //
	Description   string          `json:"description"`    //
	Tags          []string        `json:"tags"`           //
	ViewCount     int             `json:"view_count"`     //
	LikeCount     int             `json:"like_count"`     //
	Duration      float32         `json:"duration"`       // In seconds.
	UploadDate    string          `json:"upload_date"`    // In Unix.
	AgeLimit      int             `json:"age_limit"`      //
	IsLive        bool            `json:"is_live"`        //
	PlaylistIndex int             `json:"playlist_index"` //
	Uploader      *SChannel       `json:""`               //
	Thumbnails    SThumbnailSlice `json:"thumbnails"`     //
	Streams       SStreamSlice    `json:"formats"`        //
	//AutomaticCaptions      map[string][]*SCaption `json:"automatic_captions"`
}

type SVideoSlice []*SVideo

// Returns if the video is age restricted.
func (Me *SVideo) IsAgeRestricted() bool {
	return Me.AgeLimit > 0
}

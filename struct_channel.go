package gotube

// This holds basic information about a channel.
type SChannel struct {
	ID   string `json:"uploader_id"`  // The channel ID.
	URL  string `json:"uploader_url"` // The channel URL.
	Name string `json:"uploader"`     // The channel name.
}
